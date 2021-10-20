package lsp

//go:generate go run go.bug.st/lsp/generator jsonrpc_protocol.go -w

import (
	"errors"
	"fmt"
	"strings"

	"go.bug.st/json"
)

// A Message as defined by JSON-RPC. The language server protocol
// always uses “2.0” as the jsonrpc version.
type Message struct {
	JSONRPC string `json:"jsonrpc,required"`
}

type RequestID json.RawMessage

// A RequestMessage to describe a request between the client and the server. Every
// processed request must send a response back to the sender of the request.
type RequestMessage struct {
	Message

	// The request id.
	ID json.RawMessage `json:"id,required"`

	// The method to be invoked.
	Method string `json:"method,required"`

	// The method's params.
	Params json.RawMessage `json:"params,omitempty"`
}

// A ResponseMessageSuccess sent as a result of a request. If a request doesn’t provide a
// result value the receiver of a request still needs to return a response message
// to conform to the JSON RPC specification. The result property of the ResponseMessageSuccess
// should be set to null in this case to signal a successful request.
type ResponseMessageSuccess struct {
	Message

	// The request id.
	ID json.RawMessage `json:"id,required"`

	// The result of a request. This member is REQUIRED on success.
	// This member MUST NOT exist if there was an error invoking the method.
	Result json.RawMessage `json:"result,required"`
}

type ResponseMessageError struct {
	Message

	// The request id.
	ID json.RawMessage `json:"id,required"`

	// The error object in case a request fails.
	Error ResponseError `json:"error,required"`
}

// ResponseError is the error object in case a request fails.
type ResponseError struct {
	// A number indicating the error type that occurred.
	Code int `json:"code,required"`

	// A string providing a short description of the error.
	Message string `json:"message"`

	// A primitive or structured value that contains additional
	// information about the error. Can be omitted.
	Data json.RawMessage `json:"data,omitempty"`
}

const (
	// Defined by JSON RPC

	ErrorCodesParseError     = -32700
	ErrorCodesInvalidRequest = -32600
	ErrorCodesMethodNotFound = -32601
	ErrorCodesInvalidParams  = -32602
	ErrorCodesInternalError  = -32603

	// This is the start range of JSON RPC reserved error codes.
	// It doesn't denote a real error code. No LSP error codes should
	// be defined between the start and end range. For backwards
	// compatibility the `ServerNotInitialized` and the `UnknownErrorCode`
	// are left in the range.
	//
	// @since 3.16.0
	ErrorCodesJsonrpcReservedErrorRangeStart = -32099

	ErrorCodesServerNotInitialized = -32002
	ErrorCodesUnknownErrorCode     = -32001

	// This is the start range of JSON RPC reserved error codes.
	// It doesn't denote a real error code.
	ErrorCodesJsonrpcReservedErrorRangeEnd = -32000

	// This is the start range of LSP reserved error codes.
	// It doesn't denote a real error code.
	//
	// @since 3.16.0
	ErrorCodesLspReservedErrorRangeStart = -32899

	ErrorCodesContentModified  = -32801
	ErrorCodesRequestCancelled = -32800

	// This is the end range of LSP reserved error codes.
	// It doesn't denote a real error code.
	//
	// @since 3.16.0
	ErrorCodesLspReservedErrorRangeEnd = -32800
)

// NotificationMessage A processed notification message must not send
// a response back. They work like events.
type NotificationMessage struct {
	Message

	// The method to be invoked.
	Method string `json:"method,required"`

	// The notification's params.
	Params json.RawMessage `json:"params,omitempty"`
}

// CancelParams The base protocol offers support for request cancellation. To
// cancel a request, a notification message with the following properties is sent.
// A request that got canceled still needs to return from the server and send a
// response back. It can not be left open / hanging. This is in line with the
// JSON RPC protocol that requires that every request sends a response back.
// In addition it allows for returning partial results on cancel. If the request
// returns an error response on cancellation it is advised to set the error code
// to ErrorCodesRequestCancelled.
type CancelParams struct {
	// ID The request id to cancel.
	ID json.RawMessage `json:"id,required"`
}

// ProgressParams The base protocol offers also support to report progress in a generic fashion.
// This mechanism can be used to report any kind of progress including work done
// progress (usually used to report progress in the user interface using a progress
// bar) and partial result progress to support streaming of results.
// Progress is reported against a token. The token is different than the request
// ID which allows to report progress out of band and also for notification.
type ProgressParams struct {
	// Token The progress token provided by the client or server.
	Token ProgressToken `json:"token,required"`

	// The progress data.
	Value json.RawMessage `json:"value,required"`
}

// ProgressToken is a progress token
type ProgressToken json.RawMessage

// lsp:generate string|Null as StringOrNull

// lsp:generate DocumentURI|Null as DocumentURIOrNull

// Array represent an Array
type Array []json.RawMessage

// Object represents an object
type Object json.RawMessage

// MarshalJSON implements json.Marshaler
func (n Object) MarshalJSON() ([]byte, error) {
	return n, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (n *Object) UnmarshalJSON(data []byte) error {
	if n == nil {
		return errors.New("lsp.Object: UnmarshalJSON on nil pointer")
	}
	if len(data) < 2 || data[0] != '{' {
		return fmt.Errorf("lsp.Object: expected starting object '{' but founr '%c'", data[0])
	}
	if last := data[len(data)-1]; last != '}' {
		return fmt.Errorf("lsp.Object: object not closed (expected '}' but founr '%c')", last)
	}
	*n = append((*n)[0:0], data...)
	return nil
}

func (n *Object) String() string {
	return string(*n)
}

// Null is a "null" value
type Null struct{}

// MarshalJSON implements json.Marshaler
func (Null) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (n *Null) UnmarshalJSON(data []byte) error {
	if strings.TrimSpace(string(data)) == "null" {
		return nil
	}
	return errors.New("expected 'null'")
}

func (*Null) String() string {
	return "null"
}
