//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestSumTypes(t *testing.T) {
	comJSON := "{\"title\":\"command_title\",\"command\":\"command\"}"
	caJSON := "{\"title\":\"codeaction_title\",\"kind\":\"quickfix\",\"isPreferred\":true,\"command\":{\"title\":\"command_title\",\"command\":\"command\"}}"

	{
		var c CommandOrCodeAction
		c.Set(Command{
			Title:   "command_title",
			Command: "command",
		})
		data, err := json.Marshal(c)
		require.NoError(t, err)
		require.Equal(t, comJSON, string(data))
	}
	{
		var c CommandOrCodeAction
		c.Set(CodeAction{
			Title:       "codeaction_title",
			Kind:        CodeActionKindQuickFix,
			IsPreferred: true,
			Command: &Command{
				Title:   "command_title",
				Command: "command",
			},
		})
		data, err := json.Marshal(c)
		require.NoError(t, err)
		require.Equal(t, caJSON, string(data))
	}

	{
		var c CommandOrCodeAction
		err := json.Unmarshal([]byte(comJSON), &c)
		require.NoError(t, err)
		res := c.Get()
		require.IsType(t, Command{}, res)
		require.Equal(t, "command", res.(Command).Command)
	}
	{
		var c CommandOrCodeAction
		err := json.Unmarshal([]byte(caJSON), &c)
		require.NoError(t, err)
		res := c.Get()
		require.IsType(t, CodeAction{}, res)
		require.Equal(t, &Command{Title: "command_title", Command: "command"}, res.(CodeAction).Command)
	}

	// Let's try an array of CommandOrCodeActions...
	{
		jsonIn := json.RawMessage("[" + caJSON + "," + comJSON + "," + caJSON + "," + caJSON + "," + comJSON + "]")
		res, err := DecodeServerResponseResult("textDocument/codeAction", jsonIn)
		require.NoError(t, err)
		require.IsType(t, []CommandOrCodeAction{}, res)
		resArray := res.([]CommandOrCodeAction)
		require.IsType(t, CodeAction{}, resArray[0].Get())
		require.IsType(t, Command{}, resArray[1].Get())
		require.IsType(t, CodeAction{}, resArray[2].Get())
		require.IsType(t, CodeAction{}, resArray[3].Get())
		require.IsType(t, Command{}, resArray[4].Get())

		data, err := json.Marshal(resArray)
		require.NoError(t, err)
		require.Equal(t, string(jsonIn), string(data))
	}

	// some real-world examples
	{
		jsonIn := json.RawMessage(`
		[
			{
				"diagnostics": [
					{
						"code":"undeclared_var_use_suggest",
						"message":"Use of undeclared identifier 'ads'; did you mean 'abs'? (fix available)",
						"range": {
							"end":  {"character":5, "line":14},
							"start":{"character":2, "line":14}
						},
						"severity":1,
						"source":"clang"
					}
				],
				"edit": {
					"changes": {
						"file:///tmp/arduino-language-server616865191/sketch/Blink.ino.cpp": [
							{
								"newText":"abs",
								"range": {
									"end":  {"character":5, "line":14},
									"start":{"character":2, "line":14}
								}
							}
						]
					}
				},
				"isPreferred":true,
				"kind":"quickfix",
				"title":"change 'ads' to 'abs'"
			}
		]`)
		res, err := DecodeServerResponseResult("textDocument/codeAction", jsonIn)
		require.NoError(t, err)
		require.IsType(t, []CommandOrCodeAction{}, res)
		resArray := res.([]CommandOrCodeAction)
		require.IsType(t, CodeAction{}, resArray[0].Get())
	}
}
