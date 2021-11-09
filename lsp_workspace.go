//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import "encoding/json"

type DidChangeWatchedFilesParams struct {
	// The actual file events.
	Changes []FileEvent `json:"changes,required"`
}

// An event describing a file change.
type FileEvent struct {
	// The file's URI.
	URI DocumentURI `json:"uri,required"`

	// The change type.
	Type int `json:"type,required"`
}

type ExecuteCommandParams struct {
	*WorkDoneProgressParams

	// The identifier of the actual command handler.
	Command string `json:"command,required"`

	// Arguments that the command should be invoked with.
	Arguments []interface{} `json:"arguments,required"`
}

type ApplyWorkspaceEditParams struct {
	// An optional label of the workspace edit. This label is
	// presented in the user interface for example on an undo
	// stack to undo the workspace edit.
	Label string `json:"label,required"`

	// The edits to apply.
	Edit WorkspaceEdit `json:"edit,required"`
}

type ApplyWorkspaceEditResult struct {
	// Indicates whether the edit was applied or not.
	Applied bool `json:"applied,required"`

	// An optional textual description for why the edit was not applied.
	// This may be used by the server for diagnostic logging or to provide
	// a suitable error for a request that triggered the edit.
	FailureReason string `json:"failureReason,omitempty"`

	// Depending on the client's failure handling strategy `failedChange`
	// might contain the index of the change that failed. This property is
	// only available if the client signals a `failureHandling` strategy
	// in its client capabilities.
	FailedChange int `json:"failedChange,omitempty"`
}

type DidChangeWorkspaceFoldersParams struct {
	// The actual workspace folder change event.
	Event WorkspaceFoldersChangeEvent `json:"event,required"`
}

// The workspace folder change event.
type WorkspaceFoldersChangeEvent struct {
	// The array of added workspace folders
	Added []WorkspaceFolder `json:"added,required"`

	// The array of the removed workspace folders
	Eemoved []WorkspaceFolder `json:"removed,required"`
}

type DidChangeConfigurationParams struct {
	// The actual changed settings
	Settings json.RawMessage `json:"settings,required"`
}
