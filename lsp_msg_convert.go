package lsp

import (
	"go.bug.st/json"
)

//   ___   REQUEST
//      \
//  /___/
//  \
func DecodeClientRequestParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		var res InitializeParams
		return &res, json.Unmarshal(req, &res)
	case "shutdown":
		return nil, nil
	case "workspace/symbol":
		var res WorkspaceSymbolParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/executeCommand":
		var res ExecuteCommandParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/willCreateFiles":
		var res CreateFilesParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/willRenameFiles":
		var res RenameFilesParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/willDeleteFiles":
		var res DeleteFilesParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/willSaveWaitUntil":
		var res WillSaveTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/completion":
		var res CompletionParams
		return &res, json.Unmarshal(req, &res)
	case "completionItem/resolve":
		var res CompletionItem
		return &res, json.Unmarshal(req, &res)
	case "textDocument/hover":
		var res HoverParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/signatureHelp":
		var res SignatureHelpParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/declaration":
		var res DeclarationParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/definition":
		var res DefinitionParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/typeDefinition":
		var res TypeDefinitionParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/implementation":
		var res ImplementationParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/references":
		var res ReferenceParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/documentHighlight":
		var res DocumentHighlightParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/documentSymbol":
		var res DocumentSymbolParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/codeAction":
		var res CodeActionParams
		return &res, json.Unmarshal(req, &res)
	case "codeAction/resolve":
		var res CodeAction
		return &res, json.Unmarshal(req, &res)
	case "textDocument/codeLens":
		var res CodeLensParams
		return &res, json.Unmarshal(req, &res)
	case "codeLens/resolve":
		var res CodeLens
		return &res, json.Unmarshal(req, &res)
	case "textDocument/documentLink":
		var res DocumentLinkParams
		return &res, json.Unmarshal(req, &res)
	case "documentLink/resolve":
		var res DocumentLink
		return &res, json.Unmarshal(req, &res)
	case "textDocument/documentColor":
		var res DocumentColorParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/colorPresentation":
		var res ColorPresentationParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/formatting":
		var res DocumentFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/rangeFormatting":
		var res DocumentRangeFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/onTypeFormatting":
		var res DocumentOnTypeFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/rename":
		var res RenameParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/prepareRename":
		var res PrepareRenameParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/foldingRange":
		var res FoldingRangeParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/selectionRange":
		var res SelectionRangeParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/prepareCallHierarchy":
		var res CallHierarchyPrepareParams
		return &res, json.Unmarshal(req, &res)
	case "callHierarchy/incomingCalls":
		var res CallHierarchyIncomingCallsParams
		return &res, json.Unmarshal(req, &res)
	case "callHierarchy/outgoingCalls":
		var res CallHierarchyOutgoingCallsParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/semanticTokens/full":
		var res SemanticTokensParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/semanticTokens/full/delta":
		var res SemanticTokensDeltaParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/semanticTokens/range":
		var res SemanticTokensRangeParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/semanticTokens/refresh":
		return nil, nil
	case "textDocument/linkedEditingRange":
		var res LinkedEditingRangeParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/moniker":
		var res MonikerParams
		return &res, json.Unmarshal(req, &res)
	default:
		panic("unimplemented message")
	}
}

