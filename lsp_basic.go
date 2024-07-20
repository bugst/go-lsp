//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"

	"go.bug.st/json"
	"go.bug.st/lsp/jsonrpc"
)

type WorkDoneProgressOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken jsonrpc.ProgressToken `json:"workDoneToken,omitempty"`
}

type PartialResultParams struct {
	// An optional token that a server can use to report partial results (e.g.
	// streaming) to the client.
	PartialResultToken jsonrpc.ProgressToken `json:"partialResultToken,omitempty"`
}

type WorkDoneProgressCreateParams struct {
	// The token to be used to report progress.
	Token json.RawMessage `json:"token,required"`
}

type WorkDoneProgressCancelParams struct {
	// The token to be used to report progress.
	Token json.RawMessage `json:"token,required"`
}

type ProgressParams struct {
	// The progress token provided by the client or server.
	Token json.RawMessage `json:"token,required"`

	// The progress data.
	Value json.RawMessage `json:"value,required"`
}

func (p *ProgressParams) TryToDecodeWellKnownValues() interface{} {
	var begin WorkDoneProgressBegin
	if err := json.Unmarshal(p.Value, &begin); err == nil {
		return begin
	}
	var report WorkDoneProgressReport
	if err := json.Unmarshal(p.Value, &report); err == nil {
		return report
	}
	var end WorkDoneProgressEnd
	if err := json.Unmarshal(p.Value, &end); err == nil {
		return end
	}
	return nil
}

type WorkDoneProgressBegin struct {
	// Kind string `json:"kind,required"` /* automatically set to 'begin' */

	// Mandatory title of the progress operation. Used to briefly inform about
	// the kind of operation being performed.
	//
	// Examples: "Indexing" or "Linking dependencies".
	Title string `json:"title,required"`

	// Controls if a cancel button should show to allow the user to cancel the
	// long running operation. Clients that don't support cancellation are
	// allowed to ignore the setting.
	Cancellable bool `json:"cancellable,omitempty"`

	// Optional, more detailed associated progress message. Contains
	// complementary information to the `title`.
	//
	// Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	// If unset, the previous progress message (if any) is still valid.
	Message string `json:"message,omitempty"`

	// Optional progress percentage to display (value 100 is considered 100%).
	// If not provided infinite progress is assumed and clients are allowed
	// to ignore the `percentage` value in subsequent in report notifications.
	//
	// The value should be steadily rising. Clients are free to ignore values
	// that are not following this rule. The value range is [0, 100]
	Percentage *float64 `json:"percentage,omitempty"`
}

