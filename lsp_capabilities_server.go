//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"

	"go.bug.st/json"
)

type ServerCapabilities struct {
	// Defines how text documents are synced. Is either a detailed structure
	// defining each notification or for backwards compatibility the
	// TextDocumentSyncKind number. If omitted it defaults to
	// `TextDocumentSyncKind.None`.
	TextDocumentSync *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`

	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides hover support.
	HoverProvider *HoverOptions `json:"hoverProvider,omitempty"`

	// The server provides signature help support.
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`

	// The server provides go to declaration support.
	//
	// @since 3.14.0
	DeclarationProvider *DeclarationOptions `json:"declarationProvider,omitempty"`

	// The server provides goto definition support.
	DefinitionProvider *DefinitionOptions `json:"definitionProvider,omitempty"`

	// The server provides goto type definition support.
	//
	// @since 3.6.0
	TypeDefinitionProvider *TypeDefinitionOptions `json:"typeDefinitionProvider,omitempty"`

	// The server provides goto implementation support.
	//
	// @since 3.6.0
	ImplementationProvider *ImplementationOptions `json:"implementationProvider,omitempty"`

	// The server provides find references support.
	ReferencesProvider *ReferenceOptions `json:"referencesProvider,omitempty"`

	// The server provides document highlight support.
	DocumentHighlightProvider *DocumentHighlightOptions `json:"documentHighlightProvider,omitempty"`

	// The server provides document symbol support.
	DocumentSymbolProvider *DocumentSymbolOptions `json:"documentSymbolProvider,omitempty"`

	// The server provides code actions. The `CodeActionOptions` return type is
	// only valid if the client signals code action literal support via the
	// property `textDocument.codeAction.codeActionLiteralSupport`.
	CodeActionProvider *CodeActionOptions `json:"codeActionProvider,omitempty"`

	// The server provides code lens.
	CodeLensProvider *CodeLensOptions `json:"codeLensProvider,omitempty"`

	// The server provides document link support.
	DocumentLinkProvider *DocumentLinkOptions `json:"documentLinkProvider,omitempty"`

	// The server provides color provider support.
	//
	// @since 3.6.0
	ColorProvider *DocumentColorOptions `json:"colorProvider,omitempty"`

	// The server provides document formatting.
	DocumentFormattingProvider *DocumentFormattingOptions `json:"documentFormattingProvider,omitempty"`

	// The server provides document range formatting.
	DocumentRangeFormattingProvider *DocumentRangeFormattingOptions `json:"documentRangeFormattingProvider,omitempty"`

	// The server provides document formatting on typing.
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`

	// The server provides rename support. RenameOptions may only be
	// specified if the client states that it supports
	// `prepareSupport` in its initial `initialize` request.
	RenameProvider *RenameOptions `json:"renameProvider,omitempty"`

	// The server provides folding provider support.
	//
	// @since 3.10.0
	FoldingRangeProvider *FoldingRangeOptions `json:"foldingRangeProvider,omitempty"`

	// The server provides execute command support.
	ExecuteCommandProvider *ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`

	// The server provides selection range support.
	//
	// @since 3.15.0
	SelectionRangeProvider *SelectionRangeOptions `json:"selectionRangeProvider,omitempty"`

	// The server provides linked editing range support.
	//
	// @since 3.16.0
	LinkedEditingRangeProvider *LinkedEditingRangeOptions `json:"linkedEditingRangeProvider,omitempty"`

	// The server provides call hierarchy support.
	//
	// @since 3.16.0
	CallHierarchyProvider *CallHierarchyOptions `json:"callHierarchyProvider,omitempty"`

	// The server provides semantic tokens support.
	//
	// @since 3.16.0
	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`

	// Whether server provides moniker support.
	//
	// @since 3.16.0
	MonikerProvider *MonikerOptions `json:"monikerProvider,omitempty"`

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider *WorkspaceSymbolOptions `json:"workspaceSymbolProvider,omitempty"`

	// Workspace specific server capabilities
	Workspace *struct {
		// The server supports workspace folder.
		//
		// @since 3.6.0
		WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`

		// The server is interested in file notifications/requests.
		//
		// @since 3.16.0
		FileOperations *struct {
			// The server is interested in receiving didCreateFiles
			// notifications.
			DidCreate *FileOperationRegistrationOptions `json:"didCreate,omitempty"`

			// The server is interested in receiving willCreateFiles requests.
			WillCreate *FileOperationRegistrationOptions `json:"willCreate,omitempty"`

			// The server is interested in receiving didRenameFiles
			// notifications.
			DidRename *FileOperationRegistrationOptions `json:"didRename,omitempty"`

			// The server is interested in receiving willRenameFiles requests.
			WillRename *FileOperationRegistrationOptions `json:"willRename,omitempty"`

			// The server is interested in receiving didDeleteFiles file
			// notifications.
			DidDelete *FileOperationRegistrationOptions `json:"didDelete,omitempty"`

			// The server is interested in receiving willDeleteFiles file
			// requests.
			WillDelete *FileOperationRegistrationOptions `json:"willDelete,omitempty"`
		} `json:"fileOperations,omitempty"`
	} `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

type TextDocumentSyncKind int

const TextDocumentSyncKindNone TextDocumentSyncKind = 0
const TextDocumentSyncKindFull TextDocumentSyncKind = 1
const TextDocumentSyncKindIncremental TextDocumentSyncKind = 2

type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server. If omitted open
	// close notification should not be sent.
	OpenClose bool `json:"openClose,omitempty"`

	// Change notifications are sent to the server. See
	// TextDocumentSyncKind.None, TextDocumentSyncKind.Full and
	// TextDocumentSyncKind.Incremental. If omitted it defaults to
	// TextDocumentSyncKind.None.
	Change TextDocumentSyncKind `json:"change,omitempty"`

	// If present will save notifications are sent to the server. If omitted
	// the notification should not be sent.
	WillSave bool `json:"willSave,omitempty"`

	// If present will save wait until requests are sent to the server. If
	// omitted the request should not be sent.
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`

	// If present save notifications are sent to the server. If omitted the
	// notification should not be sent.
	Save *SaveOptions `json:"save,omitempty"`
}

