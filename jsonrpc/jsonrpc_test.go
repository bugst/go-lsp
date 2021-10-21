package jsonrpc

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
	simulatedResult := `
{
	"jsonrpc": "2.0",
	"id": "1",
	"result": {"fakedata":999}
}`
	testdata := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(notifData), notifData)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData), reqData)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData2), reqData2)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(reqData3), reqData3)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(cancelReqData3), cancelReqData3)
	testdata += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(simulatedResult), simulatedResult)
	resp := ""
	output := &bytes.Buffer{}
	var wg sync.WaitGroup
	conn := NewConnection(
		bufio.NewReader(strings.NewReader(testdata)),
		output,
		func(ctx context.Context, method string, params json.RawMessage, respCallback func(result json.RawMessage, err *ResponseError)) {
			resp += fmt.Sprintf("REQ method=%v params=%v\n", method, params)
			if method == "tocancel" {
				wg.Add(1)
				go func() {
					select {
					case <-ctx.Done():
					case <-time.After(time.Second):
						t.Log("Request has not been canceled!")
						t.Fail()
					}
					respCallback(nil, &ResponseError{Code: 1, Message: "error message"})
					wg.Done()
				}()
				return
			}
			respCallback(NullResult, nil)
		},
		func(ctx context.Context, method string, params json.RawMessage) {
			resp += fmt.Sprintf("NOT method=%v params=%v\n", method, params)
		},
		func(e error) {
			if e == io.EOF {
				return
			}
			resp += fmt.Sprintf("error=%s\n", e)
		},
	)
	msg, err := json.Marshal(struct{ Field bool }{true})
	require.NoError(t, err)
	cancel, wait, err := conn.SendRequest("helloworld", msg)
	require.NoError(t, err)
	cancel() // send cancel (will be catched later)

	conn.Run()                 // Exits when input is fully consumed
	respRes, respErr := wait() // Wait response to request
	require.Nil(t, respErr)
	require.Equal(t, "{\"fakedata\":999}", string(respRes))

	wg.Wait() // Wait for all pending responses to get through
	conn.Close()
	require.Equal(t,
		"NOT method=initialized params=[91 49 50 51 93]\n"+
			"REQ method=textDocument/didOpen params=[123 10 9 125]\n"+
			"REQ method=textDocument/didClose params=[123 10 9 125]\n"+
			"REQ method=tocancel params=[123 10 9 125]\n"+
			"", resp)

	require.Equal(t, "Content-Length: 72\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":\"1\",\"method\":\"helloworld\",\"params\":{\"Field\":true}}Content-Length: 64\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"$/cancelRequest\",\"params\":{\"id\":\"1\"}}Content-Length: 38\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":null}Content-Length: 38\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":2,\"result\":null}Content-Length: 69\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":3,\"error\":{\"code\":1,\"message\":\"error message\"}}", output.String())
	// fmt.Println(output.String())
}

func TestUnmarshallingRequestVsResponse(t *testing.T) {
	require := require.New(t)
	var notification NotificationMessage
	var request RequestMessage
	var response ResponseMessage

	notificationData := []byte(`{
		"jsonrpc": "2.0",
		"method": "textDocument/didOpen",
		"params": {
		}
	}`)
	require.NoError(json.Unmarshal(notificationData, &notification))
	require.Error(json.Unmarshal(notificationData, &request))
	require.Error(json.Unmarshal(notificationData, &response))

	requestData := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "textDocument/didOpen",
		"params": {
		}
	}`)
	require.NoError(json.Unmarshal(requestData, &notification))
	// a Request can be unmarshalled as a Notification or a Result too !! BEWARE
	require.NoError(json.Unmarshal(requestData, &request))
	require.NoError(json.Unmarshal(requestData, &response))

	responseSuccessData := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
		}
	}`)
	require.Error(json.Unmarshal(responseSuccessData, &notification))
	require.Error(json.Unmarshal(responseSuccessData, &request))
	require.NoError(json.Unmarshal(responseSuccessData, &response))

	responseErrorData := []byte(`{
		"jsonrpc": "2.0",
		"id": 1,
		"error": {
			"code": 1
		}
	}`)
	require.Error(json.Unmarshal(responseErrorData, &notification))
	require.Error(json.Unmarshal(responseErrorData, &request))
	require.NoError(json.Unmarshal(responseErrorData, &response))
}
