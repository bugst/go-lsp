package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"

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
}

// RequestHandler handles requests from a jsonrpc Connection.
type RequestHandler func(ctx context.Context, method string, params *ArrayOrObject, respCallback func(result Any, err error))

// NotificationHandler handles notifications from a jsonrpc Connection.
type NotificationHandler func(ctx context.Context, method string, params *ArrayOrObject)

// NewConnection starts a new
func NewConnection(in io.Reader, out io.Writer, requestHandler RequestHandler, notificationHandler NotificationHandler, errorHandler func(error)) *Connection {
	conn := &Connection{
		in:                  bufio.NewReader(in),
		out:                 out,
		requestHandler:      requestHandler,
		notificationHandler: notificationHandler,
		errorHandler:        errorHandler,
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
		c.handleRequest(jsonData)
	}
}

func (c *Connection) handleRequest(jsonData []byte) {
	var req RequestMessage
	if err := json.Unmarshal(jsonData, &req); err == nil {
		c.requestHandler(context.Background(), req.Method, req.Params, func(result Any, resultErr error) {
			resp := &ResponseMessage{
				Message: Message{JSONRPC: "2.0"},
				ID:      req.ID,
				Result:  result,
			}
			_ = resultErr // TODO...
			if sendErr := c.Send(resp); sendErr != nil {
				c.errorHandler(fmt.Errorf("error sending response: %s", sendErr))
				c.Close()
			}
		})

		return
	}

	var notif NotificationMessage
	if err := json.Unmarshal(jsonData, &notif); err == nil {
		c.notificationHandler(context.Background(), notif.Method, notif.Params)
		return
	}

	c.errorHandler(fmt.Errorf("invalid request: %s", string(jsonData)))
	c.Close()
}

func (c *Connection) Close() {

}

func (c *Connection) Send(data interface{}) error {
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
