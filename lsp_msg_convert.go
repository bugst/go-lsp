package lsp

import (
	"go.bug.st/json"
)

func DecodeRequestParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		var res InitializeParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/codeAction":
		var res CodeActionParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/hover":
		var res HoverParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/documentSymbol":
		var res DocumentSymbolParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/formatting":
		var res DocumentFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/rangeFormatting":
		var res DocumentRangeFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/signatureHelp":
		var res SignatureHelpParams
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
	case "textDocument/documentHighlight":
		var res DocumentHighlightParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/references":
		var res ReferenceParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/onTypeFormatting":
		var res DocumentOnTypeFormattingParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/rename":
		var res RenameParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/executeCommand":
		var res ExecuteCommandParams
		return &res, json.Unmarshal(req, &res)
	case "window/showMessageRequest":
		var res ShowMessageRequestParams
		return &res, json.Unmarshal(req, &res)
	}
	return nil, nil
}

func DecodeNotificationParams(method string, req json.RawMessage) (interface{}, error) {
	switch method {
	case "textDocument/didOpen":
		var res DidOpenTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didClose":
		var res DidCloseTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didChange":
		var res DidChangeTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/didSave":
		var res DidSaveTextDocumentParams
		return &res, json.Unmarshal(req, &res)
	case "textDocument/completion":
		var res CompletionParams
		return &res, json.Unmarshal(req, &res)
	case "workspace/didChangeWatchedFiles":
		var res DidChangeWatchedFilesParams
		return &res, json.Unmarshal(req, &res)
	case "window/showMessage":
		var res ShowMessageParams
		return &res, json.Unmarshal(req, &res)
	}
	return nil, nil
}

func DecodeResponseResult(method string, resp json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		var res InitializeResult
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
	case "textDocument/codeAction":
		// result: (Command | CodeAction)[] | null
		if string(resp) == "null" {
			return []Command{}, nil
		}
		var resCommands []Command
		if err := json.Unmarshal(resp, &resCommands); err == nil {
			return &resCommands, nil
		}
		var resCodeActions []CodeAction
		return &resCodeActions, json.Unmarshal(resp, &resCodeActions)
	case "completionItem/resolve":
		var res CompletionItem
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/signatureHelp":
		// result: SignatureHelp | null
		var res SignatureHelp
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/hover":
		// result: Hover | null
		var res Hover
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/definition",
		"textDocument/typeDefinition",
		"textDocument/implementation":
		// result: Location | Location[] | LocationLink[] | null
		fallthrough
	case "textDocument/references":
		// result: Location[] | null
		if string(resp) == "null" {
			return []Location{}, nil
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
	case "textDocument/documentHighlight":
		// result: DocumentHighlight[] | null
		if string(resp) == "null" {
			return []DocumentHighlight{}, nil
		}
		var res []DocumentHighlight
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/formatting",
		"textDocument/rangeFormatting",
		"textDocument/onTypeFormatting":
		// result: TextEdit[] | null
		if string(resp) == "null" {
			return []TextEdit{}, nil
		}
		var res []TextEdit
		return &res, json.Unmarshal(resp, &res)
	case "textDocument/documentSymbol":
		// result: DocumentSymbol[] | SymbolInformation[] | null
		if string(resp) == "null" {
			return []DocumentSymbol{}, nil
		}
		var documentSymbols []DocumentSymbol
		if err := json.Unmarshal(resp, &documentSymbols); err == nil {
			return documentSymbols, nil
		}
		var symbolInfomation []SymbolInformation
		return &symbolInfomation, json.Unmarshal(resp, &symbolInfomation)
	case "textDocument/rename":
		// result: WorkspaceEdit | null
		if string(resp) == "null" {
			return nil, nil
		}
		var res WorkspaceEdit
		return &res, json.Unmarshal(resp, &res)
	case "workspace/symbol":
		// result: SymbolInformation[] | null
		if string(resp) == "null" {
			return []SymbolInformation{}, nil
		}
		var res []SymbolInformation
		return &res, json.Unmarshal(resp, &res)
	case "window/showMessageRequest":
		// result: the selected MessageActionItem | null if none got selected
		if string(resp) == "null" {
			return nil, nil
		}
		var res MessageActionItem
		return &res, json.Unmarshal(resp, &res)
	case "workspace/executeCommand":
		// result: any | null
		var res interface{}
		return &res, json.Unmarshal(resp, &res)
	case "workspace/applyEdit":
		var res ApplyWorkspaceEditResult
		return &res, json.Unmarshal(resp, &res)
	}
	return nil, nil
}

func EncodeMessage(msg interface{}) json.RawMessage {
	raw, _ := json.Marshal(msg)
	return raw
}
