package lsp

type InitializeParams struct {
	WorkDoneProgressParams

	// The process Id of the parent process that started the server. Is null if
	// the process has not been started by another process. If the parent
	// process is not alive then the server should exit (see exit notification)
	// its process.
	ProcessID IntOrNull `json:"processId,required"`

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
	RootPath *StringOrNull `json:"rootPath,omitempty"`

	// The rootUri of the workspace. Is null if no
	// folder is open. If both `rootPath` and `rootUri` are set
	// `rootUri` wins.
	//
	// @deprecated in favour of `workspaceFolders`
	RootURI DocumentURIOrNull `json:"rootUri,required"`

	// User provided initialization options.
	InitializationOptions *Any `json:"initializationOptions,omitempty"`

	// The capabilities provided by the client (editor or tool)
	// TODO: capabilities: ClientCapabilities

	// The initial trace setting. If omitted trace is disabled ('off').
	// TODO: trace?: TraceValue

	// The workspace folders configured in the client when the server starts.
	// This property is only available if the client supports workspace folders.
	// It can be `null` if the client supports workspace folders but none are
	// configured.
	//
	// @since 3.6.0
	// TODO: workspaceFolders?: WorkspaceFolder[] | null
}

// InitializedParams The initialized notification is sent from the client
// to the server after the client received the result of the initialize
// equest but before the client is sending any other request or notification
// to the server. The server can use the initialized notification for
// example to dynamically register capabilities. The initialized
// notification may only be sent once.
type InitializedParams struct{}