//   ___
//      \
//  /___/  RESPONSE
//  \
func DecodeServerResponseResult(method string, resp json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		var res InitializeResult
		return &res, json.Unmarshal(resp, &res)
	case "shutdown":
		return nil, nil
	case "workspace/symbol":
		// result: SymbolInformation[] | null
		if string(resp) == "null" {
			return []SymbolInformation{}, nil
		}
		var res []SymbolInformation
		return &res, json.Unmarshal(resp, &res)
	case "workspace/executeCommand":
		// result: any | null
		return resp, nil // passthrough
	case "workspace/willCreateFiles":
		// result: WorkspaceEdit | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res WorkspaceEdit
		return &res, json.Unmarshal(resp, &res)
	case "workspace/willRenameFiles":
		// result: WorkspaceEdit | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res WorkspaceEdit
		return &res, json.Unmarshal(resp, &res)
	case "workspace/willDeleteFiles":
		// result: WorkspaceEdit | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res WorkspaceEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/willSaveWaitUntil":
		// result: TextEdit[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []TextEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/completion":
		// result: CompletionItem[] | CompletionList | null.
		// If a CompletionItem[] is provided it is interpreted to be complete. So it is the same as { isIncomplete: false, items }
		var completionItems []CompletionItem
		if err := json.Unmarshal(resp, &completionItems); err == nil {
			return &CompletionList{
				IsIncomplete: false,
				Items:        completionItems,
			}, nil
		}
		if string(resp) == "null" {
			return &CompletionList{
				IsIncomplete: false,
				Items:        []CompletionItem{},
			}, nil
		}
		var res CompletionList
		return &res, json.Unmarshal(resp, &res)
	case "completionItem/resolve":
		var res CompletionItem
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/hover":
		// result: Hover | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res Hover
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/signatureHelp":
		// result: SignatureHelp | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res SignatureHelp
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/declaration":
		// result: Location | Location[] | LocationLink[] |null
		if string(resp) == "null" {
			return nil, nil
		}
		var location Location
		if err := json.Unmarshal(resp, &location); err == nil {
			return &location, nil
		}
		var locations []Location
		if err := json.Unmarshal(resp, &locations); err == nil {
			return &locations, nil
		}
		var locationLinks []LocationLink
		return &locationLinks, json.Unmarshal(resp, &locationLinks)
	case "textDocument/definition":
		// result: Location | Location[] | LocationLink[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var location Location
		if err := json.Unmarshal(resp, &location); err == nil {
			return &location, nil
		}
		var locations []Location
		if err := json.Unmarshal(resp, &locations); err == nil {
			return &locations, nil
		}
		var locationLinks []LocationLink
		return &locationLinks, json.Unmarshal(resp, &locationLinks)
	case "textDocument/typeDefinition":
		// result: Location | Location[] | LocationLink[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var location Location
		if err := json.Unmarshal(resp, &location); err == nil {
			return &location, nil
		}
		var locations []Location
		if err := json.Unmarshal(resp, &locations); err == nil {
			return &locations, nil
		}
		var locationLinks []LocationLink
		return &locationLinks, json.Unmarshal(resp, &locationLinks)
	case "textDocument/implementation":
		// result: Location | Location[] | LocationLink[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var location Location
		if err := json.Unmarshal(resp, &location); err == nil {
			return &location, nil
		}
		var locations []Location
		if err := json.Unmarshal(resp, &locations); err == nil {
			return &locations, nil
		}
		var locationLinks []LocationLink
		return &locationLinks, json.Unmarshal(resp, &locationLinks)
	case "textDocument/references":
		// result: Location[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var locations []Location
		return &locations, json.Unmarshal(resp, &locations)
	case "textDocument/documentHighlight":
		// result: DocumentHighlight[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []DocumentHighlight
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/documentSymbol":
		// result: DocumentSymbol[] | SymbolInformation[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var documentSymbols []DocumentSymbol
		if err := json.Unmarshal(resp, &documentSymbols); err == nil {
			return documentSymbols, nil
		}
		var symbolInfomation []SymbolInformation
		return &symbolInfomation, json.Unmarshal(resp, &symbolInfomation)
	case "textDocument/codeAction":
		// result: (Command | CodeAction)[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []CommandOrCodeAction
		return res, json.Unmarshal(resp, &res)
	case "codeAction/resolve":
		var res CodeAction
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/codeLens":
		// result: CodeLens[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []CodeLens
		return &res, json.Unmarshal(resp, &res)
	case "codeLens/resolve":
		var res CodeLens
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/documentLink":
		// result: DocumentLink[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []DocumentLink
		return &res, json.Unmarshal(resp, &res)
	case "documentLink/resolve":
		var res DocumentLink
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/documentColor":
		var res []ColorInformation
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/colorPresentation":
		var res []ColorPresentation
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/formatting":
		// result: TextEdit[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []TextEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/rangeFormatting":
		// result: TextEdit[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []TextEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/onTypeFormatting":
		// result: TextEdit[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []TextEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/rename":
		// result: WorkspaceEdit | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res WorkspaceEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/prepareRename":
		// result: Range | { range: Range, placeholder: string } | { defaultBehavior: boolean } | null
		panic("unimplemented")
	case "textDocument/foldingRange":
		// result: FoldingRange[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []FoldingRange
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/selectionRange":
		// result: SelectionRange[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []SelectionRange
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/prepareCallHierarchy":
		// result: CallHierarchyItem[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []CallHierarchyItem
		return &res, json.Unmarshal(resp, &res)
	case "callHierarchy/incomingCalls":
		// result: CallHierarchyIncomingCall[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []CallHierarchyIncomingCall
		return &res, json.Unmarshal(resp, &res)
	case "callHierarchy/outgoingCalls":
		// result: CallHierarchyOutgoingCall[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []CallHierarchyOutgoingCall
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/semanticTokens/full":
		// result: SemanticTokens | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res SemanticTokens
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/semanticTokens/full/delta":
		// result: SemanticTokens | SemanticTokensDelta | null
		if string(resp) == "null" {
			return nil, nil
		}
		var delta SemanticTokensDelta
		if err := json.Unmarshal(resp, &delta); err == nil {
			return &delta, nil
		}
		var res SemanticTokens
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/semanticTokens/range":
		// result: SemanticTokens | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res SemanticTokens
		return &res, json.Unmarshal(resp, &res)
	case "workspace/semanticTokens/refresh":
		return nil, nil
	case "textDocument/linkedEditingRange":
		// result: LinkedEditingRanges | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res LinkedEditingRanges
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/moniker":
		// result: Moniker[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []Moniker
		return &res, json.Unmarshal(resp, &res)
	default:
		panic("unimplemented message")
	}
}

