package lsp

import (
	"context"
	"io"

	"go.bug.st/json"
	"go.bug.st/lsp/jsonrpc"
)

// ServerMessagesHandler interface has all the methods that an LSP Client should
// implement to correctly parse incoming messages
type ServerMessagesHandler interface {
	// Response <- Request

	WindowShowMessageRequest(context.Context, jsonrpc.FunctionLogger, *ShowMessageRequestParams) (*MessageActionItem, *jsonrpc.ResponseError)
	WindowShowDocument(context.Context, jsonrpc.FunctionLogger, *ShowDocumentParams) (*ShowDocumentResult, *jsonrpc.ResponseError)
	WindowWorkDoneProgressCreate(context.Context, jsonrpc.FunctionLogger, *WorkDoneProgressCreateParams) *jsonrpc.ResponseError
	ClientRegisterCapability(context.Context, jsonrpc.FunctionLogger, *RegistrationParams) *jsonrpc.ResponseError
	ClientUnregisterCapability(context.Context, jsonrpc.FunctionLogger, *UnregistrationParams) *jsonrpc.ResponseError
	WorkspaceWorkspaceFolders(context.Context, jsonrpc.FunctionLogger) ([]WorkspaceFolder, *jsonrpc.ResponseError)
	WorkspaceConfiguration(context.Context, jsonrpc.FunctionLogger, *ConfigurationParams) ([]json.RawMessage, *jsonrpc.ResponseError)
	WorkspaceApplyEdit(context.Context, jsonrpc.FunctionLogger, *ApplyWorkspaceEditParams) (*ApplyWorkspaceEditResult, *jsonrpc.ResponseError)
	WorkspaceCodeLensRefresh(context.Context, jsonrpc.FunctionLogger) *jsonrpc.ResponseError

	// Notifications <-

	Progress(jsonrpc.FunctionLogger, *ProgressParams)
	LogTrace(jsonrpc.FunctionLogger, *LogTraceParams)
	WindowShowMessage(jsonrpc.FunctionLogger, *ShowMessageParams)
	WindowLogMessage(jsonrpc.FunctionLogger, *LogMessageParams)
	TelemetryEvent(jsonrpc.FunctionLogger, json.RawMessage)
	TextDocumentPublishDiagnostics(jsonrpc.FunctionLogger, *PublishDiagnosticsParams)
}

// Client is an LSP Client
type Client struct {
	conn         *jsonrpc.Connection
	handler      ServerMessagesHandler
	errorHandler func(e error)
}

func NewClient(in io.Reader, out io.Writer, handler ServerMessagesHandler) *Client {
	client := &Client{
		errorHandler: func(e error) {},
	}
	client.handler = handler
	client.conn = jsonrpc.NewConnection(
		in, out,
		client.requestDispatcher,
		client.notificationDispatcher,
		client.errorHandler)
	return client
}

func (client *Client) SetLogger(l jsonrpc.Logger) {
	client.conn.SetLogger(l)
}

func (client *Client) SetErrorHandler(handler func(e error)) {
	client.errorHandler = handler
}

func (client *Client) Run() {
	client.conn.Run()
}

func (client *Client) notificationDispatcher(logger jsonrpc.FunctionLogger, method string, req json.RawMessage) {
	switch method {
	case "$/progress":
		var param ProgressParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		client.handler.Progress(logger, &param)
	case "$/cancelRequrest":
		panic("should not reach here")
	case "$/logTrace":
		var param LogTraceParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		client.handler.LogTrace(logger, &param)
	case "window/showMessage":
		var param ShowMessageParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		client.handler.WindowShowMessage(logger, &param)
	case "window/logMessage":
		var param LogMessageParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		client.handler.WindowLogMessage(logger, &param)
	case "telemetry/event":
		// params: ‘object’ | ‘number’ | ‘boolean’ | ‘string’;
		client.handler.TelemetryEvent(logger, req) // passthrough
	case "textDocument/publishDiagnostics":
		var param PublishDiagnosticsParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		client.handler.TextDocumentPublishDiagnostics(logger, &param)
	default:
		panic("unimplemented message")
	}
}

