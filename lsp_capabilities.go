package lsp

import "encoding/json"

// ClientCapabilities Workspace specific client capabilities.
type ClientCapabilities struct {
	Workspace *struct {
		// The client supports applying batch edits
		// to the workspace by supporting the request
		// 'workspace/applyEdit'
		ApplyEdit *bool `json:"applyEdit,omitempty"`

		// Capabilities specific to `WorkspaceEdit`s
		WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`

		// Capabilities specific to the `workspace/didChangeConfiguration`
		// notification.
		// didChangeConfiguration?: DidChangeConfigurationClientCapabilities;

		// Capabilities specific to the `workspace/didChangeWatchedFiles`
		// notification.
		// didChangeWatchedFiles?: DidChangeWatchedFilesClientCapabilities;

		// Capabilities specific to the `workspace/symbol` request.
		// symbol?: WorkspaceSymbolClientCapabilities;

		// Capabilities specific to the `workspace/executeCommand` request.
		// executeCommand?: ExecuteCommandClientCapabilities;

		// The client has support for workspace folders.
		//
		// @since 3.6.0
		WorkspaceFolders *bool `json:"workspaceFolders,omitempty"`

		// The client supports `workspace/configuration` requests.
		//
		// @since 3.6.0
		Configuration *bool `json:"configuration,omitempty"`

		// Capabilities specific to the semantic token requests scoped to the
		// workspace.
		//
		// @since 3.16.0
		//  semanticTokens?: SemanticTokensWorkspaceClientCapabilities;

		// Capabilities specific to the code lens requests scoped to the
		// workspace.
		//
		// @since 3.16.0
		// codeLens?: CodeLensWorkspaceClientCapabilities;

		// The client has support for file requests/notifications.
		//
		// @since 3.16.0
		FileOperations *struct {

			// Whether the client supports dynamic registration for file
			// requests/notifications.
			DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`

			// The client has support for sending didCreateFiles notifications.
			DidCreate *bool `json:"didCreate,omitempty"`

			// The client has support for sending willCreateFiles requests.
			WillCreate *bool `json:"willCreate,omitempty"`

			// The client has support for sending didRenameFiles notifications.
			DidRename *bool `json:"didRename,omitempty"`

			// The client has support for sending willRenameFiles requests.
			WillRename *bool `json:"willRename,omitempty"`

			// The client has support for sending didDeleteFiles notifications.
			DidDelete *bool `json:"didDelete,omitempty"`

			// The client has support for sending willDeleteFiles requests.
			WillDelete *bool `json:"willDelete,omitempty"`
		} `json:"fileOperations,omitempty"`
	} `json:"workspace,omitempty"`

	// Text document specific client capabilities.
	// textDocument?: TextDocumentClientCapabilities;

	// Window specific client capabilities.
	Window *struct {

		// Whether client supports handling progress notifications. If set
		// servers are allowed to report in `workDoneProgress` property in the
		// request specific server capabilities.
		//
		// @since 3.15.0
		WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`

		// Capabilities specific to the showMessage request
		//
		// @since 3.16.0
		// showMessage?: ShowMessageRequestClientCapabilities;

		// Client capabilities for the show document request.
		//
		// @since 3.16.0
		// showDocument?: ShowDocumentClientCapabilities;
	} `json:"window,omitempty"`

	// General client capabilities.
	//
	// @since 3.16.0
	General *struct {

		// Client capability that signals how the client
		// handles stale requests (e.g. a request
		// for which the client will not process the response
		// anymore since the information is outdated).
		//
		// @since 3.17.0
		StaleRequestSupport *struct {
			// The client will actively cancel the request.
			Cancel bool `json:"cancel,required"`

			// The list of requests for which the client
			// will retry the request if it receives a
			// response with error code `ContentModified``
			RetryOnContentModified []string `json:"retryOnContentModified,required"`
		} `json:"staleRequestSupport,omitempty"`

		// Client capabilities specific to regular expressions.
		//
		// @since 3.16.0
		// regularExpressions?: RegularExpressionsClientCapabilities;

		// Client capabilities specific to the client's markdown parser.
		//
		// @since 3.16.0
		// markdown?: MarkdownClientCapabilities;
	} `json:"general,omitempty"`

	// Experimental client capabilities.
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

type WorkspaceEditClientCapabilities struct {

	// The client supports versioned document changes in `WorkspaceEdit`s
	DocumentChanges *bool `json:"documentChanges,omitempty"`

	// The resource operations the client supports. Clients should at least
	// support 'create', 'rename' and 'delete' files and folders.
	//
	// @since 3.13.0
	ResourceOperations *[]ResourceOperationKind `json:"resourceOperations,omitempty"`

	// The failure handling strategy of a client if applying the workspace edit
	// fails.
	//
	// @since 3.13.0
	FailureHandling *FailureHandlingKind `json:"failureHandling,omitempty"`

	// Whether the client normalizes line endings to the client specific
	// setting.
	// If set to `true` the client will normalize line ending characters
	// in a workspace edit to the client specific new line character(s).
	//
	// @since 3.16.0
	NormalizesLineEndings *bool `json:"normalizesLineEndings,omitempty"`

	// Whether the client in general supports change annotations on text edits,
	// create file, rename file and delete file changes.
	//
	// @since 3.16.0
	ChangeAnnotationSupport *struct {
		// Whether the client groups edits with equal labels into tree nodes,
		// for instance all edits labelled with "Changes in Strings" would
		// be a tree node.
		GroupsOnLabel *bool `json:"groupsOnLabel,omitempty"`
	} `json:"changeAnnotationSupport,omitempty"`
}

type ResourceOperationKind string

// ResourceOperationKindCreate: Supports creating new files and folders.
const ResourceOperationKindCreate = ResourceOperationKind("create")

// ResourceOperationKindRename: Supports renaming existing files and folders.
const ResourceOperationKindRename = ResourceOperationKind("rename")

// ResourceOperationKindDelete: Supports deleting existing files and folders.
const ResourceOperationKindDelete = ResourceOperationKind("delete")

type FailureHandlingKind string

// FailureHandlingKindAbort: Applying the workspace change is simply aborted if one of the changes
// provided fails. All operations executed before the failing operation
// stay executed.
const FailureHandlingKindAbort FailureHandlingKind = "abort"

// FailureHandlingKindTransactional: All operations are executed transactional. That means they either all
// succeed or no changes at all are applied to the workspace.
const FailureHandlingKindTransactional FailureHandlingKind = "transactional"

// FailureHandlingKindTextOnlyTransactional: If the workspace edit contains only textual file changes they are
// executed transactional. If resource changes (create, rename or delete
// file) are part of the change the failure handling strategy is abort.
const FailureHandlingKindTextOnlyTransactional FailureHandlingKind = "textOnlyTransactional"

// FailureHandlingKindUndo: The client tries to undo the operations already executed. But there is no
// guarantee that this is succeeding.
const FailureHandlingKindUndo FailureHandlingKind = "undo"
