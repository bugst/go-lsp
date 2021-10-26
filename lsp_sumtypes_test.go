package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestSumTypes(t *testing.T) {
	com_json := "{\"title\":\"command_title\",\"command\":\"command\"}"
	ca_json := "{\"title\":\"codeaction_title\",\"kind\":\"quickfix\",\"isPreferred\":true,\"command\":{\"title\":\"command_title\",\"command\":\"command\"}}"

	{
		var c CommandOrCodeAction
		c.Set(Command{
			Title:   "command_title",
			Command: "command",
		})
		data, err := json.Marshal(c)
		require.NoError(t, err)
		require.Equal(t, com_json, string(data))
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
		require.Equal(t, ca_json, string(data))
	}

	{
		var c CommandOrCodeAction
		err := json.Unmarshal([]byte(com_json), &c)
		require.NoError(t, err)
		res := c.Get()
		require.IsType(t, Command{}, res)
		require.Equal(t, "command", res.(Command).Command)
	}
	{
		var c CommandOrCodeAction
		err := json.Unmarshal([]byte(ca_json), &c)
		require.NoError(t, err)
		res := c.Get()
		require.IsType(t, CodeAction{}, res)
		require.Equal(t, &Command{Title: "command_title", Command: "command"}, res.(CodeAction).Command)
	}

	// Let's try an array of CommandOrCodeActions...
	{
		jsonIn := json.RawMessage("[" + ca_json + "," + com_json + "," + ca_json + "," + ca_json + "," + com_json + "]")
		res, err := DecodeResponseResult("textDocument/codeAction", jsonIn)
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
}