//   ___  REQUEST
//  /
//  \___\
//      /
func DecodeServerRequestParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "window/showMessageRequest":
		var res ShowMessageRequestParams
		return &res, json.Unmarshal(req, &res)
	case "window/showDocument":
		var res ShowDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "window/workDoneProgress/create":
		var res WorkDoneProgressCreateParams
		return &res, json.Unmarshal(req, &res)
	case "client/registerCapability":
		var res RegistrationParams
		return &res, json.Unmarshal(req, &res)
	case "client/unregisterCapability":
		var res UnregistrationParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/workspaceFolders":
		return nil, nil
	case "workspace/configuration":
		var res ConfigurationParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/applyEdit":
		var res ApplyWorkspaceEditParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/codeLens/refresh":
		return nil, nil
	default:
		panic("unimplemented message")
	}
}

//   ___
//  /
//  \___\  RESPONSE
//      /
func DecodeClientResponseResult(method string, resp json.RawMessage) (interface{}, error) {
	switch method {
	case "window/showMessageRequest":
		// result: the selected MessageActionItem | null if none got selected
		if string(resp) == "null" {
			return nil, nil
		}
		var res MessageActionItem
		return &res, json.Unmarshal(resp, &res)
	case "window/showDocument":
		var res ShowDocumentResult
		return &res, json.Unmarshal(resp, &res)
	case "window/workDoneProgress/create":
		return nil, nil
	case "client/registerCapability":
		return nil, nil
	case "client/unregisterCapability":
		return nil, nil
	case "workspace/workspaceFolders":
		// result: WorkspaceFolder[] | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res []WorkspaceFolder
		return res, json.Unmarshal(resp, &res)
	case "workspace/configuration":
		// result: any[]
		var res []json.RawMessage
		return res, json.Unmarshal(resp, &res)
	case "workspace/applyEdit":
		var res ApplyWorkspaceEditResult
		return &res, json.Unmarshal(resp, &res)
	case "workspace/codeLens/refresh":
		return nil, nil
	default:
		panic("unimplemented message")
	}
}

// ____\  NOTIFICATION
//     /
func DecodeClientNotificationParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "$/progress":
		var res ProgressParams
		return &res, json.Unmarshal(req, &res)
	case "$/cancelRequrest":
		panic("should not reach here")
	case "initialized":
		var res InitializeParams
		return &res, json.Unmarshal(req, &res)
	case "exit":
		return nil, nil
	case "$/setTrace":
		var res SetTraceParams
		return &res, json.Unmarshal(req, &res)
	case "window/workDoneProgress/cancel":
		var res WorkDoneProgressCancelParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didChangeWorkspaceFolders":
		var res DidChangeWorkspaceFoldersParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didChangeConfiguration":
		var res DidChangeConfigurationParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didChangeWatchedFiles":
		var res DidChangeWatchedFilesParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didCreateFiles":
		var res CreateFilesParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didRenameFiles":
		var res RenameFilesParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didDeleteFiles":
		var res DeleteFilesParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didOpen":
		var res DidOpenTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didChange":
		var res DidChangeTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/willSave":
		var res WillSaveTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didSave":
		var res DidSaveTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didClose":
		var res DidCloseTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	default:
		panic("unimplemented message")
	}
}

// NOTIFICATION /___
//              \
func DecodeServerNotificationParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "$/progress":
		var res ProgressParams
		return &res, json.Unmarshal(req, &res)
	case "$/cancelRequrest":
		panic("should not reach here")
	case "$/logTrace":
		var res LogTraceParams
		return &res, json.Unmarshal(req, &res)
	case "window/showMessage":
		var res ShowMessageParams
		return &res, json.Unmarshal(req, &res)
	case "window/logMessage":
		var res LogMessageParams
		return &res, json.Unmarshal(req, &res)
	case "telemetry/event":
		// params: ‘object’ | ‘number’ | ‘boolean’ | ‘string’;
		return req, nil // passthrough
	case "textDocument/publishDiagnostics":
		var res PublishDiagnosticsParams
		return &res, json.Unmarshal(req, &res)
	default:
		panic("unimplemented message")
	}
}

func EncodeMessage(msg interface{}) json.RawMessage {
	raw, _ := json.Marshal(msg)
	return raw
}