type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText bool `json:"includeText,omitempty"`
}

func (s *SaveOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = SaveOptions{}
		}
		return nil
	}

	type __ SaveOptions
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = SaveOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or SaveOptions")
}

type CompletionOptions struct {
	WorkDoneProgressOptions

	// Most tools trigger completion request automatically without explicitly
	// requesting it using a keyboard shortcut (e.g. Ctrl+Space). Typically they
	// do so when the user starts to type an identifier. For example if the user
	// types `c` in a JavaScript file code complete will automatically pop up
	// present `console` besides others as a completion item. Characters that
	// make up identifiers don't need to be listed here.
	//
	// If code complete should automatically be trigger on characters not being
	// valid inside an identifier (for example `.` in JavaScript) list them in
	// `triggerCharacters`.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`

	// The list of all possible characters that commit a completion. This field
	// can be used if clients don't support individual commit characters per
	// completion item. See client capability
	// `completion.completionItem.commitCharactersSupport`.
	//
	// If a server provides both `allCommitCharacters` and commit characters on
	// an individual completion item the ones on the completion item win.
	//
	// @since 3.2.0
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`

	// The server provides support to resolve additional
	// information for a completion item.
	ResolveProvider bool `json:"resolveProvider,omitempty"`

	// The server supports the following `CompletionItem` specific
	// capabilities.
	//
	// @since 3.17.0 - proposed state
	CompletionItem *CompletionItemOptions `json:"completionItem,omitempty"`
}

type CompletionItemOptions struct {
	// The server has support for completion item label
	// details (see also `CompletionItemLabelDetails`) when receiving
	// a completion item in a resolve call.
	//
	// @since 3.17.0 - proposed state
	LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
}

type HoverOptions struct {
	*WorkDoneProgressOptions
}

func (s *HoverOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = HoverOptions{}
		}
		return nil
	}

	type __ HoverOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = HoverOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or HoverOptions")
}

type SignatureHelpOptions struct {
	WorkDoneProgressOptions

	// The characters that trigger signature help
	// automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`

	// List of characters that re-trigger signature help.
	//
	// These trigger characters are only active when signature help is already
	// showing. All trigger characters are also counted as re-trigger
	// characters.
	//
	// @since 3.15.0
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// boolean|DeclarationOptions|DeclarationRegistrationOptions
type DeclarationOptions struct {
	*WorkDoneProgressOptions
	*StaticRegistrationOptions
	*TextDocumentRegistrationOptions
}

func (s *DeclarationOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DeclarationOptions{}
		}
		return nil
	}

	type __ DeclarationOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DeclarationOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DeclarationOptions")
}

