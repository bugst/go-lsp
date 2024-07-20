//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"

	"go.bug.st/json"
)

type PublishDiagnosticsParams struct {
	// The URI for which diagnostic information is reported.
	URI DocumentURI `json:"uri,required"`

	// Optional the version number of the document the diagnostics are published
	// for.
	//
	// @since 3.15.0
	Version int `json:"version,omitempty"`

	// An array of diagnostic information items.
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	// The range at which the message applies.
	Range Range `json:"range,required"`

	// The diagnostic's severity. Can be omitted. If omitted it is up to the
	// client to interpret diagnostics as error, warning, info or hint.
	Severity DiagnosticSeverity `json:"severity,omitempty"`

	// The diagnostic's code, which might appear in the user interface.
	// integer | string
	Code json.RawMessage `json:"code,omitempty"`

	// An optional property to describe the error code.
	//
	// @since 3.16.0
	CodeDescription *CodeDescription `json:"codeDescription,omitempty"`

	// A human-readable string describing the source of this
	// diagnostic, e.g. 'typescript' or 'super lint'.
	Source string `json:"source,omitempty"`

	// The diagnostic's message.
	Message string `json:"message,required"`

	// Additional metadata about the diagnostic.
	//
	// @since 3.15.0
	Tags []DiagnosticTag `json:"tags,omitempty"`

	// An array of related diagnostic information, e.g. when symbol-names within
	// a scope collide all definitions can be marked via this property.
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`

	// A data entry field that is preserved between a
	// `textDocument/publishDiagnostics` notification and
	// `textDocument/codeAction` request.
	//
	// @since 3.16.0
	Data json.RawMessage `json:"data,omitempty"`
}

type DiagnosticSeverity int

// DiagnosticSeverityError Reports an error.
const DiagnosticSeverityError DiagnosticSeverity = 1

// DiagnosticSeverityWarning Reports a warning.
const DiagnosticSeverityWarning DiagnosticSeverity = 2

// DiagnosticSeverityInformation Reports an information.
const DiagnosticSeverityInformation DiagnosticSeverity = 3

// DiagnosticSeverityHint Reports a hint.
const DiagnosticSeverityHint DiagnosticSeverity = 4

func (ds DiagnosticSeverity) String() string {
	switch ds {
	case DiagnosticSeverityError:
		return "Error"
	case DiagnosticSeverityWarning:
		return "Warning"
	case DiagnosticSeverityInformation:
		return "Information"
	case DiagnosticSeverityHint:
		return "Hint"
	default:
		return fmt.Sprintf("Unknown (%d)", int(ds))
	}
}

// DiagnosticTag The diagnostic tags.
//
// @since 3.15.0
type DiagnosticTag int

// DiagnosticTagUnnecessary Unused or unnecessary code.
//
// Clients are allowed to render diagnostics with this tag faded out
// instead of having an error squiggle.
const DiagnosticTagUnnecessary DiagnosticTag = 1

// DiagnosticTagDeprecated Deprecated or obsolete code.
//
// Clients are allowed to rendered diagnostics with this tag strike through.
const DiagnosticTagDeprecated DiagnosticTag = 2

// DiagnosticRelatedInformation Represents a related message and source code location for a diagnostic.
// This should be used to point to code locations that cause or are related to
// a diagnostics, e.g when duplicating a symbol in a scope.
type DiagnosticRelatedInformation struct {
	// The location of this related diagnostic information.
	Location Location `json:"location,required"`

	// The message of this related diagnostic information.
	Message string `json:"message,required"`
}

func (dt DiagnosticTag) String() string {
	switch dt {
	case DiagnosticTagUnnecessary:
		return "Unnecessary"
	case DiagnosticTagDeprecated:
		return "Deprecated"
	default:
		return fmt.Sprintf("Unknown (%d)", int(dt))
	}
}
