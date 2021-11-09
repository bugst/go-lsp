//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestMarshalUnmarshalWithBooleanSumType(t *testing.T) {
	var x ServerCapabilities
	err := json.Unmarshal([]byte(`
	{
		"declarationProvider":true,
		"hoverProvider":true,
		"textDocumentSync":{
			"save":true
		}
	}`), &x)
	require.NoError(t, err)
	require.Equal(t, "&{DeclarationOptions:<nil> StaticRegistrationOptions:<nil> TextDocumentRegistrationOptions:<nil>}", fmt.Sprintf("%+v", x.DeclarationProvider))
	require.Equal(t, "&{WorkDoneProgressOptions:<nil>}", fmt.Sprintf("%+v", x.HoverProvider))
	require.Equal(t, "&{IncludeText:false}", fmt.Sprintf("%+v", x.TextDocumentSync.Save))
	err = json.Unmarshal([]byte(`
	{
		"declarationProvider":{},
		"hoverProvider":{},
		"textDocumentSync":{
			"save":{}
		}
	}`), &x)
	require.NoError(t, err)
	require.Equal(t, "&{DeclarationOptions:<nil> StaticRegistrationOptions:<nil> TextDocumentRegistrationOptions:<nil>}", fmt.Sprintf("%+v", x.DeclarationProvider))
	require.Equal(t, "&{WorkDoneProgressOptions:<nil>}", fmt.Sprintf("%+v", x.HoverProvider))
	require.Equal(t, "&{IncludeText:false}", fmt.Sprintf("%+v", x.TextDocumentSync.Save))
	y := ServerCapabilities{
		DeclarationProvider: &DeclarationRegistrationOptions{
			DeclarationOptions: &DeclarationOptions{
				WorkDoneProgressOptions: &WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
			StaticRegistrationOptions: &StaticRegistrationOptions{
				ID: "123",
			},
			TextDocumentRegistrationOptions: &TextDocumentRegistrationOptions{
				DocumentSelector: &DocumentSelector{
					DocumentFilter{Language: "lan", Scheme: "sch", Pattern: "patt"},
					DocumentFilter{Language: "lang2"},
				},
			},
		},
	}
	d, err := json.MarshalIndent(y, "", "  ")
	require.NoError(t, err)
	require.Equal(t, "{\n  \"declarationProvider\": {\n    \"workDoneProgress\": true,\n    \"id\": \"123\",\n    \"documentSelector\": [\n      {\n        \"language\": \"lan\",\n        \"scheme\": \"sch\",\n        \"pattern\": \"patt\"\n      },\n      {\n        \"language\": \"lang2\"\n      }\n    ]\n  }\n}", string(d))
	y = ServerCapabilities{
		DeclarationProvider: &DeclarationRegistrationOptions{
			DeclarationOptions: &DeclarationOptions{
				WorkDoneProgressOptions: &WorkDoneProgressOptions{
					WorkDoneProgress: true,
				},
			},
		},
	}
	d, err = json.MarshalIndent(y, "", "  ")
	require.NoError(t, err)
	require.Equal(t, "{\n  \"declarationProvider\": {\n    \"workDoneProgress\": true\n  }\n}", string(d))
}

func TestInitializeResult(t *testing.T) {
	clangd11ServerCapabilities := []byte(`
	{
	  "capabilities": {
		"astProvider": true,
		"callHierarchyProvider": true,
		"codeActionProvider": {
		  "codeActionKinds": [
			"quickfix",
			"refactor",
			"info"
		  ]
		},
		"compilationDatabase": {
		  "automaticReload": true
		},
		"completionProvider": {
		  "allCommitCharacters": [ " ", "\t", "(", ")", "[", "]", "{", "}", "<", ">", ":", ";", ",", "+", "-", "/", "*", "%", "^", "&", "#", "?", ".", "=", "\"", "'", "|" ],
		  "resolveProvider": false,
		  "triggerCharacters": [ ".", "<", ">", ":", "\"", "/" ]
		},
		"declarationProvider": true,
		"definitionProvider": true,
		"documentFormattingProvider": true,
		"documentHighlightProvider": true,
		"documentLinkProvider": {
		  "resolveProvider": false
		},
		"documentOnTypeFormattingProvider": {
		  "firstTriggerCharacter": "\n",
		  "moreTriggerCharacter": []
		},
		"documentRangeFormattingProvider": true,
		"documentSymbolProvider": true,
		"executeCommandProvider": {
		  "commands": [
			"clangd.applyFix",
			"clangd.applyTweak"
		  ]
		},
		"hoverProvider": true,
		"implementationProvider": true,
		"memoryUsageProvider": true,
		"referencesProvider": true,
		"renameProvider": {
		  "prepareProvider": true
		},
		"selectionRangeProvider": true,
		"semanticTokensProvider": {
		  "full": {
			"delta": true
		  },
		  "legend": {
			"tokenModifiers": [],
			"tokenTypes": [
			  "variable", "variable", "parameter", "function", "method",
			  "function", "property", "variable", "class", "enum",
			  "enumMember", "type", "dependent", "dependent", "namespace",
			  "typeParameter", "concept", "type", "macro", "comment"
			]
		  },
		  "range": false
		},
		"signatureHelpProvider": {
		  "triggerCharacters": [ "(", "," ]
		},
		"textDocumentSync": {
		  "change": 2,
		  "openClose": true,
		  "save": true
		},
		"typeHierarchyProvider": true,
		"workspaceSymbolProvider": true
	  },
	  "serverInfo": {
		"name": "clangd",
		"version": "clangd version 12.0.0 (https://github.com/llvm/llvm-project e841bd5f335864b8c4d81cbf4df08460ef39f2ae)"
      }
	}`)
	var sc InitializeResult
	err := json.Unmarshal(clangd11ServerCapabilities, &sc)
	require.NoError(t, err)
	fmt.Println(sc.Capabilities.SemanticTokensProvider)

	_, err = json.MarshalIndent(&InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: &TextDocumentSyncOptions{
				OpenClose: true,
			}, //{Kind: &TDSKIncremental},
			HoverProvider: &HoverOptions{}, // true,
			CompletionProvider: &CompletionOptions{
				TriggerCharacters: []string{".", "\u003e", ":"},
			},
			SignatureHelpProvider: &SignatureHelpOptions{
				TriggerCharacters: []string{"(", ","},
			},
			DefinitionProvider: &DefinitionOptions{}, // true,
			// ReferencesProvider:              &ReferenceOptions{},  // TODO: true
			DocumentHighlightProvider:       &DocumentHighlightOptions{}, //true,
			DocumentSymbolProvider:          &DocumentSymbolOptions{},    //true,
			WorkspaceSymbolProvider:         &WorkspaceSymbolOptions{},   //true,
			CodeActionProvider:              &CodeActionOptions{ResolveProvider: true},
			DocumentFormattingProvider:      &DocumentFormattingOptions{},      //true,
			DocumentRangeFormattingProvider: &DocumentRangeFormattingOptions{}, //true,
			DocumentOnTypeFormattingProvider: &DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "\n",
			},
			RenameProvider: &RenameOptions{PrepareProvider: false}, // TODO: true
			ExecuteCommandProvider: &ExecuteCommandOptions{
				Commands: []string{"clangd.applyFix", "clangd.applyTweak"},
			},
		},
	}, "", "  ")
	require.NoError(t, err)
}