// General text document registration options.
type TextDocumentRegistrationOptions struct {
	// A document selector to identify the scope of the registration. If set to
	// null the document selector provided on the client side will be used.
	DocumentSelector *DocumentSelector `json:"documentSelector,omitempty"`
}

type DocumentSelector []DocumentFilter

type DocumentFilter struct {
	// A language id, like `typescript`.
	Language string `json:"language,omitempty"`

	// A Uri [scheme](#Uri.scheme), like `file` or `untitled`.
	Scheme string `json:"scheme,omitempty"`

	// A glob pattern, like `*.{ts,js}`.
	//
	// Glob patterns can have the following syntax:
	// - `*` to match one or more characters in a path segment
	// - `?` to match on one character in a path segment
	// - `**` to match any number of path segments, including none
	// - `{}` to group sub patterns into an OR expression. (e.g. `**​/*.{ts,js}`
	//   matches all TypeScript and JavaScript files)
	// - `[]` to declare a range of characters to match in a path segment
	//   (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	// - `[!...]` to negate a range of characters to match in a path segment
	//   (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but
	//   not `example.0`)
	Pattern string `json:"pattern,omitempty"`
}

// Static registration options to be returned in the initialize request.
type StaticRegistrationOptions struct {
	// The id used to register the request. The id can be used to deregister
	// the request again. See also Registration#id.
	ID string `json:"id,omitempty"`
}

//boolean|DefinitionOptions
type DefinitionOptions struct {
	*WorkDoneProgressOptions
}

func (s *DefinitionOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DefinitionOptions{}
		}
		return nil
	}

	type __ DefinitionOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DefinitionOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DefinitionOptions")
}

// boolean | TypeDefinitionOptions | TypeDefinitionRegistrationOptions
type TypeDefinitionOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *TypeDefinitionOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = TypeDefinitionOptions{}
		}
		return nil
	}

	type __ TypeDefinitionOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = TypeDefinitionOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or TypeDefinitionOptions")
}

// boolean | ImplementationOptions | ImplementationRegistrationOptions
type ImplementationOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *ImplementationOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = ImplementationOptions{}
		}
		return nil
	}

	type __ ImplementationOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = ImplementationOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or ImplementationOptions")
}

// boolean | ReferenceOptions
type ReferenceOptions struct {
	*WorkDoneProgressOptions
}

func (s *ReferenceOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = ReferenceOptions{}
		}
		return nil
	}

	type __ ReferenceOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = ReferenceOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or ReferenceOptions")
}

// boolean | DocumentHighlightOptions
type DocumentHighlightOptions struct {
	*WorkDoneProgressOptions
}

func (s *DocumentHighlightOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DocumentHighlightOptions{}
		}
		return nil
	}

	type __ DocumentHighlightOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DocumentHighlightOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DocumentHighlightOptions")
}

// boolean | DocumentSymbolOptions
type DocumentSymbolOptions struct {
	*WorkDoneProgressOptions

	// A human-readable string that is shown when multiple outlines trees
	// are shown for the same document.
	//
	// @since 3.16.0
	Label string `json:"label,omitempty"`
}

func (s *DocumentSymbolOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DocumentSymbolOptions{}
		}
		return nil
	}

	type __ DocumentSymbolOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DocumentSymbolOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DocumentSymbolOptions")
}

