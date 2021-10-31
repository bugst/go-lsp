package lsp

import (
	"context"
	"io"

	"go.bug.st/json"
	"go.bug.st/lsp/jsonrpc"
)

// ClientMessagesHandler interface has all the methods that an LSP Server should
// implement to correctly parse incoming messages
type ClientMessagesHandler interface {
	// Request -> Response

	Initialize(context.Context, *InitializeParams) (*InitializeResult, *jsonrpc.ResponseError)
	Shutdown(context.Context) *jsonrpc.ResponseError
	WorkspaceSymbol(context.Context, *WorkspaceSymbolParams) ([]SymbolInformation, *jsonrpc.ResponseError)
	WorkspaceExecuteCommand(context.Context, *ExecuteCommandParams) (json.RawMessage, *jsonrpc.ResponseError)
	WorkspaceWillCreateFiles(context.Context, *CreateFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	WorkspaceWillRenameFiles(context.Context, *RenameFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	WorkspaceWillDeleteFiles(context.Context, *DeleteFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	TextDocumentWillSaveWaitUntil(context.Context, *WillSaveTextDocumentParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentCompletion(context.Context, *CompletionParams) (*CompletionList, *jsonrpc.ResponseError)
	CompletionItemResolve(context.Context, *CompletionItem) (*CompletionItem, *jsonrpc.ResponseError)
	TextDocumentHover(context.Context, *HoverParams) (*Hover, *jsonrpc.ResponseError)
	TextDocumentSignatureHelp(context.Context, *SignatureHelpParams) (*SignatureHelp, *jsonrpc.ResponseError)
	TextDocumentDeclaration(context.Context, *DeclarationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentDefinition(context.Context, *DefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentTypeDefinition(context.Context, *TypeDefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentImplementation(context.Context, *ImplementationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentReferences(context.Context, *ReferenceParams) ([]Location, *jsonrpc.ResponseError)
	TextDocumentDocumentHighlight(context.Context, *DocumentHighlightParams) ([]DocumentHighlight, *jsonrpc.ResponseError)
	TextDocumentDocumentSymbol(context.Context, *DocumentSymbolParams) ([]DocumentSymbol, []SymbolInformation, *jsonrpc.ResponseError)
	TextDocumentCodeAction(context.Context, *CodeActionParams) ([]CommandOrCodeAction, *jsonrpc.ResponseError)
	CodeActionResolve(context.Context, *CodeAction) (*CodeAction, *jsonrpc.ResponseError)
	TextDocumentCodeLens(context.Context, *CodeLensParams) ([]CodeLens, *jsonrpc.ResponseError)
	CodeLensResolve(context.Context, *CodeLens) (*CodeLens, *jsonrpc.ResponseError)
	TextDocumentDocumentLink(context.Context, *DocumentLinkParams) ([]DocumentLink, *jsonrpc.ResponseError)
	DocumentLinkResolve(context.Context, *DocumentLink) (*DocumentLink, *jsonrpc.ResponseError)
	TextDocumentDocumentColor(context.Context, *DocumentColorParams) ([]ColorInformation, *jsonrpc.ResponseError)
	TextDocumentColorPresentation(context.Context, *ColorPresentationParams) ([]ColorPresentation, *jsonrpc.ResponseError)
	TextDocumentFormatting(context.Context, *DocumentFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentRangeFormatting(context.Context, *DocumentRangeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentOnTypeFormatting(context.Context, *DocumentOnTypeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentRename(context.Context, *RenameParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	//TextDocumentPrepareRename(context.Context, *PrepareRenameParams) (???, *jsonrpc.ResponseError)
	TextDocumentFoldingRange(context.Context, *FoldingRangeParams) ([]FoldingRange, *jsonrpc.ResponseError)
	TextDocumentSelectionRange(context.Context, *SelectionRangeParams) ([]SelectionRange, *jsonrpc.ResponseError)
	TextDocumentPrepareCallHierarchy(context.Context, *CallHierarchyPrepareParams) ([]CallHierarchyItem, *jsonrpc.ResponseError)
	CallHierarchyIncomingCalls(context.Context, *CallHierarchyIncomingCallsParams) ([]CallHierarchyIncomingCall, *jsonrpc.ResponseError)
	CallHierarchyOutgoingCalls(context.Context, *CallHierarchyOutgoingCallsParams) ([]CallHierarchyOutgoingCall, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensFull(context.Context, *SemanticTokensParams) (*SemanticTokens, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensFullDelta(context.Context, *SemanticTokensDeltaParams) (*SemanticTokens, *SemanticTokensDelta, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensRange(context.Context, *SemanticTokensRangeParams) (*SemanticTokens, *jsonrpc.ResponseError)
	WorkspaceSemanticTokensRefresh(context.Context) *jsonrpc.ResponseError
	TextDocumentLinkedEditingRange(context.Context, *LinkedEditingRangeParams) (*LinkedEditingRanges, *jsonrpc.ResponseError)
	TextDocumentMoniker(context.Context, *MonikerParams) ([]Moniker, *jsonrpc.ResponseError)

	// Notifications ->

	Progress(*ProgressParams)
	// CancelRequrest(*jsonrpc.CancelParams) - automatically handled by the rpc library
	Initialized(*InitializeParams)
	Exit()
	SetTrace(*SetTraceParams)
	WindowWorkDoneProgressCancel(*WorkDoneProgressCancelParams)
	WorkspaceDidChangeWorkspaceFolders(*DidChangeWorkspaceFoldersParams)
	WorkspaceDidChangeConfiguration(*DidChangeConfigurationParams)
	WorkspaceDidChangeWatchedFiles(*DidChangeWatchedFilesParams)
	WorkspaceDidCreateFiles(*CreateFilesParams)
	WorkspaceDidRenameFiles(*RenameFilesParams)
	WorkspaceDidDeleteFiles(*DeleteFilesParams)
	TextDocumentDidOpen(*DidOpenTextDocumentParams)
	TextDocumentDidChange(*DidChangeTextDocumentParams)
	TextDocumentWillSave(*WillSaveTextDocumentParams)
	TextDocumentDidSave(*DidSaveTextDocumentParams)
	TextDocumentDidClose(*DidCloseTextDocumentParams)
}

// Server is an LSP Server
type Server struct {
	conn         *jsonrpc.Connection
	handler      ClientMessagesHandler
	errorHandler func(e error)
}

// ServerHandler is an LSP Server message handler
type ServerHandler interface {
	Initialize(ctx context.Context, conn jsonrpc.Connection, params InitializeParams)
}

func NewServer(in io.Reader, out io.Writer, handler ClientMessagesHandler) *Server {
	serv := &Server{
		errorHandler: func(e error) {},
	}
	serv.handler = handler
	serv.conn = jsonrpc.NewConnection(
		in, out,
		serv.requestDispatcher,
		serv.notificationDispatcher,
		serv.errorHandler)
	return serv
}

func (serv *Server) SetErrorHandler(handler func(e error)) {
	serv.errorHandler = handler
}

func (serv *Server) Run() {
	serv.conn.Run()
}

func (serv *Server) notificationDispatcher(ctx context.Context, method string, req json.RawMessage) {
	switch method {
	case "$/progress":
		var param ProgressParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.Progress(&param)
	case "$/cancelRequrest":
		panic("should not reach here")
	case "initialized":
		var param InitializeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.Initialized(&param)
	case "exit":
		serv.handler.Exit()
	case "$/setTrace":
		var param SetTraceParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.SetTrace(&param)
	case "window/workDoneProgress/cancel":
		var param WorkDoneProgressCancelParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WindowWorkDoneProgressCancel(&param)
	case "workspace/didChangeWorkspaceFolders":
		var param DidChangeWorkspaceFoldersParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeWorkspaceFolders(&param)
	case "workspace/didChangeConfiguration":
		var param DidChangeConfigurationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeConfiguration(&param)
	case "workspace/didChangeWatchedFiles":
		var param DidChangeWatchedFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeWatchedFiles(&param)
	case "workspace/didCreateFiles":
		var param CreateFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidCreateFiles(&param)
	case "workspace/didRenameFiles":
		var param RenameFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidRenameFiles(&param)
	case "workspace/didDeleteFiles":
		var param DeleteFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidDeleteFiles(&param)
	case "textDocument/didOpen":
		var param DidOpenTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidOpen(&param)
	case "textDocument/didChange":
		var param DidChangeTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidChange(&param)
	case "textDocument/willSave":
		var param WillSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentWillSave(&param)
	case "textDocument/didSave":
		var param DidSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidSave(&param)
	case "textDocument/didClose":
		var param DidCloseTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidClose(&param)
	default:
		panic("unimplemented message")
	}
}

func (serv *Server) requestDispatcher(ctx context.Context, method string, req json.RawMessage, respCallback func(json.RawMessage, *jsonrpc.ResponseError)) {
	resp := func(res interface{}, err *jsonrpc.ResponseError) {
		respCallback(EncodeMessage(res), err)
	}
	resp2 := func(res1 interface{}, res2 interface{}, err *jsonrpc.ResponseError) {
		if res1 == nil {
			respCallback(EncodeMessage(res2), err)
		} else {
			respCallback(EncodeMessage(res1), err)
		}
	}
	switch method {
	case "initialize":
		var param InitializeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.Initialize(ctx, &param))
	case "shutdown":
		resp(nil, serv.handler.Shutdown(ctx))
	case "workspace/symbol":
		var param WorkspaceSymbolParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceSymbol(ctx, &param))
	case "workspace/executeCommand":
		var param ExecuteCommandParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceExecuteCommand(ctx, &param))
	case "workspace/willCreateFiles":
		var param CreateFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillCreateFiles(ctx, &param))
	case "workspace/willRenameFiles":
		var param RenameFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillRenameFiles(ctx, &param))
	case "workspace/willDeleteFiles":
		var param DeleteFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillDeleteFiles(ctx, &param))
	case "textDocument/willSaveWaitUntil":
		var param WillSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentWillSaveWaitUntil(ctx, &param))
	case "textDocument/completion":
		var param CompletionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCompletion(ctx, &param))
	case "completionItem/resolve":
		var param CompletionItem
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CompletionItemResolve(ctx, &param))
	case "textDocument/hover":
		var param HoverParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentHover(ctx, &param))
	case "textDocument/signatureHelp":
		var param SignatureHelpParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSignatureHelp(ctx, &param))
	case "textDocument/declaration":
		var param DeclarationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDeclaration(ctx, &param))
	case "textDocument/definition":
		var param DefinitionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDefinition(ctx, &param))
	case "textDocument/typeDefinition":
		var param TypeDefinitionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentTypeDefinition(ctx, &param))
	case "textDocument/implementation":
		var param ImplementationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentImplementation(ctx, &param))
	case "textDocument/references":
		var param ReferenceParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentReferences(ctx, &param))
	case "textDocument/documentHighlight":
		var param DocumentHighlightParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentHighlight(ctx, &param))
	case "textDocument/documentSymbol":
		var param DocumentSymbolParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDocumentSymbol(ctx, &param))
	case "textDocument/codeAction":
		var param CodeActionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCodeAction(ctx, &param))
	case "codeAction/resolve":
		var param CodeAction
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CodeActionResolve(ctx, &param))
	case "textDocument/codeLens":
		var param CodeLensParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCodeLens(ctx, &param))
	case "codeLens/resolve":
		var param CodeLens
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CodeLensResolve(ctx, &param))
	case "textDocument/documentLink":
		var param DocumentLinkParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentLink(ctx, &param))
	case "documentLink/resolve":
		var param DocumentLink
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.DocumentLinkResolve(ctx, &param))
	case "textDocument/documentColor":
		var param DocumentColorParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentColor(ctx, &param))
	case "textDocument/colorPresentation":
		var param ColorPresentationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentColorPresentation(ctx, &param))
	case "textDocument/formatting":
		var param DocumentFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentFormatting(ctx, &param))
	case "textDocument/rangeFormatting":
		var param DocumentRangeFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentRangeFormatting(ctx, &param))
	case "textDocument/onTypeFormatting":
		var param DocumentOnTypeFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentOnTypeFormatting(ctx, &param))
	case "textDocument/rename":
		var param RenameParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentRename(ctx, &param))
	case "textDocument/prepareRename":
		var param PrepareRenameParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		// resp(serv.handler.TextDocumentPrepareRename(ctx, &param))
		panic("unimplemented")
	case "textDocument/foldingRange":
		var param FoldingRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentFoldingRange(ctx, &param))
	case "textDocument/selectionRange":
		var param SelectionRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSelectionRange(ctx, &param))
	case "textDocument/prepareCallHierarchy":
		var param CallHierarchyPrepareParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentPrepareCallHierarchy(ctx, &param))
	case "callHierarchy/incomingCalls":
		var param CallHierarchyIncomingCallsParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CallHierarchyIncomingCalls(ctx, &param))
	case "callHierarchy/outgoingCalls":
		var param CallHierarchyOutgoingCallsParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CallHierarchyOutgoingCalls(ctx, &param))
	case "textDocument/semanticTokens/full":
		var param SemanticTokensParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSemanticTokensFull(ctx, &param))
	case "textDocument/semanticTokens/full/delta":
		var param SemanticTokensDeltaParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentSemanticTokensFullDelta(ctx, &param))
	case "textDocument/semanticTokens/range":
		var param SemanticTokensRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSemanticTokensRange(ctx, &param))
	case "workspace/semanticTokens/refresh":
		resp(nil, serv.handler.WorkspaceSemanticTokensRefresh(ctx))
	case "textDocument/linkedEditingRange":
		var param LinkedEditingRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentLinkedEditingRange(ctx, &param))
	case "textDocument/moniker":
		var param MonikerParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentMoniker(ctx, &param))
	default:
		panic("unimplemented message")
	}
}

