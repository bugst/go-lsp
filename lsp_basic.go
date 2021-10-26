package lsp

import (
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

type ProgressParams struct {
	// The progress token provided by the client or server.
	Token json.RawMessage `json:"token,required"`

	// The progress data.
	Value json.RawMessage `json:"value,required"`
}

func (p *ProgressParams) TryToDecodeWellKnownValues() interface{} {
	var disc struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(p.Value, &disc); err == nil {
		switch disc.Kind {
		case "begin":
			var res WorkDoneProgressBegin
			if err := json.Unmarshal(p.Value, &res); err == nil {
				return res
			}
		case "report":
			var res WorkDoneProgressReport
			if err := json.Unmarshal(p.Value, &res); err == nil {
				return res
			}
		case "end":
			var res WorkDoneProgressEnd
			if err := json.Unmarshal(p.Value, &res); err == nil {
				return res
			}
		}
	}
	return nil
}

type WorkDoneProgressBegin struct {
	Kind string `json:"kind,required"` /* 'begin' */

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
	Percentage int `json:"percentage,omitempty"`
}

type WorkDoneProgressReport struct {
	Kind string `json:"kind,required"` // 'report'

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
	Percentage int `json:"percentage,omitempty"`
}

type WorkDoneProgressEnd struct {
	Kind string `json:"kind,required"` // 'end'

	// Optional, a final message indicating to for example indicate the outcome
	// of the operation.
	Message string `json:"message,omitempty"`
}
