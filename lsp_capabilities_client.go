package lsp

import (
	"fmt"

	"go.bug.st/json"
)

// ClientCapabilities Workspace specific client capabilities.
type ClientCapabilities struct {
	Workspace *struct {
		// The client supports applying batch edits
		// to the workspace by supporting the request
		// 'workspace/applyEdit'
		ApplyEdit bool `json:"applyEdit,omitempty"`

		// Capabilities specific to `WorkspaceEdit`s
		WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`

		// Capabilities specific to the `workspace/didChangeConfiguration`
		// notification.
		DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`

		// Capabilities specific to the `workspace/didChangeWatchedFiles`
		// notification.
		DidChangeWatchedFiles *DidChangeWatchedFilesClientCapabilities `json:"didChangeWatchedFiles,omitempty"`

		// Capabilities specific to the `workspace/symbol` request.
		Symbol *WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`

		// Capabilities specific to the `workspace/executeCommand` request.
		ExecuteCommand *ExecuteCommandClientCapabilities `json:"executeCommand,omitempty"`

		// The client has support for workspace folders.
		//
		// @since 3.6.0
		WorkspaceFolders bool `json:"workspaceFolders,omitempty"`

		// The client supports `workspace/configuration` requests.
		//
		// @since 3.6.0
		Configuration bool `json:"configuration,omitempty"`

		// Capabilities specific to the semantic token requests scoped to the
		// workspace.
		//
		// @since 3.16.0
		SemanticTokens *SemanticTokensWorkspaceClientCapabilities `json:"semanticTokens,omitempty"`

		// Capabilities specific to the code lens requests scoped to the
		// workspace.
		//
		// @since 3.16.0
		CodeLens *CodeLensWorkspaceClientCapabilities `json:"codeLens,omitempty"`

		// The client has support for file requests/notifications.
		//
		// @since 3.16.0
		FileOperations *struct {
			// Whether the client supports dynamic registration for file
			// requests/notifications.
			DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

			// The client has support for sending didCreateFiles notifications.
			DidCreate bool `json:"didCreate,omitempty"`

			// The client has support for sending willCreateFiles requests.
			WillCreate bool `json:"willCreate,omitempty"`

			// The client has support for sending didRenameFiles notifications.
			DidRename bool `json:"didRename,omitempty"`

			// The client has support for sending willRenameFiles requests.
			WillRename bool `json:"willRename,omitempty"`

			// The client has support for sending didDeleteFiles notifications.
			DidDelete bool `json:"didDelete,omitempty"`

			// The client has support for sending willDeleteFiles requests.
			WillDelete bool `json:"willDelete,omitempty"`
		} `json:"fileOperations,omitempty"`
	} `json:"workspace,omitempty"`

	// Text document specific client capabilities.
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`

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
		ShowMessage *ShowMessageRequestClientCapabilities `json:"showMessage,omitempty"`

		// Client capabilities for the show document request.
		//
		// @since 3.16.0
		ShowDocument *ShowDocumentClientCapabilities `json:"showDocument,omitempty"`
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
		RegularExpressions *RegularExpressionsClientCapabilities `json:"regularExpressions,omitempty"`

		// Client capabilities specific to the client's markdown parser.
		//
		// @since 3.16.0
		Markdown *MarkdownClientCapabilities `json:"markdown,omitempty"`
	} `json:"general,omitempty"`

	// Experimental client capabilities.
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

type WorkspaceEditClientCapabilities struct {

	// The client supports versioned document changes in `WorkspaceEdit`s
	DocumentChanges bool `json:"documentChanges,omitempty"`

	// The resource operations the client supports. Clients should at least
	// support 'create', 'rename' and 'delete' files and folders.
	//
	// @since 3.13.0
	ResourceOperations []ResourceOperationKind `json:"resourceOperations,omitempty"`

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
	NormalizesLineEndings bool `json:"normalizesLineEndings,omitempty"`

	// Whether the client in general supports change annotations on text edits,
	// create file, rename file and delete file changes.
	//
	// @since 3.16.0
	ChangeAnnotationSupport *struct {
		// Whether the client groups edits with equal labels into tree nodes,
		// for instance all edits labelled with "Changes in Strings" would
		// be a tree node.
		GroupsOnLabel bool `json:"groupsOnLabel,omitempty"`
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

type DidChangeConfigurationClientCapabilities struct {
	// Did change configuration notification supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DidChangeWatchedFilesClientCapabilities struct {
	// Did change watched files notification supports dynamic registration.
	// Please note that the current protocol doesn't support static
	// configuration for file changes from the server side.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type WorkspaceSymbolClientCapabilities struct {
	// Symbol request supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Specific capabilities for the `SymbolKind` in the `workspace/symbol`
	// request.
	SymbolKind *struct {
		// The symbol kind values the client supports. When this
		// property exists the client also guarantees that it will
		// handle values outside its set gracefully and falls back
		// to a default value when unknown.
		//
		// If this property is not present the client only supports
		// the symbol kinds from `File` to `Array` as defined in
		// the initial version of the protocol.
		ValueSet *[]SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`

	// The client supports tags on `SymbolInformation`.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// @since 3.16.0
	TagSupport *struct {
		// The tags supported by the client.
		ValueSet []SymbolTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
}

type SymbolKind int

const SymbolKindFile SymbolKind = 1
const SymbolKindModule SymbolKind = 2
const SymbolKindNamespace SymbolKind = 3
const SymbolKindPackage SymbolKind = 4
const SymbolKindClass SymbolKind = 5
const SymbolKindMethod SymbolKind = 6
const SymbolKindProperty SymbolKind = 7
const SymbolKindField SymbolKind = 8
const SymbolKindConstructor SymbolKind = 9
const SymbolKindEnum SymbolKind = 10
const SymbolKindInterface SymbolKind = 11
const SymbolKindFunction SymbolKind = 12
const SymbolKindVariable SymbolKind = 13
const SymbolKindConstant SymbolKind = 14
const SymbolKindString SymbolKind = 15
const SymbolKindNumber SymbolKind = 16
const SymbolKindBoolean SymbolKind = 17
const SymbolKindArray SymbolKind = 18
const SymbolKindObject SymbolKind = 19
const SymbolKindKey SymbolKind = 20
const SymbolKindNull SymbolKind = 21
const SymbolKindEnumMember SymbolKind = 22
const SymbolKindStruct SymbolKind = 23
const SymbolKindEvent SymbolKind = 24
const SymbolKindOperator SymbolKind = 25
const SymbolKindTypeParameter SymbolKind = 26

func (s SymbolKind) String() string {
	switch s {
	case SymbolKindFile:
		return "SymbolKind:File"
	case SymbolKindModule:
		return "SymbolKind:Module"
	case SymbolKindNamespace:
		return "SymbolKind:Namespace"
	case SymbolKindPackage:
		return "SymbolKind:Package"
	case SymbolKindClass:
		return "SymbolKind:Class"
	case SymbolKindMethod:
		return "SymbolKind:Method"
	case SymbolKindProperty:
		return "SymbolKind:Property"
	case SymbolKindField:
		return "SymbolKind:Field"
	case SymbolKindConstructor:
		return "SymbolKind:Constructor"
	case SymbolKindEnum:
		return "SymbolKind:Enum"
	case SymbolKindInterface:
		return "SymbolKind:Interface"
	case SymbolKindFunction:
		return "SymbolKind:Function"
	case SymbolKindVariable:
		return "SymbolKind:Variable"
	case SymbolKindConstant:
		return "SymbolKind:Constant"
	case SymbolKindString:
		return "SymbolKind:String"
	case SymbolKindNumber:
		return "SymbolKind:Number"
	case SymbolKindBoolean:
		return "SymbolKind:Boolean"
	case SymbolKindArray:
		return "SymbolKind:Array"
	case SymbolKindObject:
		return "SymbolKind:Object"
	case SymbolKindKey:
		return "SymbolKind:Key"
	case SymbolKindNull:
		return "SymbolKind:Null"
	case SymbolKindEnumMember:
		return "SymbolKind:EnumMember"
	case SymbolKindStruct:
		return "SymbolKind:Struct"
	case SymbolKindEvent:
		return "SymbolKind:Event"
	case SymbolKindOperator:
		return "SymbolKind:Operator"
	case SymbolKindTypeParameter:
		return "SymbolKind:TypeParameter"
	default:
		return fmt.Sprintf("SymbolKind:%d", s)
	}
}

type SymbolTag int

const SymbolTagDeprecated SymbolTag = 1

type ExecuteCommandClientCapabilities struct {
	// Execute command supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type SemanticTokensWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from
	// the server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// semantic tokens currently shown. It should be used with absolute care
	// and is useful for situation where a server for example detect a project
	// wide change that requires such a calculation.
	RefreshSupport bool `json:"refreshSupport,omitempty"`
}

type CodeLensWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from the
	// server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// code lenses currently shown. It should be used with absolute care and is
	// useful for situation where a server for example detect a project wide
	// change that requires such a calculation.
	RefreshSupport bool `json:"refreshSupport,omitempty"`
}

type TextDocumentClientCapabilities struct {
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`

	// Capabilities specific to the `textDocument/completion` request.
	Completion *CompletionClientCapabilities `json:"completion,omitempty"`

	// Capabilities specific to the `textDocument/hover` request.
	Hover *HoverClientCapabilities `json:"hover,omitempty"`

	// Capabilities specific to the `textDocument/signatureHelp` request.
	SignatureHelp *SignatureHelpClientCapabilities `json:"signatureHelp,omitempty"`

	// Capabilities specific to the `textDocument/declaration` request.
	//
	// @since 3.14.0
	Declaration *DeclarationClientCapabilities `json:"declaration,omitempty"`

	// Capabilities specific to the `textDocument/definition` request.
	Definition *DefinitionClientCapabilities `json:"definition,omitempty"`

	// Capabilities specific to the `textDocument/typeDefinition` request.
	//
	// @since 3.6.0
	TypeDefinition *TypeDefinitionClientCapabilities `json:"typeDefinition,omitempty"`

	// Capabilities specific to the `textDocument/implementation` request.
	//
	// @since 3.6.0
	Implementation *ImplementationClientCapabilities `json:"implementation,omitempty"`

	// Capabilities specific to the `textDocument/references` request.
	References *ReferenceClientCapabilities `json:"references,omitempty"`

	// Capabilities specific to the `textDocument/documentHighlight` request.
	DocumentHighlight *DocumentHighlightClientCapabilities `json:"documentHighlight,omitempty"`

	// Capabilities specific to the `textDocument/documentSymbol` request.
	DocumentSymbol *DocumentSymbolClientCapabilities `json:"documentSymbol,omitempty"`

	// Capabilities specific to the `textDocument/codeAction` request.
	CodeAction *CodeActionClientCapabilities `json:"codeAction,omitempty"`

	// Capabilities specific to the `textDocument/codeLens` request.
	CodeLens *CodeLensClientCapabilities `json:"codeLens,omitempty"`

	// Capabilities specific to the `textDocument/documentLink` request.
	DocumentLink *DocumentLinkClientCapabilities `json:"documentLink,omitempty"`

	// Capabilities specific to the `textDocument/documentColor` and the
	// `textDocument/colorPresentation` request.
	//
	// @since 3.6.0
	ColorProvider *DocumentColorClientCapabilities `json:"colorProvider,omitempty"`

	// Capabilities specific to the `textDocument/formatting` request.
	Formatting *DocumentFormattingClientCapabilities `json:"formatting,omitempty"`

	// Capabilities specific to the `textDocument/rangeFormatting` request.
	RangeFormatting *DocumentRangeFormattingClientCapabilities `json:"rangeFormatting,omitempty"`

	// Capabilities specific to the `textDocument/onTypeFormatting` request.
	OnTypeFormatting *DocumentOnTypeFormattingClientCapabilities `json:"onTypeFormatting,omitempty"`

	// Capabilities specific to the `textDocument/rename` request.
	Rename *RenameClientCapabilities `json:"rename,omitempty"`

	// Capabilities specific to the `textDocument/publishDiagnostics`
	// notification.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`

	// Capabilities specific to the `textDocument/foldingRange` request.
	//
	// @since 3.10.0
	FoldingRange *FoldingRangeClientCapabilities `json:"foldingRange,omitempty"`

	// Capabilities specific to the `textDocument/selectionRange` request.
	//
	// @since 3.15.0
	SelectionRange *SelectionRangeClientCapabilities `json:"selectionRange,omitempty"`

	// Capabilities specific to the `textDocument/linkedEditingRange` request.
	//
	// @since 3.16.0
	LinkedEditingRange *LinkedEditingRangeClientCapabilities `json:"linkedEditingRange,omitempty"`

	// Capabilities specific to the various call hierarchy requests.
	//
	// @since 3.16.0
	CallHierarchy *CallHierarchyClientCapabilities `json:"callHierarchy,omitempty"`

	// Capabilities specific to the various semantic token requests.
	//
	// @since 3.16.0
	SemanticTokens *SemanticTokensClientCapabilities `json:"semanticTokens,omitempty"`

	// Capabilities specific to the `textDocument/moniker` request.
	//
	// @since 3.16.0
	Moniker *MonikerClientCapabilities `json:"moniker,omitempty"`
}

type TextDocumentSyncClientCapabilities struct {
	// Whether text document synchronization supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports sending will save notifications.
	WillSave bool `json:"willSave,omitempty"`

	// The client supports sending a will save request and
	// waits for a response providing text edits which will
	// be applied to the document before it is saved.
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`

	// The client supports did save notifications.
	DidSave bool `json:"didSave,omitempty"`
}

type CompletionKindCapabilities struct {
	// Whether completion supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports the following `CompletionItem` specific
	// capabilities.
	CompletionItem *struct {
		// Client supports snippets as insert text.
		//
		// A snippet can define tab stops and placeholders with `$1`, `$2`
		// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
		// the end of the snippet. Placeholders with equal identifiers are
		// linked, that is typing in one will update others too.
		SnippetSupport bool `json:"snippetSupport,omitempty"`

		// Client supports commit characters on a completion item.
		CommitCharactersSupport bool `json:"commitCharactersSupport,omitempty"`

		// Client supports the follow content formats for the documentation
		// property. The order describes the preferred format of the client.
		DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

		// Client supports the deprecated property on a completion item.
		DeprecatedSupport bool `json:"deprecatedSupport,omitempty"`

		// Client supports the preselect property on a completion item.
		PreselectSupport bool `json:"preselectSupport,omitempty"`

		// Client supports the tag property on a completion item. Clients
		// supporting tags have to handle unknown tags gracefully. Clients
		// especially need to preserve unknown tags when sending a completion
		// item back to the server in a resolve call.
		//
		// @since 3.15.0
		TagSupport *struct {
			// The tags supported by the client.
			ValueSet []CompletionItemTag `json:"valueSet"`
		} `json:"tagSupport,omitempty"`

		// Client supports insert replace edit to control different behavior if
		// a completion item is inserted in the text or should replace text.
		//
		// @since 3.16.0
		InsertReplaceSupport bool `json:"insertReplaceSupport,omitempty"`

		// Indicates which properties a client can resolve lazily on a
		// completion item. Before version 3.16.0 only the predefined properties
		// `documentation` and `detail` could be resolved lazily.
		//
		// @since 3.16.0
		ResolveSupport *struct {
			// The properties that a client can resolve lazily.
			Properties []string `json:"properties"`
		} `json:"resolveSupport,omitempty"`

		// The client supports the `insertTextMode` property on
		// a completion item to override the whitespace handling mode
		// as defined by the client (see `insertTextMode`).
		//
		// @since 3.16.0
		InsertTextModeSupport *struct {
			ValueSet []InsertTextMode `json:"valueSet"`
		} `json:"insertTextModeSupport,omitempty"`

		// The client has support for completion item label
		// details (see also `CompletionItemLabelDetails`).
		//
		// @since 3.17.0 - proposed state
		LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
	} `json:"completionItem,omitempty"`

	CompletionItemKind *struct {
		// The completion item kind values the client supports. When this
		// property exists the client also guarantees that it will
		// handle values outside its set gracefully and falls back
		// to a default value when unknown.
		//
		// If this property is not present the client only supports
		// the completion items kinds from `Text` to `Reference` as defined in
		// the initial version of the protocol.
		ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
	} `json:"completionItemKind,omitempty"`

	// The client supports to send additional context information for a
	// `textDocument/completion` request.
	ContextSupport bool `json:"contextSupport,omitempty"`

	// The client's default when the completion item doesn't provide a
	// `insertTextMode` property.
	//
	// @since 3.17.0
	InsertTextMode *InsertTextMode `json:"insertTextMode,omitempty"`
}

type MarkupKind string

// Plain text is supported as a content format
const MarkupKindPlainText MarkupKind = "plaintext"

// Markdown is supported as a content format
const MarkupKindMarkdown MarkupKind = "markdown"

// Completion item tags are extra annotations that tweak the rendering of a
// completion item.
//
// @since 3.15.0
type CompletionItemTag int

// Render a completion as obsolete, usually using a strike-out.
const CompletionItemTagDeprecated CompletionItemTag = 1

// How whitespace and indentation is handled during completion
// item insertion.
//
// @since 3.16.0
type InsertTextMode int

// The insertion or replace strings is taken as it is. If the
// value is multi line the lines below the cursor will be
// inserted using the indentation defined in the string value.
// The client will not apply any kind of adjustments to the
// string.
const InsertTextModeAsIs InsertTextMode = 1

// The editor adjusts leading whitespace of new lines so that
// they match the indentation up to the cursor of the line for
// which the item is accepted.
//
// Consider a line like this: <2tabs><cursor><3tabs>foo. Accepting a
// multi line completion item is indented using 2 tabs and all
// following lines inserted will be indented using 2 tabs as well.
const InsertTextModeAdjustIndentation InsertTextMode = 2

// The kind of a completion entry.
type CompletionItemKind int

const CompletionItemKindText CompletionItemKind = 1
const CompletionItemKindMethod CompletionItemKind = 2
const CompletionItemKindFunction CompletionItemKind = 3
const CompletionItemKindConstructor CompletionItemKind = 4
const CompletionItemKindField CompletionItemKind = 5
const CompletionItemKindVariable CompletionItemKind = 6
const CompletionItemKindClass CompletionItemKind = 7
const CompletionItemKindInterface CompletionItemKind = 8
const CompletionItemKindModule CompletionItemKind = 9
const CompletionItemKindProperty CompletionItemKind = 10
const CompletionItemKindUnit CompletionItemKind = 11
const CompletionItemKindValue CompletionItemKind = 12
const CompletionItemKindEnum CompletionItemKind = 13
const CompletionItemKindKeyword CompletionItemKind = 14
const CompletionItemKindSnippet CompletionItemKind = 15
const CompletionItemKindColor CompletionItemKind = 16
const CompletionItemKindFile CompletionItemKind = 17
const CompletionItemKindReference CompletionItemKind = 18
const CompletionItemKindFolder CompletionItemKind = 19
const CompletionItemKindEnumMember CompletionItemKind = 20
const CompletionItemKindConstant CompletionItemKind = 21
const CompletionItemKindStruct CompletionItemKind = 22
const CompletionItemKindEvent CompletionItemKind = 23
const CompletionItemKindOperator CompletionItemKind = 24
const CompletionItemKindTypeParameter CompletionItemKind = 25

type CompletionClientCapabilities struct {
	// Whether completion supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports the following `CompletionItem` specific
	// capabilities.
	CompletionItem *struct {
		// Client supports snippets as insert text.
		//
		// A snippet can define tab stops and placeholders with `$1`, `$2`
		// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
		// the end of the snippet. Placeholders with equal identifiers are
		// linked, that is typing in one will update others too.
		SnippetSupport bool `json:"snippetSupport,omitempty"`

		// Client supports commit characters on a completion item.
		CommitCharactersSupport bool `json:"commitCharactersSupport,omitempty"`

		// Client supports the follow content formats for the documentation
		// property. The order describes the preferred format of the client.
		DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

		// Client supports the deprecated property on a completion item.
		DeprecatedSupport bool `json:"deprecatedSupport,omitempty"`

		// Client supports the preselect property on a completion item.
		PreselectSupport bool `json:"preselectSupport,omitempty"`

		// Client supports the tag property on a completion item. Clients
		// supporting tags have to handle unknown tags gracefully. Clients
		// especially need to preserve unknown tags when sending a completion
		// item back to the server in a resolve call.
		//
		// @since 3.15.0
		TagSupport *struct {
			// The tags supported by the client.
			ValueSet []CompletionItemTag `json:"valueSet"`
		} `json:"tagSupport,omitempty"`

		// Client supports insert replace edit to control different behavior if
		// a completion item is inserted in the text or should replace text.
		//
		// @since 3.16.0
		InsertReplaceSupport bool `json:"insertReplaceSupport,omitempty"`

		// Indicates which properties a client can resolve lazily on a
		// completion item. Before version 3.16.0 only the predefined properties
		// `documentation` and `detail` could be resolved lazily.
		//
		// @since 3.16.0
		ResolveSupport *struct {
			// The properties that a client can resolve lazily.
			Properties []string `json:"properties"`
		} `json:"resolveSupport,omitempty"`

		// The client supports the `insertTextMode` property on
		// a completion item to override the whitespace handling mode
		// as defined by the client (see `insertTextMode`).
		//
		// @since 3.16.0
		InsertTextModeSupport *struct {
			ValueSet []InsertTextMode `json:"valueSet"`
		} `json:"insertTextModeSupport,omitempty"`

		// The client has support for completion item label
		// details (see also `CompletionItemLabelDetails`).
		//
		// @since 3.17.0 - proposed state
		LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
	} `json:"completionItem,omitempty"`

	CompletionItemKind *struct {
		// The completion item kind values the client supports. When this
		// property exists the client also guarantees that it will
		// handle values outside its set gracefully and falls back
		// to a default value when unknown.
		//
		// If this property is not present the client only supports
		// the completion items kinds from `Text` to `Reference` as defined in
		// the initial version of the protocol.
		ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
	} `json:"completionItemKind,omitempty"`

	// The client supports to send additional context information for a
	// `textDocument/completion` request.
	ContextSupport bool `json:"contextSupport,omitempty"`

	// The client's default when the completion item doesn't provide a
	// `insertTextMode` property.
	//
	// @since 3.17.0
	InsertTextMode *InsertTextMode `json:"insertTextMode,omitempty"`
}

type HoverClientCapabilities struct {
	// Whether hover supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Client supports the follow content formats if the content
	// property refers to a `literal of type MarkupContent`.
	// The order describes the preferred format of the client.
	ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
}

type SignatureHelpClientCapabilities struct {
	// Whether signature help supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports the following `SignatureInformation`
	// specific properties.
	SignatureInformation *struct {
		// Client supports the follow content formats for the documentation
		// property. The order describes the preferred format of the client.
		DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

		// Client capabilities specific to parameter information.
		ParameterInformation *struct {
			// The client supports processing label offsets instead of a
			// simple label string.
			//
			// @since 3.14.0
			LabelOffsetSupport bool `json:"labelOffsetSupport,omitempty"`
		} `json:"parameterInformation,omitempty"`

		// The client supports the `activeParameter` property on
		// `SignatureInformation` literal.
		//
		// @since 3.16.0
		ActiveParameterSupport bool `json:"activeParameterSupport,omitempty"`
	} `json:"signatureInformation,omitempty"`

	// The client supports to send additional context information for a
	// `textDocument/signatureHelp` request. A client that opts into
	// contextSupport will also support the `retriggerCharacters` on
	// `SignatureHelpOptions`.
	//
	// @since 3.15.0
	ContextSupport bool `json:"contextSupport,omitempty"`
}

type DeclarationClientCapabilities struct {
	// Whether declaration supports dynamic registration. If this is set to
	// `true` the client supports the new `DeclarationRegistrationOptions`
	// return value for the corresponding server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports additional metadata in the form of declaration links.
	LinkSupport bool `json:"linkSupport,omitempty"`
}

type DefinitionClientCapabilities struct {
	// Whether definition supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports additional metadata in the form of definition links.
	//
	// @since 3.14.0
	LinkSupport bool `json:"linkSupport,omitempty"`
}

type TypeDefinitionClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to
	// `true` the client supports the new `TypeDefinitionRegistrationOptions`
	// return value for the corresponding server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports additional metadata in the form of definition links.
	//
	// @since 3.14.0
	LinkSupport bool `json:"linkSupport,omitempty"`
}

type ImplementationClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to
	// `true` the client supports the new `ImplementationRegistrationOptions`
	// return value for the corresponding server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports additional metadata in the form of definition links.
	//
	// @since 3.14.0
	LinkSupport bool `json:"linkSupport,omitempty"`
}

type ReferenceClientCapabilities struct {
	// Whether references supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentHighlightClientCapabilities struct {
	// Whether document highlight supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentSymbolClientCapabilities struct {
	// Whether document symbol supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Specific capabilities for the `SymbolKind` in the
	// `textDocument/documentSymbol` request.
	SymbolKind *struct {
		// The symbol kind values the client supports. When this
		// property exists the client also guarantees that it will
		// handle values outside its set gracefully and falls back
		// to a default value when unknown.
		//
		// If this property is not present the client only supports
		// the symbol kinds from `File` to `Array` as defined in
		// the initial version of the protocol.
		ValueSet []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`

	// The client supports hierarchical document symbols.
	HierarchicalDocumentSymbolSupport bool `json:"hierarchicalDocumentSymbolSupport,omitempty"`

	// The client supports tags on `SymbolInformation`. Tags are supported on
	// `DocumentSymbol` if `hierarchicalDocumentSymbolSupport` is set to true.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// @since 3.16.0
	TagSupport *struct {
		// The tags supported by the client.
		ValueSet []SymbolTag `json:"valueSet"`
	} `json:"tagSupport,omitempty"`

	// The client supports an additional label presented in the UI when
	// registering a document symbol provider.
	//
	// @since 3.16.0
	LabelSupport bool `json:"labelSupport,omitempty"`
}

type CodeActionClientCapabilities struct {
	// Whether code action supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports code action literals as a valid
	// response of the `textDocument/codeAction` request.
	//
	// @since 3.8.0
	CodeActionLiteralSupport *struct {
		// The code action kind is supported with the following value
		// set.
		CodeActionKind struct {
			// The code action kind values the client supports. When this
			// property exists the client also guarantees that it will
			// handle values outside its set gracefully and falls back
			// to a default value when unknown.
			ValueSet []CodeActionKind `json:"valueSet"`
		} `json:"codeActionKind"`
	} `json:"codeActionLiteralSupport,omitempty"`

	// Whether code action supports the `isPreferred` property.
	//
	// @since 3.15.0
	IsPreferredSupport bool `json:"isPreferredSupport,omitempty"`

	// Whether code action supports the `disabled` property.
	//
	// @since 3.16.0
	DisabledSupport bool `json:"disabledSupport,omitempty"`

	// Whether code action supports the `data` property which is
	// preserved between a `textDocument/codeAction` and a
	// `codeAction/resolve` request.
	//
	// @since 3.16.0
	DataSupport bool `json:"dataSupport,omitempty"`

	// Whether the client supports resolving additional code action
	// properties via a separate `codeAction/resolve` request.
	//
	// @since 3.16.0
	ResolveSupport *struct {

		// The properties that a client can resolve lazily.
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`

	// Whether the client honors the change annotations in
	// text edits and resource operations returned via the
	// `CodeAction#edit` property by for example presenting
	// the workspace edit in the user interface and asking
	// for confirmation.
	//
	// @since 3.16.0
	HonorsChangeAnnotations bool `json:"honorsChangeAnnotations,omitempty"`
}

type CodeLensClientCapabilities struct {
	// Whether code lens supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentLinkClientCapabilities struct {
	// Whether document link supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Whether the client supports the `tooltip` property on `DocumentLink`.
	//
	// @since 3.15.0
	TooltipSupport bool `json:"tooltipSupport,omitempty"`
}

type DocumentColorClientCapabilities struct {
	// Whether document color supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentFormattingClientCapabilities struct {
	// Whether formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentRangeFormattingClientCapabilities struct {
	// Whether formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type DocumentOnTypeFormattingClientCapabilities struct {
	// Whether on type formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type RenameClientCapabilities struct {
	// Whether rename supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Client supports testing for validity of rename operations
	// before execution.
	//
	// @since version 3.12.0
	PrepareSupport bool `json:"prepareSupport,omitempty"`

	// Client supports the default behavior result
	// (`{ defaultBehavior: boolean }`).
	//
	// The value indicates the default behavior used by the
	// client.
	//
	// @since version 3.16.0
	PrepareSupportDefaultBehavior *PrepareSupportDefaultBehavior `json:"prepareSupportDefaultBehavior,omitempty"`

	// Whether th client honors the change annotations in
	// text edits and resource operations returned via the
	// rename request's workspace edit by for example presenting
	// the workspace edit in the user interface and asking
	// for confirmation.
	//
	// @since 3.16.0
	honorsChangeAnnotations bool `json:"honorsChangeAnnotations,omitempty"`
}

type PublishDiagnosticsClientCapabilities struct {
	// Whether the clients accepts diagnostics with related information.
	RelatedInformation bool `json:"relatedInformation,omitempty"`

	// Client supports the tag property to provide meta data about a diagnostic.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// @since 3.15.0
	TagSupport *struct {
		// The tags supported by the client.
		ValueSet []DiagnosticTag `json:"valueSet"`
	} `json:"tagSupport,omitempty"`

	// Whether the client interprets the version property of the
	// `textDocument/publishDiagnostics` notification's parameter.
	//
	// @since 3.15.0
	VersionSupport bool `json:"versionSupport,omitempty"`

	// Client supports a codeDescription property
	//
	// @since 3.16.0
	CodeDescriptionSupport bool `json:"codeDescriptionSupport,omitempty"`

	// Whether code action supports the `data` property which is
	// preserved between a `textDocument/publishDiagnostics` and
	// `textDocument/codeAction` request.
	//
	// @since 3.16.0
	DataSupport bool `json:"dataSupport,omitempty"`
}

type FoldingRangeClientCapabilities struct {
	// Whether implementation supports dynamic registration for folding range
	// providers. If this is set to `true` the client supports the new
	// `FoldingRangeRegistrationOptions` return value for the corresponding
	// server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The maximum number of folding ranges that the client prefers to receive
	// per document. The value serves as a hint, servers are free to follow the
	// limit.
	RangeLimit *int `json:"rangeLimit,omitempty"`

	// If set, the client signals that it only supports folding complete lines.
	// If set, client will ignore specified `startCharacter` and `endCharacter`
	// properties in a FoldingRange.
	LineFoldingOnly bool `json:"lineFoldingOnly,omitempty"`
}

type SelectionRangeClientCapabilities struct {
	// Whether implementation supports dynamic registration for selection range
	// providers. If this is set to `true` the client supports the new
	// `SelectionRangeRegistrationOptions` return value for the corresponding
	// server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type LinkedEditingRangeClientCapabilities struct {
	// Whether implementation supports dynamic registration.
	// If this is set to `true` the client supports the new
	// `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type CallHierarchyClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to
	// `true` the client supports the new `(TextDocumentRegistrationOptions &
	// StaticRegistrationOptions)` return value for the corresponding server
	// capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type SemanticTokensClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to
	// `true` the client supports the new `(TextDocumentRegistrationOptions &
	// StaticRegistrationOptions)` return value for the corresponding server
	// capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Which requests the client supports and might send to the server
	// depending on the server's capability. Please note that clients might not
	// show semantic tokens or degrade some of the user experience if a range
	// or full request is advertised by the client but not provided by the
	// server. If for example the client capability `requests.full` and
	// `request.range` are both set to true but the server only provides a
	// range provider the client might not render a minimap correctly or might
	// even decide to not show any semantic tokens at all.
	Requests struct {

		// The client will send the `textDocument/semanticTokens/range` request
		// if the server provides a corresponding handler.
		Range bool `json:"range,omitempty"`

		// The client will send the `textDocument/semanticTokens/full` request
		// if the server provides a corresponding handler.
		Full *struct {
			// The client will send the `textDocument/semanticTokens/full/delta`
			// request if the server provides a corresponding handler.
			Delta bool `json:"delta,omitempty"`
		} `json:"full,omitempty"`
	} `json:"requests"`

	// The token types that the client supports.
	TokenTypes []string `json:"tokenTypes"`

	// The token modifiers that the client supports.
	TokenModifiers []string `json:"tokenModifiers"`

	// The formats the clients supports.
	Formats []TokenFormat `json:"formats"`

	// Whether the client supports tokens that can overlap each other.
	OverlappingTokenSupport bool `json:"overlappingTokenSupport,omitempty"`

	// Whether the client supports tokens that can span multiple lines.
	MultilineTokenSupport bool `json:"multilineTokenSupport,omitempty"`
}

type MonikerClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to
	// `true` the client supports the new `(TextDocumentRegistrationOptions &
	// StaticRegistrationOptions)` return value for the corresponding server
	// capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// The kind of a code action.
//
// Kinds are a hierarchical list of identifiers separated by `.`,
// e.g. `"refactor.extract.function"`.
//
// The set of kinds is open and client needs to announce the kinds it supports
// to the server during initialization.
type CodeActionKind string

// Empty kind.
const CodeActionKindEmpty CodeActionKind = ""

// Base kind for quickfix actions: "quickfix".
const CodeActionKindQuickFix CodeActionKind = "quickfix"

// Base kind for refactoring actions: "refactor".
const CodeActionKindRefactor CodeActionKind = "refactor"

// Base kind for refactoring extraction actions: "refactor.extract".
//
// Example extract actions:
//
// - Extract method
// - Extract function
// - Extract variable
// - Extract interface from class
// - ...
const CodeActionKindRefactorExtract CodeActionKind = "refactor.extract"

// Base kind for refactoring inline actions: "refactor.inline".
//
// Example inline actions:
//
// - Inline function
// - Inline variable
// - Inline constant
// - ...
const CodeActionKindRefactorInline CodeActionKind = "refactor.inline"

// Base kind for refactoring rewrite actions: "refactor.rewrite".
//
// Example rewrite actions:
//
// - Convert JavaScript function to class
// - Add or remove parameter
// - Encapsulate field
// - Make method static
// - Move method to base class
// - ...
const CodeActionKindRefactorRewrite CodeActionKind = "refactor.rewrite"

// Base kind for source actions: `source`.
//
// Source code actions apply to the entire file.
const CodeActionKindSource CodeActionKind = "source"

// Base kind for an organize imports source action:
// `source.organizeImports`.
const CodeActionKindSourceOrganizeImports CodeActionKind = "source.organizeImports"

// Base kind for a "fix all" source action: `source.fixAll`.
//
// "Fix all" actions automatically fix errors that have a clear fix that
// do not require user input. They should not suppress errors or perform
// unsafe fixes such as generating new types or classes.
//
// @since 3.17.0
const CodeActionKindSourceFixAll CodeActionKind = "source.fixAll"

type PrepareSupportDefaultBehavior int

// The client's default behavior is to select the identifier
// according the to language's syntax rule.
const PrepareSupportDefaultBehaviorIdentifier = 1

type TokenFormat string

const TokenFormatRelative TokenFormat = "relative"

// Show message request client capabilities
type ShowMessageRequestClientCapabilities struct {
	// Capabilities specific to the `MessageActionItem` type.
	MessageActionItem struct {

		// Whether the client supports additional attributes which
		// are preserved and sent back to the server in the
		// request's response.
		AdditionalPropertiesSupport bool `json:"additionalPropertiesSupport,omitempty"`
	} `json:"messageActionItem"`
}

// Client capabilities for the show document request.
//
// @since 3.16.0
type ShowDocumentClientCapabilities struct {
	// The client has support for the show document
	// request.
	Support bool `json:"support"`
}

// Client capabilities specific to regular expressions.
type RegularExpressionsClientCapabilities struct {
	// The engine's name.
	Engine string `json:"engine,required"`

	// The engine's version.
	Version string `json:"version,omitempty"`
}

// Client capabilities specific to the used markdown parser.
//
// @since 3.16.0
type MarkdownClientCapabilities struct {
	// The name of the parser.
	Parser string `json:"parser,required"`

	// The version of the parser.
	Version string `json:"version,omitempty"`
}
