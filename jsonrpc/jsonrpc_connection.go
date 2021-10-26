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

	activeInRequests      map[string]*request
	activeInRequestsMutex sync.Mutex

	activeOutRequests      map[string]chan<- *outRequest
	activeOutRequestsMutex sync.Mutex
	lastOutRequestsIndex   uint64
}

type request struct {
	cancel func()
}

type outRequest struct {
	reqResult json.RawMessage
	reqError  *ResponseError
}

// RequestHandler handles requests from a jsonrpc Connection.
type RequestHandler func(ctx context.Context, method string, params json.RawMessage, respCallback func(result json.RawMessage, err *ResponseError))

// NotificationHandler handles notifications from a jsonrpc Connection.
type NotificationHandler func(ctx context.Context, method string, params json.RawMessage)

// NewConnection starts a new
func NewConnection(in io.Reader, out io.Writer, requestHandler RequestHandler, notificationHandler NotificationHandler, errorHandler func(error)) *Connection {
	conn := &Connection{
		in:                  bufio.NewReader(in),
		out:                 out,
		requestHandler:      requestHandler,
		notificationHandler: notificationHandler,
		errorHandler:        errorHandler,
		activeInRequests:    map[string]*request{},
		activeOutRequests:   map[string]chan<- *outRequest{},
	}
	return conn
}

func (c *Connection) Run() {
	in := textproto.NewReader(c.in)
	for {
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
	c.activeInRequests[id] = &request{
		cancel: cancel,
	}
	c.activeInRequestsMutex.Unlock()

	c.requestHandler(ctx, req.Method, req.Params, func(result json.RawMessage, resultErr *ResponseError) {
		c.activeInRequestsMutex.Lock()
		c.activeInRequests[id].cancel()
		delete(c.activeInRequests, id)
		c.activeInRequestsMutex.Unlock()

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

	c.notificationHandler(context.Background(), notif.Method, notif.Params)
}

func (c *Connection) handleIncomingResponse(req *ResponseMessage) {
	var id string
	if err := json.Unmarshal(req.ID, &id); err != nil {
		c.errorHandler(fmt.Errorf("invalid ID in request response '%v': %w", req.ID, err))
		c.Close()
		return
	}

	c.activeOutRequestsMutex.Lock()
	resultChan, ok := c.activeOutRequests[id]
	c.activeOutRequestsMutex.Unlock()

	if !ok {
		c.errorHandler(fmt.Errorf("invalid ID in request response '%s': double answer or request not sent", id))
		c.Close()
		return
	}

	delete(c.activeOutRequests, id)
	resultChan <- &outRequest{
		reqResult: req.Result,
		reqError:  req.Error,
	}
}

func (c *Connection) cancelIncomingRequest(id json.RawMessage) {
	c.activeInRequestsMutex.Lock()
	if req, ok := c.activeInRequests[string(id)]; ok {
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

	resultChan := make(chan *outRequest, 1)
	c.activeOutRequestsMutex.Lock()
	err = c.send(req)
	if err == nil {
		c.activeOutRequests[id] = resultChan
	}
	c.activeOutRequestsMutex.Unlock()
	if err != nil {
		return nil, nil, fmt.Errorf("sending request: %w", err)
	}

	// Wait the response or send cancel request if requested from context
	select {
	case <-ctx.Done():
		c.activeOutRequestsMutex.Lock()
		_, active := c.activeOutRequests[id]
		c.activeOutRequestsMutex.Unlock()
		if active {
			if notif, err := json.Marshal(CancelParams{ID: encodedId}); err != nil {
				// should never happen
				panic("internal error: failed json encoding")
			} else {
				_ = c.SendNotification("$/cancelRequest", notif) // ignore error (it won't matter anyway)
			}
		}
	case result := <-resultChan:
		return result.reqResult, result.reqError, nil
	}

	result := <-resultChan
	return result.reqResult, result.reqError, nil
}

func (c *Connection) SendNotification(method string, params json.RawMessage) error {
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
	return nil
}
