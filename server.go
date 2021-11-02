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

	Initialize(context.Context, jsonrpc.FunctionLogger, *InitializeParams) (*InitializeResult, *jsonrpc.ResponseError)
	Shutdown(context.Context, jsonrpc.FunctionLogger) *jsonrpc.ResponseError
	WorkspaceSymbol(context.Context, jsonrpc.FunctionLogger, *WorkspaceSymbolParams) ([]SymbolInformation, *jsonrpc.ResponseError)
	WorkspaceExecuteCommand(context.Context, jsonrpc.FunctionLogger, *ExecuteCommandParams) (json.RawMessage, *jsonrpc.ResponseError)
	WorkspaceWillCreateFiles(context.Context, jsonrpc.FunctionLogger, *CreateFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	WorkspaceWillRenameFiles(context.Context, jsonrpc.FunctionLogger, *RenameFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	WorkspaceWillDeleteFiles(context.Context, jsonrpc.FunctionLogger, *DeleteFilesParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	TextDocumentWillSaveWaitUntil(context.Context, jsonrpc.FunctionLogger, *WillSaveTextDocumentParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentCompletion(context.Context, jsonrpc.FunctionLogger, *CompletionParams) (*CompletionList, *jsonrpc.ResponseError)
	CompletionItemResolve(context.Context, jsonrpc.FunctionLogger, *CompletionItem) (*CompletionItem, *jsonrpc.ResponseError)
	TextDocumentHover(context.Context, jsonrpc.FunctionLogger, *HoverParams) (*Hover, *jsonrpc.ResponseError)
	TextDocumentSignatureHelp(context.Context, jsonrpc.FunctionLogger, *SignatureHelpParams) (*SignatureHelp, *jsonrpc.ResponseError)
	TextDocumentDeclaration(context.Context, jsonrpc.FunctionLogger, *DeclarationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentDefinition(context.Context, jsonrpc.FunctionLogger, *DefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentTypeDefinition(context.Context, jsonrpc.FunctionLogger, *TypeDefinitionParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentImplementation(context.Context, jsonrpc.FunctionLogger, *ImplementationParams) ([]Location, []LocationLink, *jsonrpc.ResponseError)
	TextDocumentReferences(context.Context, jsonrpc.FunctionLogger, *ReferenceParams) ([]Location, *jsonrpc.ResponseError)
	TextDocumentDocumentHighlight(context.Context, jsonrpc.FunctionLogger, *DocumentHighlightParams) ([]DocumentHighlight, *jsonrpc.ResponseError)
	TextDocumentDocumentSymbol(context.Context, jsonrpc.FunctionLogger, *DocumentSymbolParams) ([]DocumentSymbol, []SymbolInformation, *jsonrpc.ResponseError)
	TextDocumentCodeAction(context.Context, jsonrpc.FunctionLogger, *CodeActionParams) ([]CommandOrCodeAction, *jsonrpc.ResponseError)
	CodeActionResolve(context.Context, jsonrpc.FunctionLogger, *CodeAction) (*CodeAction, *jsonrpc.ResponseError)
	TextDocumentCodeLens(context.Context, jsonrpc.FunctionLogger, *CodeLensParams) ([]CodeLens, *jsonrpc.ResponseError)
	CodeLensResolve(context.Context, jsonrpc.FunctionLogger, *CodeLens) (*CodeLens, *jsonrpc.ResponseError)
	TextDocumentDocumentLink(context.Context, jsonrpc.FunctionLogger, *DocumentLinkParams) ([]DocumentLink, *jsonrpc.ResponseError)
	DocumentLinkResolve(context.Context, jsonrpc.FunctionLogger, *DocumentLink) (*DocumentLink, *jsonrpc.ResponseError)
	TextDocumentDocumentColor(context.Context, jsonrpc.FunctionLogger, *DocumentColorParams) ([]ColorInformation, *jsonrpc.ResponseError)
	TextDocumentColorPresentation(context.Context, jsonrpc.FunctionLogger, *ColorPresentationParams) ([]ColorPresentation, *jsonrpc.ResponseError)
	TextDocumentFormatting(context.Context, jsonrpc.FunctionLogger, *DocumentFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentRangeFormatting(context.Context, jsonrpc.FunctionLogger, *DocumentRangeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentOnTypeFormatting(context.Context, jsonrpc.FunctionLogger, *DocumentOnTypeFormattingParams) ([]TextEdit, *jsonrpc.ResponseError)
	TextDocumentRename(context.Context, jsonrpc.FunctionLogger, *RenameParams) (*WorkspaceEdit, *jsonrpc.ResponseError)
	//TextDocumentPrepareRename(context.Context,jsonrpc.FunctionLogger, *PrepareRenameParams) (???, *jsonrpc.ResponseError)
	TextDocumentFoldingRange(context.Context, jsonrpc.FunctionLogger, *FoldingRangeParams) ([]FoldingRange, *jsonrpc.ResponseError)
	TextDocumentSelectionRange(context.Context, jsonrpc.FunctionLogger, *SelectionRangeParams) ([]SelectionRange, *jsonrpc.ResponseError)
	TextDocumentPrepareCallHierarchy(context.Context, jsonrpc.FunctionLogger, *CallHierarchyPrepareParams) ([]CallHierarchyItem, *jsonrpc.ResponseError)
	CallHierarchyIncomingCalls(context.Context, jsonrpc.FunctionLogger, *CallHierarchyIncomingCallsParams) ([]CallHierarchyIncomingCall, *jsonrpc.ResponseError)
	CallHierarchyOutgoingCalls(context.Context, jsonrpc.FunctionLogger, *CallHierarchyOutgoingCallsParams) ([]CallHierarchyOutgoingCall, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensFull(context.Context, jsonrpc.FunctionLogger, *SemanticTokensParams) (*SemanticTokens, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensFullDelta(context.Context, jsonrpc.FunctionLogger, *SemanticTokensDeltaParams) (*SemanticTokens, *SemanticTokensDelta, *jsonrpc.ResponseError)
	TextDocumentSemanticTokensRange(context.Context, jsonrpc.FunctionLogger, *SemanticTokensRangeParams) (*SemanticTokens, *jsonrpc.ResponseError)
	WorkspaceSemanticTokensRefresh(context.Context, jsonrpc.FunctionLogger) *jsonrpc.ResponseError
	TextDocumentLinkedEditingRange(context.Context, jsonrpc.FunctionLogger, *LinkedEditingRangeParams) (*LinkedEditingRanges, *jsonrpc.ResponseError)
	TextDocumentMoniker(context.Context, jsonrpc.FunctionLogger, *MonikerParams) ([]Moniker, *jsonrpc.ResponseError)

	// Notifications ->

	Progress(jsonrpc.FunctionLogger, *ProgressParams)
	// CancelRequrest(*jsonrpc.CancelParams) - automatically handled by the rpc library
	Initialized(jsonrpc.FunctionLogger, *InitializedParams)
	Exit(jsonrpc.FunctionLogger)
	SetTrace(jsonrpc.FunctionLogger, *SetTraceParams)
	WindowWorkDoneProgressCancel(jsonrpc.FunctionLogger, *WorkDoneProgressCancelParams)
	WorkspaceDidChangeWorkspaceFolders(jsonrpc.FunctionLogger, *DidChangeWorkspaceFoldersParams)
	WorkspaceDidChangeConfiguration(jsonrpc.FunctionLogger, *DidChangeConfigurationParams)
	WorkspaceDidChangeWatchedFiles(jsonrpc.FunctionLogger, *DidChangeWatchedFilesParams)
	WorkspaceDidCreateFiles(jsonrpc.FunctionLogger, *CreateFilesParams)
	WorkspaceDidRenameFiles(jsonrpc.FunctionLogger, *RenameFilesParams)
	WorkspaceDidDeleteFiles(jsonrpc.FunctionLogger, *DeleteFilesParams)
	TextDocumentDidOpen(jsonrpc.FunctionLogger, *DidOpenTextDocumentParams)
	TextDocumentDidChange(jsonrpc.FunctionLogger, *DidChangeTextDocumentParams)
	TextDocumentWillSave(jsonrpc.FunctionLogger, *WillSaveTextDocumentParams)
	TextDocumentDidSave(jsonrpc.FunctionLogger, *DidSaveTextDocumentParams)
	TextDocumentDidClose(jsonrpc.FunctionLogger, *DidCloseTextDocumentParams)
}

// Server is an LSP Server
type Server struct {
	conn         *jsonrpc.Connection
	handler      ClientMessagesHandler
	errorHandler func(e error)
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

func (serv *Server) SetLogger(l jsonrpc.Logger) {
	serv.conn.SetLogger(l)
}

func (serv *Server) SetErrorHandler(handler func(e error)) {
	serv.errorHandler = handler
}

func (serv *Server) Run() {
	serv.conn.Run()
}

func (serv *Server) notificationDispatcher(logger jsonrpc.FunctionLogger, method string, req json.RawMessage) {
	switch method {
	case "$/progress":
		var param ProgressParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.Progress(logger, &param)
	case "$/cancelRequrest":
		panic("should not reach here")
	case "initialized":
		var param InitializedParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.Initialized(logger, &param)
	case "exit":
		serv.handler.Exit(logger)
	case "$/setTrace":
		var param SetTraceParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.SetTrace(logger, &param)
	case "window/workDoneProgress/cancel":
		var param WorkDoneProgressCancelParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WindowWorkDoneProgressCancel(logger, &param)
	case "workspace/didChangeWorkspaceFolders":
		var param DidChangeWorkspaceFoldersParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeWorkspaceFolders(logger, &param)
	case "workspace/didChangeConfiguration":
		var param DidChangeConfigurationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeConfiguration(logger, &param)
	case "workspace/didChangeWatchedFiles":
		var param DidChangeWatchedFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidChangeWatchedFiles(logger, &param)
	case "workspace/didCreateFiles":
		var param CreateFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidCreateFiles(logger, &param)
	case "workspace/didRenameFiles":
		var param RenameFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidRenameFiles(logger, &param)
	case "workspace/didDeleteFiles":
		var param DeleteFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.WorkspaceDidDeleteFiles(logger, &param)
	case "textDocument/didOpen":
		var param DidOpenTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidOpen(logger, &param)
	case "textDocument/didChange":
		var param DidChangeTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidChange(logger, &param)
	case "textDocument/willSave":
		var param WillSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentWillSave(logger, &param)
	case "textDocument/didSave":
		var param DidSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidSave(logger, &param)
	case "textDocument/didClose":
		var param DidCloseTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		serv.handler.TextDocumentDidClose(logger, &param)
	default:
		panic("unimplemented message")
	}
}

func (serv *Server) requestDispatcher(ctx context.Context, logger jsonrpc.FunctionLogger, method string, req json.RawMessage, respCallback func(json.RawMessage, *jsonrpc.ResponseError)) {
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
		resp(serv.handler.Initialize(ctx, logger, &param))
	case "shutdown":
		resp(nil, serv.handler.Shutdown(ctx, logger))
	case "workspace/symbol":
		var param WorkspaceSymbolParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceSymbol(ctx, logger, &param))
	case "workspace/executeCommand":
		var param ExecuteCommandParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceExecuteCommand(ctx, logger, &param))
	case "workspace/willCreateFiles":
		var param CreateFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillCreateFiles(ctx, logger, &param))
	case "workspace/willRenameFiles":
		var param RenameFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillRenameFiles(ctx, logger, &param))
	case "workspace/willDeleteFiles":
		var param DeleteFilesParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.WorkspaceWillDeleteFiles(ctx, logger, &param))
	case "textDocument/willSaveWaitUntil":
		var param WillSaveTextDocumentParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentWillSaveWaitUntil(ctx, logger, &param))
	case "textDocument/completion":
		var param CompletionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCompletion(ctx, logger, &param))
	case "completionItem/resolve":
		var param CompletionItem
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CompletionItemResolve(ctx, logger, &param))
	case "textDocument/hover":
		var param HoverParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentHover(ctx, logger, &param))
	case "textDocument/signatureHelp":
		var param SignatureHelpParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSignatureHelp(ctx, logger, &param))
	case "textDocument/declaration":
		var param DeclarationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDeclaration(ctx, logger, &param))
	case "textDocument/definition":
		var param DefinitionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDefinition(ctx, logger, &param))
	case "textDocument/typeDefinition":
		var param TypeDefinitionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentTypeDefinition(ctx, logger, &param))
	case "textDocument/implementation":
		var param ImplementationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentImplementation(ctx, logger, &param))
	case "textDocument/references":
		var param ReferenceParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentReferences(ctx, logger, &param))
	case "textDocument/documentHighlight":
		var param DocumentHighlightParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentHighlight(ctx, logger, &param))
	case "textDocument/documentSymbol":
		var param DocumentSymbolParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentDocumentSymbol(ctx, logger, &param))
	case "textDocument/codeAction":
		var param CodeActionParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCodeAction(ctx, logger, &param))
	case "codeAction/resolve":
		var param CodeAction
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CodeActionResolve(ctx, logger, &param))
	case "textDocument/codeLens":
		var param CodeLensParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentCodeLens(ctx, logger, &param))
	case "codeLens/resolve":
		var param CodeLens
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CodeLensResolve(ctx, logger, &param))
	case "textDocument/documentLink":
		var param DocumentLinkParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentLink(ctx, logger, &param))
	case "documentLink/resolve":
		var param DocumentLink
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.DocumentLinkResolve(ctx, logger, &param))
	case "textDocument/documentColor":
		var param DocumentColorParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentDocumentColor(ctx, logger, &param))
	case "textDocument/colorPresentation":
		var param ColorPresentationParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentColorPresentation(ctx, logger, &param))
	case "textDocument/formatting":
		var param DocumentFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentFormatting(ctx, logger, &param))
	case "textDocument/rangeFormatting":
		var param DocumentRangeFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentRangeFormatting(ctx, logger, &param))
	case "textDocument/onTypeFormatting":
		var param DocumentOnTypeFormattingParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentOnTypeFormatting(ctx, logger, &param))
	case "textDocument/rename":
		var param RenameParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentRename(ctx, logger, &param))
	case "textDocument/prepareRename":
		var param PrepareRenameParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		// resp(serv.handler.TextDocumentPrepareRename(ctx,logger, &param))
		panic("unimplemented")
	case "textDocument/foldingRange":
		var param FoldingRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentFoldingRange(ctx, logger, &param))
	case "textDocument/selectionRange":
		var param SelectionRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSelectionRange(ctx, logger, &param))
	case "textDocument/prepareCallHierarchy":
		var param CallHierarchyPrepareParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentPrepareCallHierarchy(ctx, logger, &param))
	case "callHierarchy/incomingCalls":
		var param CallHierarchyIncomingCallsParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CallHierarchyIncomingCalls(ctx, logger, &param))
	case "callHierarchy/outgoingCalls":
		var param CallHierarchyOutgoingCallsParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.CallHierarchyOutgoingCalls(ctx, logger, &param))
	case "textDocument/semanticTokens/full":
		var param SemanticTokensParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSemanticTokensFull(ctx, logger, &param))
	case "textDocument/semanticTokens/full/delta":
		var param SemanticTokensDeltaParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp2(serv.handler.TextDocumentSemanticTokensFullDelta(ctx, logger, &param))
	case "textDocument/semanticTokens/range":
		var param SemanticTokensRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentSemanticTokensRange(ctx, logger, &param))
	case "workspace/semanticTokens/refresh":
		resp(nil, serv.handler.WorkspaceSemanticTokensRefresh(ctx, logger))
	case "textDocument/linkedEditingRange":
		var param LinkedEditingRangeParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentLinkedEditingRange(ctx, logger, &param))
	case "textDocument/moniker":
		var param MonikerParams
		if err := json.Unmarshal(req, &param); err != nil {
			serv.errorHandler(err)
			return
		}
		resp(serv.handler.TextDocumentMoniker(ctx, logger, &param))
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