// boolean | CodeActionOptions
type CodeActionOptions struct {
	*WorkDoneProgressOptions

	// CodeActionKinds that this server may return.
	//
	// The list of kinds may be generic, such as `CodeActionKind.Refactor`,
	// or the server may list out every specific kind they provide.
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`

	// The server provides support to resolve additional
	// information for a code action.
	//
	// @since 3.16.0
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

func (s *CodeActionOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = CodeActionOptions{}
		}
		return nil
	}

	type __ CodeActionOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = CodeActionOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or CodeActionOptions")
}

type CodeLensOptions struct {
	*WorkDoneProgressOptions

	// Code lens has a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type DocumentLinkOptions struct {
	*WorkDoneProgressOptions

	// Document links have a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// boolean | DocumentColorOptions | DocumentColorRegistrationOptions
type DocumentColorOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *DocumentColorOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DocumentColorOptions{}
		}
		return nil
	}

	type __ DocumentColorOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DocumentColorOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DocumentColorOptions")
}

// boolean | DocumentFormattingOptions
type DocumentFormattingOptions struct {
	*WorkDoneProgressOptions
}

func (s *DocumentFormattingOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DocumentFormattingOptions{}
		}
		return nil
	}

	type __ DocumentFormattingOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DocumentFormattingOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DocumentFormattingOptions")
}

// boolean | DocumentFormattingOptions
type DocumentRangeFormattingOptions struct {
	*WorkDoneProgressOptions
}

func (s *DocumentRangeFormattingOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = DocumentRangeFormattingOptions{}
		}
		return nil
	}

	type __ DocumentRangeFormattingOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = DocumentRangeFormattingOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or DocumentRangeFormattingOptions")
}

type DocumentOnTypeFormattingOptions struct {
	// A character on which formatting should be triggered, like `}`.
	FirstTriggerCharacter string `json:"firstTriggerCharacter,required"`

	// More trigger characters.
	MoreTriggerCharacter []string `json:"moreTriggerCharacter,omitempty"`
}

// boolean | RenameOptions
type RenameOptions struct {
	*WorkDoneProgressOptions

	// Renames should be checked and tested before being executed.
	PrepareProvider bool `json:"prepareProvider,omitempty"`
}

func (s *RenameOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = RenameOptions{}
		}
		return nil
	}

	type __ RenameOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = RenameOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or RenameOptions")
}

// boolean | FoldingRangeOptions | FoldingRangeRegistrationOptions
type FoldingRangeOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *FoldingRangeOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = FoldingRangeOptions{}
		}
		return nil
	}

	type __ FoldingRangeOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = FoldingRangeOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or FoldingRangeOptions")
}

type ExecuteCommandOptions struct {
	*WorkDoneProgressOptions

	// The commands to be executed on the server
	Commands []string `json:"commands"`
}

// boolean | SelectionRangeOptions | SelectionRangeRegistrationOptions
type SelectionRangeOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *SelectionRangeOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = SelectionRangeOptions{}
		}
		return nil
	}

	type __ SelectionRangeOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = SelectionRangeOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or SelectionRangeOptions")
}

// boolean | LinkedEditingRangeOptions | LinkedEditingRangeRegistrationOptions
type LinkedEditingRangeOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *LinkedEditingRangeOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = LinkedEditingRangeOptions{}
		}
		return nil
	}

	type __ LinkedEditingRangeOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = LinkedEditingRangeOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or LinkedEditingRangeOptions")
}

// boolean | CallHierarchyOptions | CallHierarchyRegistrationOptions
type CallHierarchyOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions
}

func (s *CallHierarchyOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = CallHierarchyOptions{}
		}
		return nil
	}

	type __ CallHierarchyOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = CallHierarchyOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or CallHierarchyOptions")
}

type SemanticTokensOptions struct {
	*TextDocumentRegistrationOptions
	*StaticRegistrationOptions

	// WorkDoneProgressOptions
	// The legend used by the server
	Legend SemanticTokensLegend `json:"legend,required"`

	// Server supports providing semantic tokens for a specific range
	// of a document.
	// type: boolean | { }
	Range BooleanOrEmptyStruct `json:"range,required"`

	// Server supports providing semantic tokens for a full document.
	// type: boolean | { delta?: boolean }
	Full *SemanticTokenFullOptions `json:"full,omitempty"`
}

type SemanticTokenFullOptions struct {
	// The server supports deltas for full documents.
	Delta bool `json:"delta,omitempty"`
}

type BooleanOrEmptyStruct bool

func (x *BooleanOrEmptyStruct) UnmarshalJSON(data []byte) error {
	var s struct{}
	if err := json.Unmarshal(data, &s); err == nil {
		*x = true
		return nil
	}
	var b bool
	err := json.Unmarshal(data, &b)
	*x = BooleanOrEmptyStruct(b)
	return err

}

type SemanticTokensLegend struct {
	// The token types a server uses.
	TokenTypes []string `json:"tokenTypes"`

	// The token modifiers a server uses.
	TokenModifiers []string `json:"tokenModifiers"`
}

// boolean | MonikerOptions | MonikerRegistrationOptions is defined as follows:
type MonikerOptions struct {
	*WorkDoneProgressOptions
	*TextDocumentRegistrationOptions
}

func (s *MonikerOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = MonikerOptions{}
		}
		return nil
	}

	type __ MonikerOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = MonikerOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or MonikerOptions")
}

type WorkspaceSymbolOptions struct {
	*WorkDoneProgressOptions
}

type WorkspaceSymbolRegistrationOptions struct {
	*WorkspaceSymbolOptions
}

// boolean | WorkspaceSymbolOptions where WorkspaceSymbolOptions is defined as follows:
func (s *WorkspaceSymbolOptions) UnmarshalJSON(data []byte) error {
	save := false
	if err := json.Unmarshal(data, &save); err == nil {
		if save {
			*s = WorkspaceSymbolOptions{}
		}
		return nil
	}

	type __ WorkspaceSymbolOptions // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err == nil {
		*s = WorkspaceSymbolOptions(res)
		return nil
	}
	return fmt.Errorf("expected boolean or WorkspaceSymbolOptions")
}

type WorkspaceFoldersServerCapabilities struct {
	// The server has support for workspace folders
	Supported bool `json:"supported,omitempty"`

	// Whether the server wants to receive workspace folder
	// change notifications.
	//
	// If a string is provided, the string is treated as an ID
	// under which the notification is registered on the client
	// side. The ID can be used to unregister for these events
	// using the `client/unregisterCapability` request.
	ChangeNotifications json.RawMessage `json:"changeNotifications,omitempty"`
}

// The options to register for file operations.
//
// @since 3.16.0
type FileOperationRegistrationOptions struct {
	// The actual filters.
	Filters []FileOperationFilter `json:"filters"`
}

// A filter to describe in which file operation requests or notifications
// the server is interested in.
//
// @since 3.16.0
type FileOperationFilter struct {
	// A Uri like `file` or `untitled`.
	Scheme string `json:"scheme,omitempty"`

	// The actual file operation pattern.
	Pattern FileOperationPattern `json:"pattern"`
}

// A pattern to describe in which file operation requests or notifications
// the server is interested in.
//
// @since 3.16.0
type FileOperationPattern struct {
	// The glob pattern to match. Glob patterns can have the following syntax:
	// - `*` to match one or more characters in a path segment
	// - `?` to match on one character in a path segment
	// - `**` to match any number of path segments, including none
	// - `{}` to group sub patterns into an OR expression. (e.g. `**​/*.{ts,js}`
	//   matches all TypeScript and JavaScript files)
	// - `[]` to declare a range of characters to match in a path segment
	//   (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	// - `[!...]` to negate a range of characters to match in a path segment
	//   (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but
	//   not `example.0`)
	Glob string `json:"glob"`

	// Whether to match files or folders with this pattern.
	//
	// Matches both if undefined.
	Matches *FileOperationPatternKind `json:"matches,omitempty"`

	// Additional options used during matching.
	Options *FileOperationPatternOptions `json:"options,omitempty"`
}

// A pattern kind describing if a glob pattern matches a file a folder or
// both.
//
// @since 3.16.0
type FileOperationPatternKind string

// The pattern matches a file only.
const FileOperationPatternKindFile FileOperationPatternKind = "file"

// The pattern matches a folder only.
const FileOperationPatternKindFolder FileOperationPatternKind = "folder"

// Matching options for the file operation pattern.
//
// @since 3.16.0
type FileOperationPatternOptions struct {
	// The pattern should be matched ignoring casing.
	IgnoreCase bool `json:"ignoreCase,omitempty"`
}
