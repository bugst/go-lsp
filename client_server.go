package lsp

import "context"

// Client is an LSP client
type Client struct{}

// Server is an LSP Server
type Server struct{}

// ServerHandler is an LSP Server message handler
type ServerHandler interface {
	Initialize(ctx context.Context, conn Connection, params InitializeParams)
}