func (client *Client) requestDispatcher(ctx context.Context, logger jsonrpc.FunctionLogger, method string, req json.RawMessage, respCallback func(json.RawMessage, *jsonrpc.ResponseError)) {
	resp := func(res interface{}, err *jsonrpc.ResponseError) {
		respCallback(EncodeMessage(res), err)
	}
	switch method {
	case "window/showMessageRequest":
		var param ShowMessageRequestParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(client.handler.WindowShowMessageRequest(ctx, logger, &param))
	case "window/showDocument":
		var param ShowDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(client.handler.WindowShowDocument(ctx, logger, &param))
	case "window/workDoneProgress/create":
		var param WorkDoneProgressCreateParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(nil, client.handler.WindowWorkDoneProgressCreate(ctx, logger, &param))
	case "client/registerCapability":
		var param RegistrationParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(nil, client.handler.ClientRegisterCapability(ctx, logger, &param))
	case "client/unregisterCapability":
		var param UnregistrationParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(nil, client.handler.ClientUnregisterCapability(ctx, logger, &param))
	case "workspace/workspaceFolders":
		resp(client.handler.WorkspaceWorkspaceFolders(ctx, logger))
	case "workspace/configuration":
		var param ConfigurationParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(client.handler.WorkspaceConfiguration(ctx, logger, &param))
	case "workspace/applyEdit":
		var param ApplyWorkspaceEditParams
		if err := json.Unmarshal(req, &param); err != nil {
			client.errorHandler(err)
			return
		}
		resp(client.handler.WorkspaceApplyEdit(ctx, logger, &param))
	case "workspace/codeLens/refresh":
		resp(nil, client.handler.WorkspaceCodeLensRefresh(ctx, logger))
	default:
		panic("unimplemented message")
	}
}

// Requests to Server

func (client *Client) Initialize(ctx context.Context, param *InitializeParams) (*InitializeResult, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "initialize", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res InitializeResult
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) Shutdown(ctx context.Context) (*jsonrpc.ResponseError, error) {
	_, respErr, err := client.conn.SendRequest(ctx, "shutdown", EncodeMessage(jsonrpc.NullResult))
	return respErr, err
}

