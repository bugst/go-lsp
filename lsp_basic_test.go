//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestDecodingProgressReports(t *testing.T) {
	{
		perc := 15.0
		beginParam := &WorkDoneProgressBegin{
			Title:       "title",
			Cancellable: true,
			Message:     "mesg",
			Percentage:  &perc,
		}

		beginJSON := `{
			"cancellable" : true,
			"kind" : "begin",
			"message" : "mesg",
			"percentage" : 15,
			"title" : "title"
		}`
		beginData, err := json.Marshal(beginParam)
		require.NoError(t, err)
		require.JSONEq(t, beginJSON, string(beginData))

		err = json.Unmarshal([]byte(beginJSON), &beginParam)
		require.NoError(t, err)
	}
	{
		beginJSON := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "begin",
				"title": "hello",
				"cancellable": true
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(beginJSON)))
		require.NoError(t, err)
		require.IsType(t, &ProgressParams{}, r)
		pp := r.(*ProgressParams)
		begin := pp.TryToDecodeWellKnownValues()
		require.IsType(t, WorkDoneProgressBegin{}, begin)
		require.Equal(t, "hello", begin.(WorkDoneProgressBegin).Title)
	}
	{
		perc := 15.0
		reportParam := &WorkDoneProgressReport{
			Cancellable: true,
			Message:     "mesg",
			Percentage:  &perc,
		}

		reportJSON := `{
			"cancellable" : true,
			"kind" : "report",
			"message" : "mesg",
			"percentage" : 15
		}`
		reportData, err := json.Marshal(reportParam)
		require.NoError(t, err)
		require.JSONEq(t, reportJSON, string(reportData))

		err = json.Unmarshal([]byte(reportJSON), &reportParam)
		require.NoError(t, err)
	}
	{
		reportJSON := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "report",
				"message": "hello",
				"cancellable": true
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(reportJSON)))
		require.NoError(t, err)
		require.IsType(t, &ProgressParams{}, r)
		pp := r.(*ProgressParams)
		report := pp.TryToDecodeWellKnownValues()
		require.IsType(t, WorkDoneProgressReport{}, report)
		require.Equal(t, "hello", report.(WorkDoneProgressReport).Message)
	}
	{
		endParam := &WorkDoneProgressEnd{
			Message: "mesg",
		}

		endJSON := `{
			"kind" : "end",
			"message" : "mesg"
		}`
		endData, err := json.Marshal(endParam)
		require.NoError(t, err)
		require.JSONEq(t, endJSON, string(endData))

		err = json.Unmarshal([]byte(endJSON), &endParam)
		require.NoError(t, err)
	}
	{
		endJSON := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "end",
				"message": "bye"
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(endJSON)))
		require.NoError(t, err)
		require.IsType(t, &ProgressParams{}, r)
		pp := r.(*ProgressParams)
		end := pp.TryToDecodeWellKnownValues()
		require.IsType(t, WorkDoneProgressEnd{}, end)
		require.Equal(t, "bye", end.(WorkDoneProgressEnd).Message)
	}
}
