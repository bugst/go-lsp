//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package jsonrpc

import "go.bug.st/json"

type Logger interface {
	LogOutgoingRequest(id string, method string, params json.RawMessage)
	LogIncomingRequest(id string, method string, params json.RawMessage) FunctionLogger
	LogOutgoingResponse(id string, method string, resp json.RawMessage, respErr *ResponseError)
	LogIncomingResponse(id string, method string, resp json.RawMessage, respErr *ResponseError)
	LogOutgoingNotification(method string, params json.RawMessage)
	LogIncomingNotification(method string, params json.RawMessage) FunctionLogger
	LogIncomingCancelRequest(id string)
	LogOutgoingCancelRequest(id string)
}

type FunctionLogger interface {
	Logf(format string, a ...interface{})
}

type NullLogger struct{}

func (NullLogger) LogOutgoingRequest(id string, method string, params json.RawMessage) {
}

func (NullLogger) LogIncomingRequest(id string, method string, params json.RawMessage) FunctionLogger {
	return &NullFunctionLogger{}
}

func (NullLogger) LogOutgoingResponse(id string, method string, resp json.RawMessage, respErr *ResponseError) {
}

func (NullLogger) LogIncomingResponse(id string, method string, resp json.RawMessage, respErr *ResponseError) {
}

func (NullLogger) LogOutgoingNotification(method string, params json.RawMessage) {
}

func (NullLogger) LogIncomingNotification(method string, params json.RawMessage) FunctionLogger {
	return &NullFunctionLogger{}
}

func (NullLogger) LogIncomingCancelRequest(id string) {}

func (NullLogger) LogOutgoingCancelRequest(id string) {}

type NullFunctionLogger struct{}

func (NullFunctionLogger) Logf(format string, a ...interface{}) {}
