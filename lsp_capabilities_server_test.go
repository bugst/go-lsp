package lsp

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
