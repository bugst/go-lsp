//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package jsonrpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/json"
)

// Connection is a JSON RPC connection for LSP protocol
type Connection struct {
	in                  *bufio.Reader
	out                 io.Writer
	outMutex            sync.Mutex
	errorHandler        func(error)
	requestHandler      RequestHandler
	notificationHandler NotificationHandler
	logger              Logger
	loggerMutex         sync.Mutex

	activeInRequests      map[string]*inRequest
	activeInRequestsMutex sync.Mutex

	activeOutRequests      map[string]*outRequest
	activeOutRequestsMutex sync.Mutex
	lastOutRequestsIndex   uint64
}

type inRequest struct {
	cancel func()
}

type outRequest struct {
	resultChan chan<- *outResponse
	method     string
}

type outResponse struct {
	reqResult json.RawMessage
	reqError  *ResponseError
}

// RequestHandler handles requests from a jsonrpc Connection.
type RequestHandler func(ctx context.Context, logger FunctionLogger, method string, params json.RawMessage, respCallback func(result json.RawMessage, err *ResponseError))

// NotificationHandler handles notifications from a jsonrpc Connection.
type NotificationHandler func(logger FunctionLogger, method string, params json.RawMessage)

// NewConnection starts a new
func NewConnection(in io.Reader, out io.Writer, requestHandler RequestHandler, notificationHandler NotificationHandler, errorHandler func(error)) *Connection {
	conn := &Connection{
		in:                  bufio.NewReader(in),
		out:                 out,
		requestHandler:      requestHandler,
		notificationHandler: notificationHandler,
		errorHandler:        errorHandler,
		activeInRequests:    map[string]*inRequest{},
		activeOutRequests:   map[string]*outRequest{},
		logger:              NullLogger{},
	}
	return conn
}

func (c *Connection) SetLogger(l Logger) {
	c.loggerMutex.Lock()
	c.logger = l
	c.loggerMutex.Unlock()
}

func (c *Connection) Run() {
	in := textproto.NewReader(c.in)
	for {
		start := time.Now()

		head, err := in.ReadMIMEHeader()
		if err != nil {
			c.errorHandler(err)
			c.Close()
			return
		}

		httpHeader := http.Header(head)
		l := httpHeader.Get("Content-Length")
		dataLen, err := strconv.Atoi(l)
		if err != nil {
			c.errorHandler(err)
			c.Close()
			return
		}

		jsonData := make([]byte, dataLen)
		if n, err := io.ReadFull(in.R, jsonData); err != nil {
			c.errorHandler(err)
			c.Close()
			return
		} else if n != dataLen {
			c.errorHandler(fmt.Errorf("expected %d bytes but %d have been read", dataLen, n))
		}

		elapsed := time.Since(start)
		c.loggerMutex.Lock()
		c.logger.LogIncomingDataDelay(elapsed)
		c.loggerMutex.Unlock()

		c.handleIncomingData(jsonData)
	}
}

func (c *Connection) handleIncomingData(jsonData []byte) {
	var req RequestMessage
	var notif NotificationMessage
	var resp ResponseMessage
	if err := json.Unmarshal(jsonData, &req); err == nil {
		c.handleIncomingRequest(&req)
	} else if err := json.Unmarshal(jsonData, &notif); err == nil {
		c.handleIncomingNotification(&notif)
	} else if err := json.Unmarshal(jsonData, &resp); err == nil {
		c.handleIncomingResponse(&resp)
	} else {
		c.errorHandler(fmt.Errorf("invalid request: %s", string(jsonData)))
		c.Close()
	}
}

func (c *Connection) handleIncomingRequest(req *RequestMessage) {
	id := string(req.ID)
	ctx, cancel := context.WithCancel(context.Background())

	c.activeInRequestsMutex.Lock()
	c.activeInRequests[id] = &inRequest{
		cancel: cancel,
	}
	c.activeInRequestsMutex.Unlock()

	c.loggerMutex.Lock()
	logger := c.logger.LogIncomingRequest(id, req.Method, req.Params)
	c.loggerMutex.Unlock()

	c.requestHandler(ctx, logger, req.Method, req.Params, func(result json.RawMessage, resultErr *ResponseError) {
		c.activeInRequestsMutex.Lock()
		c.activeInRequests[id].cancel()
		delete(c.activeInRequests, id)
		c.activeInRequestsMutex.Unlock()

		c.loggerMutex.Lock()
		c.logger.LogOutgoingResponse(id, req.Method, result, resultErr)
		c.loggerMutex.Unlock()

		resp := &ResponseMessage{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
			Error:   resultErr,
		}
		if sendErr := c.send(resp); sendErr != nil {
			c.errorHandler(fmt.Errorf("error sending response: %s", sendErr))
			c.Close()
		}
	})
}

