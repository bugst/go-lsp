package lsp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"go.bug.st/json"
	"go.bug.st/lsp/jsonrpc"
)

var sketchPath = paths.New("/home/user/sketch")
var initMsg = `{
	"jsonrpc": "2.0",
	"id": 0,
	"method": "initialize",
	"params": {
	  "processId": 32346,
	  "clientInfo": {"name": "vscode","version": "1.44.0"},
	  "rootPath": "` + sketchPath.String() + `",
	  "rootUri": "` + NewDocumentURIFromPath(sketchPath).String() + `",
	  "capabilities": {
		"workspace": {
		  "applyEdit": true,
		  "workspaceEdit": {
			"documentChanges": true,
			"resourceOperations": ["create","rename","delete"],
			"failureHandling": "textOnlyTransactional"
		  },
		  "didChangeConfiguration": {"dynamicRegistration": true},
		  "didChangeWatchedFiles": {"dynamicRegistration": true},
		  "symbol": {
			"dynamicRegistration": true,
			"symbolKind": {"valueSet": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]}
		  },
		  "executeCommand": {"dynamicRegistration": true},
		  "configuration": true,
		  "workspaceFolders": true
		},
		"textDocument": {
		  "publishDiagnostics": {
			"relatedInformation": true,
			"versionSupport": false,
			"tagSupport": {"valueSet": [1,2]}
		  },
		  "synchronization": {
			"dynamicRegistration": true,
			"willSave": true,
			"willSaveWaitUntil": true,
			"didSave": true
		  },
		  "completion": {
			"dynamicRegistration": true,
			"contextSupport": true,
			"completionItem": {
			  "snippetSupport": true,
			  "commitCharactersSupport": true,
			  "documentationFormat": ["markdown","plaintext"],
			  "deprecatedSupport": true,
			  "preselectSupport": true,
			  "tagSupport": {"valueSet": [1]}
			},
			"completionItemKind": {"valueSet": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]}
		  },
		  "hover": {
			"dynamicRegistration": true,
			"contentFormat": ["markdown","plaintext"]
		  },
		  "signatureHelp": {
			"dynamicRegistration": true,
			"signatureInformation": {
			  "documentationFormat": ["markdown","plaintext"],
			  "parameterInformation": {"labelOffsetSupport": true}
			},
			"contextSupport": true
		  },
		  "definition": {"dynamicRegistration": true,"linkSupport": true},
		  "references": {"dynamicRegistration": true},
		  "documentHighlight": {"dynamicRegistration": true},
		  "documentSymbol": {
			"dynamicRegistration": true,
			"symbolKind": {"valueSet": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]},
			"hierarchicalDocumentSymbolSupport": true
		  },
		  "codeAction": {
			"dynamicRegistration": true,
			"isPreferredSupport": true,
			"codeActionLiteralSupport": {
			  "codeActionKind": {
				"valueSet": ["","quickfix","refactor","refactor.extract","refactor.inline","refactor.rewrite","source","source.organizeImports"]
			  }
			}
		  },
		  "codeLens": {"dynamicRegistration": true},
		  "formatting": {"dynamicRegistration": true},
		  "rangeFormatting": {"dynamicRegistration": true},
		  "onTypeFormatting": {"dynamicRegistration": true},
		  "rename": {"dynamicRegistration": true,"prepareSupport": true},
		  "documentLink": {"dynamicRegistration": true,"tooltipSupport": true},
		  "typeDefinition": {"dynamicRegistration": true,"linkSupport": true},
		  "implementation": {"dynamicRegistration": true,"linkSupport": true},
		  "colorProvider": {"dynamicRegistration": true},
		  "foldingRange": {"dynamicRegistration": true,"rangeLimit": 5000,"lineFoldingOnly": true},
		  "declaration": {"dynamicRegistration": true,"linkSupport": true},
		  "selectionRange": {"dynamicRegistration": true}
		},
		"window": {"workDoneProgress": true}
	  },
	  "initializationOptions": {},
	  "trace": "off",
	  "workspaceFolders": [
		{
		  "uri": "` + NewDocumentURIFromPath(sketchPath).String() + `",
		  "name": "` + sketchPath.Base() + `"
		}
	  ]
	}
  }`

func TestLSPInitializeMessages(t *testing.T) {
	testdata := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(initMsg), initMsg)
	resp := ""
	output := &bytes.Buffer{}
	var wg sync.WaitGroup
	conn := jsonrpc.NewConnection(
		bufio.NewReader(strings.NewReader(testdata)),
		output,
		func(ctx context.Context, logger jsonrpc.FunctionLogger, method string, params json.RawMessage, respCallback func(json.RawMessage, *jsonrpc.ResponseError)) {
			require.Equal(t, "initialize", method)
			var par InitializeParams
			err := json.Unmarshal(params, &par)
			require.NoError(t, err)
			spew.Config.DisableCapacities = true
			spew.Config.DisableMethods = true
			spew.Config.DisablePointerAddresses = true
			spew.Dump(par)
			respCallback(nil, nil)
		},
		func(logger jsonrpc.FunctionLogger, method string, params json.RawMessage) {
			resp += fmt.Sprintf("NOT method=%v params=%v\n", method, params)
		},
		func(e error) {
			if e == io.EOF {
				return
			}
			resp += fmt.Sprintf("error=%s\n", e)
		},
	)
	conn.Run() // Exits when input is fully consumed
	wg.Wait()  // Wait for all pending responses to get through
	conn.Close()
	// require.Equal(t,
	// 	"NOT method=initialized params=[[49 50 51]]\n"+
	// 		"REQ method=textDocument/didOpen params=[123 10 9 125]\n"+
	// 		"REQ method=textDocument/didClose params=[123 10 9 125]\n"+
	// 		"REQ method=tocancel params=[123 10 9 125]\n"+
	// 		"", resp)

	//require.Equal(t, "", output.String())
	fmt.Println(output.String())
}
