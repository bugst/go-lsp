package lsp

import (
	"go.bug.st/lsp/jsonrpc"
)

type WorkDoneProgressOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken *jsonrpc.ProgressToken `json:"workDoneToken,omitempty"`
}

type PartialResultParams struct {
	// An optional token that a server can use to report partial results (e.g.
	// streaming) to the client.
	PartialResultToken *jsonrpc.ProgressToken `json:"partialResultToken,omitempty"`
}