func (c *Connection) handleIncomingNotification(notif *NotificationMessage) {
	if notif.Method == "$/cancelRequest" {
		// Send cancelation signal and exit
		var params CancelParams
		if err := json.Unmarshal(notif.Params, &params); err != nil {
			c.errorHandler(fmt.Errorf("invalid cancelRequest: %s", err))
			return
		}
		c.cancelIncomingRequest(params.ID)
		return
	}

	c.loggerMutex.Lock()
	logger := c.logger.LogIncomingNotification(notif.Method, notif.Params)
	c.loggerMutex.Unlock()

	c.notificationHandler(logger, notif.Method, notif.Params)
}

func (c *Connection) handleIncomingResponse(resp *ResponseMessage) {
	var id string
	if err := json.Unmarshal(resp.ID, &id); err != nil {
		c.errorHandler(fmt.Errorf("invalid ID in request response '%v': %w", resp.ID, err))
		c.Close()
		return
	}

	c.activeOutRequestsMutex.Lock()
	req, ok := c.activeOutRequests[id]
	if ok {
		delete(c.activeOutRequests, id)
	}
	c.activeOutRequestsMutex.Unlock()

	if !ok {
		c.errorHandler(fmt.Errorf("invalid ID in request response '%s': double answer or request not sent", id))
		c.Close()
		return
	}

	req.resultChan <- &outResponse{
		reqResult: resp.Result,
		reqError:  resp.Error,
	}
}

func (c *Connection) cancelIncomingRequest(id json.RawMessage) {
	c.activeInRequestsMutex.Lock()
	if req, ok := c.activeInRequests[string(id)]; ok {
		c.loggerMutex.Lock()
		c.logger.LogIncomingCancelRequest(string(id))
		c.loggerMutex.Unlock()

		req.cancel()
	}
	c.activeInRequestsMutex.Unlock()
}

func (c *Connection) Close() {
}

func (c *Connection) SendRequest(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, *ResponseError, error) {
	id := fmt.Sprintf("%d", atomic.AddUint64(&c.lastOutRequestsIndex, 1))
	encodedId, err := json.Marshal(id)
	if err != nil {
		// should never happen...
		panic("internal error creating RequestMessage")
	}
	req := RequestMessage{
		JSONRPC: "2.0",
		ID:      encodedId,
		Method:  method,
		Params:  params,
	}

	c.loggerMutex.Lock()
	c.logger.LogOutgoingRequest(id, method, params)
	c.loggerMutex.Unlock()

	resultChan := make(chan *outResponse, 1)
	c.activeOutRequestsMutex.Lock()
	err = c.send(req)
	if err == nil {
		c.activeOutRequests[id] = &outRequest{
			resultChan: resultChan,
			method:     method,
		}
	}
	c.activeOutRequestsMutex.Unlock()
	if err != nil {
		return nil, nil, fmt.Errorf("sending request: %w", err)
	}

	// Wait the response or send cancel request if requested from context
	var result *outResponse
	select {
	case result = <-resultChan:
		// got result, do nothing

	case <-ctx.Done():
		c.activeOutRequestsMutex.Lock()
		_, active := c.activeOutRequests[id]
		c.activeOutRequestsMutex.Unlock()
		if active {
			if notif, err := json.Marshal(CancelParams{ID: encodedId}); err != nil {
				// should never happen
				panic("internal error: failed json encoding")
			} else {
				c.loggerMutex.Lock()
				c.logger.LogOutgoingCancelRequest(id)
				c.loggerMutex.Unlock()

				_ = c.SendNotification("$/cancelRequest", notif) // ignore error (it won't matter anyway)
			}
		}

		// After cancelation wait for result...
		result = <-resultChan
	}

	c.loggerMutex.Lock()
	c.logger.LogIncomingResponse(id, method, result.reqResult, result.reqError)
	c.loggerMutex.Unlock()

	return result.reqResult, result.reqError, nil
}

func (c *Connection) SendNotification(method string, params json.RawMessage) error {
	c.loggerMutex.Lock()
	c.logger.LogOutgoingNotification(method, params)
	c.loggerMutex.Unlock()

	notif := NotificationMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	if err := c.send(notif); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}
	return nil
}

func (c *Connection) send(data interface{}) error {
	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}

	start := time.Now()
	c.outMutex.Lock()
	defer c.outMutex.Unlock()
	if _, err := fmt.Fprintf(c.out, "Content-Length: %d\r\n\r\n", len(buff)); err != nil {
		return err
	}
	for len(buff) > 0 {
		n, err := c.out.Write(buff)
		if err != nil {
			return err
		}
		buff = buff[n:]
	}
	elapsed := time.Since(start)
	c.loggerMutex.Lock()
	c.logger.LogOutgoingDataDelay(elapsed)
	c.loggerMutex.Unlock()
	return nil
}