func (p *WorkDoneProgressBegin) UnmarshalJSON(data []byte) error {
	var temp struct {
		Kind string `json:"kind,required"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Kind != "begin" {
		return fmt.Errorf("invalid Kind field value '%s': must be 'begin'", temp.Kind)
	}
	type __ WorkDoneProgressBegin
	var res __
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*p = WorkDoneProgressBegin(res)
	return nil
}

func (p WorkDoneProgressBegin) MarshalJSON() ([]byte, error) {
	var temp struct {
		Kind        string   `json:"kind,required"`
		Title       string   `json:"title,required"`
		Cancellable bool     `json:"cancellable,omitempty"`
		Message     string   `json:"message,omitempty"`
		Percentage  *float64 `json:"percentage,omitempty"`
	}
	temp.Kind = "begin"
	temp.Title = p.Title
	temp.Cancellable = p.Cancellable
	temp.Message = p.Message
	temp.Percentage = p.Percentage
	return json.Marshal(temp)
}

func (p WorkDoneProgressBegin) String() string {
	res := "BEGIN"
	if p.Cancellable {
		res += " (cancellable)"
	}
	res += " " + p.Title
	if p.Message != "" {
		res += ", " + p.Message
	}
	if p.Percentage != nil {
		res += fmt.Sprintf(", %1.1f%%", *p.Percentage)
	}
	return res
}

type WorkDoneProgressReport struct {
	// Kind string `json:"kind,required"` /* automatically set to 'report' */

	// Controls enablement state of a cancel button. This property is only valid
	// if a cancel button got requested in the `WorkDoneProgressBegin` payload.
	//
	// Clients that don't support cancellation or don't support control the
	// button's enablement state are allowed to ignore the setting.
	Cancellable bool `json:"cancellable,omitempty"`

	// Optional, more detailed associated progress message. Contains
	// complementary information to the `title`.
	//
	// Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	// If unset, the previous progress message (if any) is still valid.
	Message string `json:"message,omitempty"`

	// Optional progress percentage to display (value 100 is considered 100%).
	// If not provided infinite progress is assumed and clients are allowed
	// to ignore the `percentage` value in subsequent in report notifications.
	//
	// The value should be steadily rising. Clients are free to ignore values
	// that are not following this rule. The value range is [0, 100]
	Percentage *float64 `json:"percentage,omitempty"`
}

func (p *WorkDoneProgressReport) UnmarshalJSON(data []byte) error {
	var temp struct {
		Kind string `json:"kind,required"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Kind != "report" {
		return fmt.Errorf("invalid Kind field value '%s': must be 'report'", temp.Kind)
	}
	type __ WorkDoneProgressReport
	var res __
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*p = WorkDoneProgressReport(res)
	return nil
}

func (p WorkDoneProgressReport) MarshalJSON() ([]byte, error) {
	var temp struct {
		Kind        string   `json:"kind,required"`
		Cancellable bool     `json:"cancellable,omitempty"`
		Message     string   `json:"message,omitempty"`
		Percentage  *float64 `json:"percentage,omitempty"`
	}
	temp.Kind = "report"
	temp.Cancellable = p.Cancellable
	temp.Message = p.Message
	temp.Percentage = p.Percentage
	return json.Marshal(temp)
}

func (p WorkDoneProgressReport) String() string {
	res := "REPORT"
	if p.Cancellable {
		res += " (cancellable)"
	}
	if p.Message != "" {
		res += ", " + p.Message
	}
	if p.Percentage != nil {
		res += fmt.Sprintf(", %1.1f%%", *p.Percentage)
	}
	return res
}

type WorkDoneProgressEnd struct {
	// Kind string `json:"kind,required"` /* automatically set to 'end' */

	// Optional, a final message indicating to for example indicate the outcome
	// of the operation.
	Message string `json:"message,omitempty"`
}

func (p *WorkDoneProgressEnd) UnmarshalJSON(data []byte) error {
	var temp struct {
		Kind string `json:"kind,required"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Kind != "end" {
		return fmt.Errorf("invalid Kind field value '%s': must be 'end'", temp.Kind)
	}
	type __ WorkDoneProgressEnd
	var res __
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*p = WorkDoneProgressEnd(res)
	return nil
}

func (p WorkDoneProgressEnd) MarshalJSON() ([]byte, error) {
	var temp struct {
		Kind    string `json:"kind,required"`
		Message string `json:"message,omitempty"`
	}
	temp.Kind = "end"
	temp.Message = p.Message
	return json.Marshal(temp)
}

func (p WorkDoneProgressEnd) String() string {
	res := "END "
	if p.Message != "" {
		res += ", " + p.Message
	}
	return res
}

// General parameters to register for a capability.
type Registration struct {
	// The id used to register the request. The id can be used to deregister
	// the request again.
	ID string `json:"id,required"`

	// The method / capability to register for.
	Method string `json:"method,required"`

	// Options necessary for the registration.
	RegisterOptions json.RawMessage `json:"registerOptions,omitempty"`
}

type RegistrationParams struct {
	Registrations []Registration `json:"registrations,required"`
}

// General parameters to unregister a capability.
type Unregistration struct {
	// The id used to unregister the request or notification. Usually an id
	// provided during the register request.
	ID string `json:"id,required"`

	// The method / capability to unregister for.
	Method string `json:"method,required"`
}

type UnregistrationParams struct {
	// This should correctly be named `unregistrations`. However changing this
	// is a breaking change and needs to wait until we deliver a 4.x version
	// of the specification.
	Unregisterations []Unregistration `json:"unregisterations,required"`
}
