package lsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestRPCConnection(t *testing.T) {
	reqData := `
{
	"jsonrpc": "2.0",
	"id": 1,
	"method": "textDocument/didOpen",
	"params": {
	}
}`
	notifData := `
{
	"jsonrpc": "2.0",
	"method": "initialized",
	"params": [123]
}`
	testdata := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(notifData), notifData)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData), reqData)
	resp := ""
	conn := &Connection{
		in: bufio.NewReader(strings.NewReader(testdata)),
		handler: func(ctx context.Context, id *IntOrString, method string, params *ArrayOrObject) error {
			resp += fmt.Sprintf("method=%v ID=%v params=%v\n", method, id, params.Value())
			return nil
		},
	}
	conn.inLoop()
	time.Sleep(100 * time.Millisecond)
	require.Len(t, conn.errors, 1)
	require.Equal(t, conn.errors[0], io.EOF)
	require.Equal(t, "method=textDocument/didOpen ID=&[49] params=[123 10 9 125]\nmethod=initialized ID=<nil> params=[[49 50 51]]\n", resp)
}

func TestIntegerOrString(t *testing.T) {
	var t1 IntOrString
	err := json.Unmarshal([]byte(`10`), &t1)
	require.NoError(t, err)
	require.Equal(t, 10, t1.Value())
	require.IsType(t, 10, t1.Value())
	d1, err := json.Marshal(t1)
	require.NoError(t, err)
	require.Equal(t, []byte(`10`), d1)

	err = json.Unmarshal([]byte(`"10"`), &t1)
	require.NoError(t, err)
	require.Equal(t, "10", t1.Value())
	require.IsType(t, "10", t1.Value())
	d1, err = json.Marshal(t1)
	require.NoError(t, err)
	require.Equal(t, []byte(`"10"`), d1)

	err = json.Unmarshal([]byte(`{"id":"10"}`), &t1)
	require.Error(t, err)
}

func TestJSONResponseMessageSupport(t *testing.T) {
	var r1 ResponseMessage
	in1 := []byte(`
	{
		"jsonrpc": "2.0",
		"id": "10",
		"result": null
	}`)
	err := json.Unmarshal(in1, &r1)
	require.NoError(t, err)
	require.IsType(t, Null{}, r1.Result.Value())
	require.Nil(t, r1.Error)
	d1, err := json.Marshal(r1)
	require.NoError(t, err)
	require.JSONEq(t, string(in1), string(d1))
	fmt.Println(string(d1))

	var r2 ResponseMessage
	in2 := []byte(`
	{
		"jsonrpc": "2.0",
		"id": "10",
		"error": {
			"code" : 9999
		}
	}`)
	err = json.Unmarshal(in2, &r2)
	require.NoError(t, err)
	require.Nil(t, r2.Result.Value())
	d2, err := json.Marshal(r2)
	require.NoError(t, err)
	require.JSONEq(t, string(in2), string(d2))
	fmt.Println(string(d2))

	var r3 ResponseMessage
	in3 := []byte(`
	{
		"jsonrpc": "2.0",
		"error": {
			"code" : 9999
		}
	}`)
	err = json.Unmarshal(in3, &r3)
	require.Error(t, err)
}