func (client *Client) WorkspaceSymbol(ctx context.Context, param *WorkspaceSymbolParams) ([]SymbolInformation, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "workspace/symbol", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: SymbolInformation[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []SymbolInformation
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) WorkspaceExecuteCommand(ctx context.Context, param *ExecuteCommandParams) (json.RawMessage, *jsonrpc.ResponseError, error) {
	return client.conn.SendRequest(ctx, "workspace/executeCommand", EncodeMessage(param))
}

func (client *Client) WorkspaceWillCreateFiles(ctx context.Context, param *CreateFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "workspace/willCreateFiles", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: WorkspaceEdit | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res WorkspaceEdit
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) WorkspaceWillRenameFiles(ctx context.Context, param *RenameFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "workspace/willRenameFiles", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: WorkspaceEdit | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res WorkspaceEdit
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) WorkspaceWillDeleteFiles(ctx context.Context, param *DeleteFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "workspace/willDeleteFiles", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: WorkspaceEdit | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res WorkspaceEdit
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentWillSaveWaitUntil(ctx context.Context, param *WillSaveTextDocumentParams) ([]TextEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/willSaveWaitUntil", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: TextEdit[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []TextEdit
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentCompletion(ctx context.Context, param *CompletionParams) (*CompletionList, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/completion", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: CompletionItem[] | CompletionList | null.
	// If a CompletionItem[] is provided it is interpreted to be complete. So it is the same as { isIncomplete: false, items }
	var completionItems []CompletionItem
	if err := json.Unmarshal(resp, &completionItems); err == nil {
		return &CompletionList{
			IsIncomplete: false,
			Items:        completionItems,
		}, respErr, nil
	}
	if string(resp) == "null" {
		return &CompletionList{
			IsIncomplete: false,
			Items:        []CompletionItem{},
		}, respErr, nil
	}
	var res CompletionList
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) CompletionItemResolve(ctx context.Context, param *CompletionItem) (*CompletionItem, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "completionItem/resolve", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res CompletionItem
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentHover(ctx context.Context, param *HoverParams) (*Hover, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/hover", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: Hover | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res Hover
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentSignatureHelp(ctx context.Context, param *SignatureHelpParams) (*SignatureHelp, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/signatureHelp", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: SignatureHelp | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res SignatureHelp
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentDeclaration(ctx context.Context, param *DeclarationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/declaration", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: Location | Location[] | LocationLink[] |null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var location Location
	if err := json.Unmarshal(resp, &location); err == nil {
		return []Location{location}, nil, respErr, nil
	}
	var locations []Location
	if err := json.Unmarshal(resp, &locations); err == nil {
		return locations, nil, respErr, nil
	}
	var locationLinks []LocationLink
	return nil, locationLinks, respErr, json.Unmarshal(resp, &locationLinks)
}

func (client *Client) TextDocumentDefinition(ctx context.Context, param *DefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/definition", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: Location | Location[] | LocationLink[] |null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var location Location
	if err := json.Unmarshal(resp, &location); err == nil {
		return []Location{location}, nil, respErr, nil
	}
	var locations []Location
	if err := json.Unmarshal(resp, &locations); err == nil {
		return locations, nil, respErr, nil
	}
	var locationLinks []LocationLink
	return nil, locationLinks, respErr, json.Unmarshal(resp, &locationLinks)
}

func (client *Client) TextDocumentTypeDefinition(ctx context.Context, param *TypeDefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/typeDefinition", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: Location | Location[] | LocationLink[] |null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var location Location
	if err := json.Unmarshal(resp, &location); err == nil {
		return []Location{location}, nil, respErr, nil
	}
	var locations []Location
	if err := json.Unmarshal(resp, &locations); err == nil {
		return locations, nil, respErr, nil
	}
	var locationLinks []LocationLink
	return nil, locationLinks, respErr, json.Unmarshal(resp, &locationLinks)
}

func (client *Client) TextDocumentImplementation(ctx context.Context, param *ImplementationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/implementation", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: Location | Location[] | LocationLink[] |null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var location Location
	if err := json.Unmarshal(resp, &location); err == nil {
		return []Location{location}, nil, respErr, nil
	}
	var locations []Location
	if err := json.Unmarshal(resp, &locations); err == nil {
		return locations, nil, respErr, nil
	}
	var locationLinks []LocationLink
	return nil, locationLinks, respErr, json.Unmarshal(resp, &locationLinks)
}

func (client *Client) TextDocumentReferences(ctx context.Context, param *ReferenceParams) ([]Location, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/references", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: Location[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []Location
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentDocumentHighlight(ctx context.Context, param *DocumentHighlightParams) ([]DocumentHighlight, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/documentHighlight", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: DocumentHighlight[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []DocumentHighlight
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentDocumentSymbol(ctx context.Context, param *DocumentSymbolParams) ([]DocumentSymbol, []SymbolInformation, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/documentSymbol", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: DocumentSymbol[] | SymbolInformation[] | null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var documentSymbols []DocumentSymbol
	if err := json.Unmarshal(resp, &documentSymbols); err == nil {
		return documentSymbols, nil, respErr, nil
	}
	var symbolInfomation []SymbolInformation
	return nil, symbolInfomation, respErr, json.Unmarshal(resp, &symbolInfomation)

}

func (client *Client) TextDocumentCodeAction(ctx context.Context, param *CodeActionParams) ([]CommandOrCodeAction, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/codeAction", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: (Command | CodeAction)[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []CommandOrCodeAction
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) CodeActionResolve(ctx context.Context, param *CodeAction) (*CodeAction, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "codeAction/resolve", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res CodeAction
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentCodeLens(ctx context.Context, param *CodeLensParams) ([]CodeLens, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/codeLens", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: CodeLens[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []CodeLens
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) CodeLensResolve(ctx context.Context, param *CodeLens) (*CodeLens, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "codeLens/resolve", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res CodeLens
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentDocumentLink(ctx context.Context, param *DocumentLinkParams) ([]DocumentLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/documentLink", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: DocumentLink[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []DocumentLink
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) DocumentLinkResolve(ctx context.Context, param *DocumentLink) (*DocumentLink, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "documentLink/resolve", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res DocumentLink
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentDocumentColor(ctx context.Context, param *DocumentColorParams) ([]ColorInformation, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/documentColor", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res []ColorInformation
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentColorPresentation(ctx context.Context, param *ColorPresentationParams) ([]ColorPresentation, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/colorPresentation", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res []ColorPresentation
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentFormatting(ctx context.Context, param *DocumentFormattingParams) ([]TextEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/formatting", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: TextEdit[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []TextEdit
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentRangeFormatting(ctx context.Context, param *DocumentRangeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/rangeFormatting", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: TextEdit[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []TextEdit
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentOnTypeFormatting(ctx context.Context, param *DocumentOnTypeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/onTypeFormatting", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: TextEdit[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []TextEdit
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentRename(ctx context.Context, param *RenameParams) (*WorkspaceEdit, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/rename", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: WorkspaceEdit | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res WorkspaceEdit
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentPrepareRename(ctx context.Context, param *PrepareRenameParams) (json.RawMessage, *jsonrpc.ResponseError, error) {
	panic("unimplemented")
	// _, _, err := client.conn.SendRequest(ctx, "textDocument/prepareRename", EncodeMessage(param))
	// if err != nil {
	// 	return nil, nil, err
	// }
	// // result: Range | { range: Range, placeholder: string } | { defaultBehavior: boolean } | null
	// return nil, nil, nil
}

func (client *Client) TextDocumentFoldingRange(ctx context.Context, param *FoldingRangeParams) ([]FoldingRange, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/foldingRange", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: FoldingRange[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []FoldingRange
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentSelectionRange(ctx context.Context, param *SelectionRangeParams) ([]SelectionRange, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/selectionRange", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: SelectionRange[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []SelectionRange
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentPrepareCallHierarchy(ctx context.Context, param *CallHierarchyPrepareParams) ([]CallHierarchyItem, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/prepareCallHierarchy", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: CallHierarchyItem[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []CallHierarchyItem
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) CallHierarchyIncomingCalls(ctx context.Context, param *CallHierarchyIncomingCallsParams) ([]CallHierarchyIncomingCall, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "callHierarchy/incomingCalls", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: CallHierarchyIncomingCall[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []CallHierarchyIncomingCall
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) CallHierarchyOutgoingCalls(ctx context.Context, param *CallHierarchyOutgoingCallsParams) ([]CallHierarchyOutgoingCall, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "callHierarchy/outgoingCalls", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: CallHierarchyOutgoingCall[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []CallHierarchyOutgoingCall
	return res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentSemanticTokensFull(ctx context.Context, param *SemanticTokensParams) (*SemanticTokens, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/semanticTokens/full", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: SemanticTokens | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res SemanticTokens
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentSemanticTokensFullDelta(ctx context.Context, param *SemanticTokensDeltaParams) (*SemanticTokens, *SemanticTokensDelta, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/semanticTokens/full/delta", EncodeMessage(param))
	if err != nil {
		return nil, nil, nil, err
	}
	// result: SemanticTokens | SemanticTokensDelta | null
	if string(resp) == "null" {
		return nil, nil, respErr, nil
	}
	var delta SemanticTokensDelta
	if err := json.Unmarshal(resp, &delta); err == nil {
		return nil, &delta, respErr, nil
	}
	var res SemanticTokens
	return &res, nil, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentSemanticTokensRange(ctx context.Context, param *SemanticTokensRangeParams) (*SemanticTokens, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/semanticTokens/range", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: SemanticTokens | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res SemanticTokens
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) WorkspaceSemanticTokensRefresh(ctx context.Context) (*jsonrpc.ResponseError, error) {
	_, respErr, err := client.conn.SendRequest(ctx, "workspace/semanticTokens/refresh", EncodeMessage(jsonrpc.NullResult))
	return respErr, err
}

func (client *Client) TextDocumentLinkedEditingRange(ctx context.Context, param *LinkedEditingRangeParams) (*LinkedEditingRanges, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/linkedEditingRange", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: LinkedEditingRanges | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res LinkedEditingRanges
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (client *Client) TextDocumentMoniker(ctx context.Context, param *MonikerParams) ([]Moniker, *jsonrpc.ResponseError, error) {
	resp, respErr, err := client.conn.SendRequest(ctx, "textDocument/moniker", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: Moniker[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []Moniker
	return res, respErr, json.Unmarshal(resp, &res)
}

// Notifications to Server

func (client *Client) Progress(param *ProgressParams) error {
	return client.conn.SendNotification("$/progress", EncodeMessage(param))
}

func (client *Client) Initialized(param *InitializedParams) error {
	return client.conn.SendNotification("initialized", EncodeMessage(param))
}

func (client *Client) Exit() error {
	return client.conn.SendNotification("exit", EncodeMessage(jsonrpc.NullResult))
}

func (client *Client) SetTrace(param *SetTraceParams) error {
	return client.conn.SendNotification("$/setTrace", EncodeMessage(param))
}

func (client *Client) WindowWorkDoneProgressCancel(param *WorkDoneProgressCancelParams) error {
	return client.conn.SendNotification("window/workDoneProgress/cancel", EncodeMessage(param))
}

func (client *Client) WorkspaceDidChangeWorkspaceFolders(param *DidChangeWorkspaceFoldersParams) error {
	return client.conn.SendNotification("workspace/didChangeWorkspaceFolders", EncodeMessage(param))
}

func (client *Client) WorkspaceDidChangeConfiguration(param *DidChangeConfigurationParams) error {
	return client.conn.SendNotification("workspace/didChangeConfiguration", EncodeMessage(param))
}

func (client *Client) WorkspaceDidChangeWatchedFiles(param *DidChangeWatchedFilesParams) error {
	return client.conn.SendNotification("workspace/didChangeWatchedFiles", EncodeMessage(param))
}

func (client *Client) WorkspaceDidCreateFiles(param *CreateFilesParams) error {
	return client.conn.SendNotification("workspace/didCreateFiles", EncodeMessage(param))
}

func (client *Client) WorkspaceDidRenameFiles(param *RenameFilesParams) error {
	return client.conn.SendNotification("workspace/didRenameFiles", EncodeMessage(param))
}

func (client *Client) WorkspaceDidDeleteFiles(param *DeleteFilesParams) error {
	return client.conn.SendNotification("workspace/didDeleteFiles", EncodeMessage(param))
}

func (client *Client) TextDocumentDidOpen(param *DidOpenTextDocumentParams) error {
	return client.conn.SendNotification("textDocument/didOpen", EncodeMessage(param))
}

func (client *Client) TextDocumentDidChange(param *DidChangeTextDocumentParams) error {
	return client.conn.SendNotification("textDocument/didChange", EncodeMessage(param))
}

func (client *Client) TextDocumentWillSave(param *WillSaveTextDocumentParams) error {
	return client.conn.SendNotification("textDocument/willSave", EncodeMessage(param))
}

func (client *Client) TextDocumentDidSave(param *DidSaveTextDocumentParams) error {
	return client.conn.SendNotification("textDocument/didSave", EncodeMessage(param))
}

func (client *Client) TextDocumentDidClose(param *DidCloseTextDocumentParams) error {
	return client.conn.SendNotification("textDocument/didClose", EncodeMessage(param))
}
