package lsp

import (
	"context"
	"io"

	"go.bug.st/json"
	"go.bug.st/lsp/jsonrpc"
)

// Client is an LSP client
type Client struct{}

// Server is an LSP Server
type Server struct {
	handler ServerHandler
	conn    *jsonrpc.Connection
}

// ServerHandler is an LSP Server message handler
type ServerHandler interface {
	Initialize(ctx context.Context, conn jsonrpc.Connection, params InitializeParams)
}

func NewServer(in io.Reader, out io.Writer, handler ServerHandler) *Server {
	serv := &Server{}
	serv.handler = handler
	serv.conn = jsonrpc.NewConnection(in, out, serv.RequestHandler, serv.NotificationHandler, serv.ErrorHandler)
	return serv
}

func (serv *Server) ErrorHandler(e error) {
}

func (serv *Server) NotificationHandler(ctx context.Context, method string, params json.RawMessage) {
}

func (serv *Server) RequestHandler(ctx context.Context, method string, params json.RawMessage, respCallback func(json.RawMessage, *jsonrpc.ResponseError)) {
	switch method {
	case "initialize":
	default:
	}
	respCallback(nil, nil)
}

func (serv *Server) Run() {
	serv.conn.Run()
}
