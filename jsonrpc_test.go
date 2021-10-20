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
	reqData2 := `
{
	"jsonrpc": "2.0",
	"id": 2,
	"method": "textDocument/didClose",
	"params": {
	}
}`
	reqData3 := `
{
	"jsonrpc": "2.0",
	"id": 3,
	"method": "tocancel",
	"params": {
	}
}`
	cancelReqData3 := `
{
	"jsonrpc": "2.0",
	"method": "$/cancelRequest",
	"params": { "id":3 }
}`
	notifData := `
{
	"jsonrpc": "2.0",
	"method": "initialized",
	"params": [123]
}`
	testdata := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(notifData), notifData)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData), reqData)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData2), reqData2)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData3), reqData3)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(cancelReqData3), cancelReqData3)
	resp := ""
	output := &bytes.Buffer{}
	var wg sync.WaitGroup
	conn := NewConnection(
		bufio.NewReader(strings.NewReader(testdata)),
		output,
		func(ctx context.Context, method string, params *ArrayOrObject, respCallback func(Any, error)) {
			resp += fmt.Sprintf("REQ method=%v params=%v\n", method, params.Value())
			if method == "tocancel" {
				wg.Add(1)
				go func() {
					select {
					case <-ctx.Done():
					case <-time.After(time.Second):
						t.Log("Request has not been canceled!")
						t.Fail()
					}
					respCallback(nil, nil)
					wg.Done()
				}()
				return
			}
			respCallback(nil, nil)
		},
		func(ctx context.Context, method string, params *ArrayOrObject) {
			resp += fmt.Sprintf("NOT method=%v params=%v\n", method, params.Value())
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
	require.Equal(t,
		"NOT method=initialized params=[[49 50 51]]\n"+
			"REQ method=textDocument/didOpen params=[123 10 9 125]\n"+
			"REQ method=textDocument/didClose params=[123 10 9 125]\n"+
			"REQ method=tocancel params=[123 10 9 125]\n"+
			"", resp)

	//require.Equal(t, "", output.String())
	fmt.Println(output.String())
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
