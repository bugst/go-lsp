package lsp

import (
	"encoding/json"
	"fmt"
)

type ServerCapabilities struct {

	// Defines how text documents are synced. Is either a detailed structure
	// defining each notification or for backwards compatibility the
	// TextDocumentSyncKind number. If omitted it defaults to
	// `TextDocumentSyncKind.None`.
	TextDocumentSync *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`

	/*
		// The server provides completion support.
		completionProvider?: CompletionOptions;

		// The server provides hover support.
		hoverProvider?: boolean|HoverOptions;

		// The server provides signature help support.
		signatureHelpProvider?: SignatureHelpOptions;

		// The server provides go to declaration support.
		//
		// @since 3.14.0
		declarationProvider?: boolean|DeclarationOptions|DeclarationRegistrationOptions;

		// The server provides goto definition support.
		definitionProvider?: boolean|DefinitionOptions;

		// The server provides goto type definition support.
		//
		// @since 3.6.0
		typeDefinitionProvider?: boolean|TypeDefinitionOptions|TypeDefinitionRegistrationOptions;

		// The server provides goto implementation support.
		//
		// @since 3.6.0
		implementationProvider?: boolean|ImplementationOptions|ImplementationRegistrationOptions;

		// The server provides find references support.
		referencesProvider?: boolean|ReferenceOptions;

		// The server provides document highlight support.
		documentHighlightProvider?: boolean|DocumentHighlightOptions;

		// The server provides document symbol support.
		documentSymbolProvider?: boolean|DocumentSymbolOptions;

		// The server provides code actions. The `CodeActionOptions` return type is
		// only valid if the client signals code action literal support via the
		// property `textDocument.codeAction.codeActionLiteralSupport`.
		codeActionProvider?: boolean|CodeActionOptions;

		// The server provides code lens.
		codeLensProvider?: CodeLensOptions;

		// The server provides document link support.
		documentLinkProvider?: DocumentLinkOptions;

		// The server provides color provider support.
		//
		// @since 3.6.0
		colorProvider?: boolean|DocumentColorOptions|DocumentColorRegistrationOptions;

		// The server provides document formatting.
		documentFormattingProvider?: boolean|DocumentFormattingOptions;

		// The server provides document range formatting.
		documentRangeFormattingProvider?: boolean|DocumentRangeFormattingOptions;

		// The server provides document formatting on typing.
		documentOnTypeFormattingProvider?: DocumentOnTypeFormattingOptions;

		// The server provides rename support. RenameOptions may only be
		// specified if the client states that it supports
		// `prepareSupport` in its initial `initialize` request.
		renameProvider?: boolean|RenameOptions;

		// The server provides folding provider support.
		//
		// @since 3.10.0
		foldingRangeProvider?: boolean|FoldingRangeOptions|FoldingRangeRegistrationOptions;

		// The server provides execute command support.
		executeCommandProvider?: ExecuteCommandOptions;

		// The server provides selection range support.
		//
		// @since 3.15.0
		selectionRangeProvider?: boolean|SelectionRangeOptions|SelectionRangeRegistrationOptions;

		// The server provides linked editing range support.
		//
		// @since 3.16.0
		linkedEditingRangeProvider?: boolean|LinkedEditingRangeOptions|LinkedEditingRangeRegistrationOptions;

		// The server provides call hierarchy support.
		//
		// @since 3.16.0
		callHierarchyProvider?: boolean|CallHierarchyOptions|CallHierarchyRegistrationOptions;

		// The server provides semantic tokens support.
		//
		// @since 3.16.0
		semanticTokensProvider?: SemanticTokensOptions|SemanticTokensRegistrationOptions;

		// Whether server provides moniker support.
		//
		// @since 3.16.0
		monikerProvider?: boolean|MonikerOptions|MonikerRegistrationOptions;

		// The server provides workspace symbol support.
		workspaceSymbolProvider?: boolean|WorkspaceSymbolOptions;

		// Workspace specific server capabilities
		workspace?: {

			// The server supports workspace folder.
			//
			// @since 3.6.0
			workspaceFolders?: WorkspaceFoldersServerCapabilities;

			// The server is interested in file notifications/requests.
			//
			// @since 3.16.0
			fileOperations?: {

				// The server is interested in receiving didCreateFiles
				// notifications.
				didCreate?: FileOperationRegistrationOptions;

				// The server is interested in receiving willCreateFiles requests.
				willCreate?: FileOperationRegistrationOptions;

				// The server is interested in receiving didRenameFiles
				// notifications.
				didRename?: FileOperationRegistrationOptions;

				// The server is interested in receiving willRenameFiles requests.
				willRename?: FileOperationRegistrationOptions;

				// The server is interested in receiving didDeleteFiles file
				// notifications.
				didDelete?: FileOperationRegistrationOptions;

				// The server is interested in receiving willDeleteFiles file
				// requests.
				willDelete?: FileOperationRegistrationOptions;
			};
		};

		// Experimental server capabilities.
		experimental?: any;
	*/
}

type TextDocumentSyncKind int

const TextDocumentSyncKindNone TextDocumentSyncKind = 0
const TextDocumentSyncKindFull TextDocumentSyncKind = 1
const TextDocumentSyncKindIncremental TextDocumentSyncKind = 2

type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server. If omitted open
	// close notification should not be sent.
	OpenClose *bool `json:"openClose,omitempty"`

	// Change notifications are sent to the server. See
	// TextDocumentSyncKind.None, TextDocumentSyncKind.Full and
	// TextDocumentSyncKind.Incremental. If omitted it defaults to
	// TextDocumentSyncKind.None.
	Change *TextDocumentSyncKind `json:"change,omitempty"`

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
		*s = SaveOptions{}
		return nil
	}

	var res SaveOptions
	if err := json.Unmarshal(data, &res); err != nil {
		*s = res
		return nil
	}
	return fmt.Errorf("expected boolean or SaveOptions")
}
