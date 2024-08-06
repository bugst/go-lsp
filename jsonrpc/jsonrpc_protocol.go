//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package jsonrpc

import (
	"fmt"

	"go.bug.st/json"
)

type RequestID json.RawMessage

// A RequestMessage to describe a request between the client and the server. Every
// processed request must send a response back to the sender of the request.
type RequestMessage struct {
	// The language server protocol always uses “2.0” as the jsonrpc version.
	JSONRPC string `json:"jsonrpc,required"`

	// The request id.
	ID json.RawMessage `json:"id,required"`

	// The method to be invoked.
	Method string `json:"method,required"`

	// The method's params.
	Params json.RawMessage `json:"params,omitempty"`
}

// A ResponseMessage sent as a result of a request. If a request doesn’t provide a
// result value the receiver of a request still needs to return a response message
// to conform to the JSON RPC specification. The result property of the ResponseMessageSuccess
// should be set to null in this case to signal a successful request.
type ResponseMessage struct {
	// The language server protocol always uses “2.0” as the jsonrpc version.
	JSONRPC string `json:"jsonrpc,required"`

	// The request id.
	ID json.RawMessage `json:"id,required"`

	// The result of a request. This member is REQUIRED on success.
	// This member MUST NOT exist if there was an error invoking the method.
	Result json.RawMessage `json:"result,omitempty"`

	// The error object in case a request fails.
	Error *ResponseError `json:"error,omitempty"`
}

var NullResult json.RawMessage = []byte("null")

// ResponseError is the error object in case a request fails.
type ResponseError struct {
	// A number indicating the error type that occurred.
	Code ErrorCode `json:"code,required"`

	// A string providing a short description of the error.
	Message string `json:"message"`

	// A primitive or structured value that contains additional
	// information about the error. Can be omitted.
	Data json.RawMessage `json:"data,omitempty"`
}

func (r *ResponseError) AsError() error {
	if r.Message == "" {
		return fmt.Errorf("error code: %d", r.Code)
	}
	return fmt.Errorf("%d %s", r.Code, r.Message)
}

type ErrorCode int

const (
	// Defined by JSON RPC

	ErrorCodesParseError     ErrorCode = -32700
	ErrorCodesInvalidRequest ErrorCode = -32600
	ErrorCodesMethodNotFound ErrorCode = -32601
	ErrorCodesInvalidParams  ErrorCode = -32602
	ErrorCodesInternalError  ErrorCode = -32603

	// This is the start range of JSON RPC reserved error codes.
	// It doesn't denote a real error code. No LSP error codes should
	// be defined between the start and end range. For backwards
	// compatibility the `ServerNotInitialized` and the `UnknownErrorCode`
	// are left in the range.
	//
	// @since 3.16.0
	ErrorCodesJsonrpcReservedErrorRangeStart ErrorCode = -32099

	ErrorCodesServerNotInitialized ErrorCode = -32002
	ErrorCodesUnknownErrorCode     ErrorCode = -32001

	// This is the start range of JSON RPC reserved error codes.
	// It doesn't denote a real error code.
	ErrorCodesJsonrpcReservedErrorRangeEnd ErrorCode = -32000

	// This is the start range of LSP reserved error codes.
	// It doesn't denote a real error code.
	//
	// @since 3.16.0
	ErrorCodesLspReservedErrorRangeStart ErrorCode = -32899

	ErrorCodesContentModified  ErrorCode = -32801
	ErrorCodesRequestCancelled ErrorCode = -32800

	// This is the end range of LSP reserved error codes.
	// It doesn't denote a real error code.
	//
	// @since 3.16.0
	ErrorCodesLspReservedErrorRangeEnd ErrorCode = -32800
)

// NotificationMessage A processed notification message must not send
// a response back. They work like events.
type NotificationMessage struct {
	// The language server protocol always uses “2.0” as the jsonrpc version.
	JSONRPC string `json:"jsonrpc,required"`

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
type ProgressToken string