// Requests to Client

func (serv *Server) WindowShowMessageRequest(ctx context.Context, param *ShowMessageRequestParams) (*MessageActionItem, *jsonrpc.ResponseError, error) {
	resp, respErr, err := serv.conn.SendRequest(ctx, "window/showMessageRequest", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: the selected MessageActionItem | null if none got selected
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res MessageActionItem
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (serv *Server) WindowShowDocument(ctx context.Context, param *ShowDocumentParams) (*ShowDocumentResult, *jsonrpc.ResponseError, error) {
	resp, respErr, err := serv.conn.SendRequest(ctx, "window/showDocument", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res ShowDocumentResult
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (serv *Server) WindowWorkDoneProgressCreate(ctx context.Context, param *WorkDoneProgressCreateParams) (*jsonrpc.ResponseError, error) {
	_, respErr, err := serv.conn.SendRequest(ctx, "window/workDoneProgress/create", EncodeMessage(param))
	return respErr, err
}

func (serv *Server) ClientRegisterCapability(ctx context.Context, param *RegistrationParams) (*jsonrpc.ResponseError, error) {
	_, respErr, err := serv.conn.SendRequest(ctx, "client/registerCapability", EncodeMessage(param))
	return respErr, err
}

func (serv *Server) ClientUnregisterCapability(ctx context.Context, param *UnregistrationParams) (*jsonrpc.ResponseError, error) {
	_, respErr, err := serv.conn.SendRequest(ctx, "client/unregisterCapability", EncodeMessage(param))
	return respErr, err
}

func (serv *Server) WorkspaceWorkspaceFolders(ctx context.Context) ([]WorkspaceFolder, *jsonrpc.ResponseError, error) {
	resp, respErr, err := serv.conn.SendRequest(ctx, "workspace/workspaceFolders", EncodeMessage(jsonrpc.NullResult))
	if err != nil {
		return nil, nil, err
	}
	// result: WorkspaceFolder[] | null
	if string(resp) == "null" {
		return nil, respErr, nil
	}
	var res []WorkspaceFolder
	return res, respErr, json.Unmarshal(resp, &res)
}

func (serv *Server) WorkspaceConfiguration(ctx context.Context, param *ConfigurationParams) ([]json.RawMessage, *jsonrpc.ResponseError, error) {
	resp, respErr, err := serv.conn.SendRequest(ctx, "workspace/configuration", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	// result: any[]
	var res []json.RawMessage
	return res, respErr, json.Unmarshal(resp, &res)
}

func (serv *Server) WorkspaceApplyEdit(ctx context.Context, param *ApplyWorkspaceEditParams) (*ApplyWorkspaceEditResult, *jsonrpc.ResponseError, error) {
	resp, respErr, err := serv.conn.SendRequest(ctx, "workspace/applyEdit", EncodeMessage(param))
	if err != nil {
		return nil, nil, err
	}
	var res ApplyWorkspaceEditResult
	return &res, respErr, json.Unmarshal(resp, &res)
}

func (serv *Server) WorkspaceCodeLensRefresh(ctx context.Context) (*jsonrpc.ResponseError, error) {
	_, respErr, err := serv.conn.SendRequest(ctx, "workspace/codeLens/refresh", EncodeMessage(jsonrpc.NullResult))
	return respErr, err
}

// Notifications to Client

func (serv *Server) Progress(param *ProgressParams) error {
	return serv.conn.SendNotification("$/progress", EncodeMessage(param))
}

func (serv *Server) LogTrace(param *LogTraceParams) error {
	return serv.conn.SendNotification("&/logTrace", EncodeMessage(param))
}

func (serv *Server) WindowShowMessage(param *ShowMessageParams) error {
	return serv.conn.SendNotification("window/showMessage", EncodeMessage(param))
}

func (serv *Server) WindowLogMessage(param *LogMessageParams) error {
	return serv.conn.SendNotification("window/logMessage", EncodeMessage(param))
}

func (serv *Server) TelemetryEvent(param json.RawMessage) error {
	return serv.conn.SendNotification("telemetry/event", param)
}

func (serv *Server) TextDocumentPublishDiagnostics(param *PublishDiagnosticsParams) error {
	return serv.conn.SendNotification("textDocument/publishDiagnostics", EncodeMessage(param))
}
