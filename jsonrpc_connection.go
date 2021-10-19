package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"

	"go.bug.st/json"
)

// Connection is a JSON RPC connection for LSP protocol
type Connection struct {
	in      *bufio.Reader
	out     io.Writer
	errors  []error
	handler RequestHandler
}

// RequestHandler handles requests from a jsonrpc Connection.
// If id is nil the message is a 'notification' otherwise is a 'request'.
type RequestHandler func(ctx context.Context, id *IntOrString, method string, params *ArrayOrObject) error

// NewConnection starts a new
func NewConnection(in io.Reader, out io.Writer, handler RequestHandler) *Connection {
	conn := &Connection{
		in:      bufio.NewReader(in),
		out:     out,
		handler: handler,
	}
	go conn.inLoop()
	return conn
}

func (c *Connection) inLoop() {
	in := textproto.NewReader(c.in)
	for {
		head, err := in.ReadMIMEHeader()
		if err != nil {
			c.errors = append(c.errors, err)
			c.Close()
			return
		}

		httpHeader := http.Header(head)
		l := httpHeader.Get("Content-Length")
		httpHeader.Del("Content-Length")
		dataLen, err := strconv.Atoi(l)
		if err != nil {
			c.errors = append(c.errors, err)
			c.Close()
			return
		}

		jsonData := make([]byte, dataLen)
		if _, err := io.ReadFull(in.R, jsonData); err != nil {
			c.errors = append(c.errors, err)
			c.Close()
			return
		}

		go c.handleRequest(jsonData)
	}
}

func (c *Connection) handleRequest(jsonData []byte) {
	var req RequestMessage
	if err := json.Unmarshal(jsonData, &req); err == nil {
		fmt.Printf("REQUEST: %+v\n", req)
		c.handler(context.Background(), &req.ID, req.Method, req.Params)
		return
	}

	var notif NotificationMessage
	if err := json.Unmarshal(jsonData, &notif); err == nil {
		fmt.Printf("NOTIFICATION: %+v\n", notif)
		c.handler(context.Background(), nil, notif.Method, notif.Params)
		return
	}

	c.errors = append(c.errors, fmt.Errorf("invalid request: %s", string(jsonData)))
	c.Close()
}

func (c *Connection) Close() {

}

func (c *Connection) Send() {

}
