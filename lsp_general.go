package lsp

import "encoding/json"

type InitializeParams struct {
	WorkDoneProgressParams

	// The process Id of the parent process that started the server. Is null if
	// the process has not been started by another process. If the parent
	// process is not alive then the server should exit (see exit notification)
	// its process.
	ProcessID *int `json:"processId,required"`

	// Information about the client
	//
	// @since 3.15.0
	ClientInfo *struct {
		// The name of the client as defined by the client.
		Name string `json:"name,required"`

		// The client's version as defined by the client.
		Version *string `json:"version,omitempty"`
	} `json:"clientInfo,omitempty"`

	// The locale the client is currently showing the user interface
	// in. This must not necessarily be the locale of the operating
	// system.
	//
	// Uses IETF language tags as the value's syntax
	// (See https://en.wikipedia.org/wiki/IETF_language_tag)
	//
	// @since 3.16.0
	Locale *string `json:"locale,omitempty"`

	// The rootPath of the workspace. Is null
	// if no folder is open.
	//
	// @deprecated in favour of `rootUri`.
	RootPath *string `json:"rootPath,omitempty"`

	// The rootUri of the workspace. Is null if no
	// folder is open. If both `rootPath` and `rootUri` are set
	// `rootUri` wins.
	//
	// @deprecated in favour of `workspaceFolders`
	RootURI *DocumentURI `json:"rootUri,required"`

	// User provided initialization options.
	InitializationOptions json.RawMessage `json:"initializationOptions,omitempty"`

	// The capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities,required"`

	// The initial trace setting. If omitted trace is disabled ('off').
	Trace *TraceValue `json:"trace,omitempy"`

	// The workspace folders configured in the client when the server starts.
	// This property is only available if the client supports workspace folders.
	// It can be `null` if the client supports workspace folders but none are
	// configured.
	//
	// @since 3.6.0
	WorkspaceFolders *[]WorkspaceFolder `json:"workspaceFolders"`
}

type TraceValue string

const TraceValueOff TraceValue = "off"
const TraceValueMessages TraceValue = "messages"
const TraceValueVerbose TraceValue = "verbose"

type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI DocumentURI `json:"uri,required"`

	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name,required"`
}

type InitializeResult struct {
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities,required"`

	// Information about the server.
	//
	// @since 3.15.0
	ServerInfo *struct {
		// The name of the server as defined by the server.
		Name string `json:"name,required"`

		// The server's version as defined by the server.
		Version *string `json:"version,omitempty"`
	} `json:"serverInfo,omitempty"`
}

// If the protocol version provided by the client can't be handled by the
// server.
//
// @deprecated This initialize error got replaced by client capabilities.
// There is no version handshake in version 3.0x
const InitializeErrorUnknownProtocolVersion ErrorCode = 1

type InitializeError struct {
	// Indicates whether the client execute the following retry logic:
	// (1) show the message provided by the ResponseError to the user
	// (2) user selects retry or cancel
	// (3) if user selected retry the initialize method is sent again.
	Retry bool `json:"retry,required"`
}

// InitializedParams The initialized notification is sent from the client
// to the server after the client received the result of the initialize
// equest but before the client is sending any other request or notification
// to the server. The server can use the initialized notification for
// example to dynamically register capabilities. The initialized
// notification may only be sent once.
type InitializedParams struct{}

type LogTraceParams struct {
	// The message to be logged.
	Message string `json:"message,required"`

	// Additional information that can be computed if the `trace` configuration
	// is set to `'verbose'`
	Verbose *string `json:"verbose,omitempty"`
}

type SetTraceParams struct {
	// The new value that should be assigned to the trace setting.
	Value TraceValue `json:"value,required"`
}
