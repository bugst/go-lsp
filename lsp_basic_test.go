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

		beginJson := `{
			"cancellable" : true,
			"kind" : "begin",
			"message" : "mesg",
			"percentage" : 15,
			"title" : "title"
		}`
		beginData, err := json.Marshal(beginParam)
		require.NoError(t, err)
		require.JSONEq(t, beginJson, string(beginData))

		err = json.Unmarshal([]byte(beginJson), &beginParam)
		require.NoError(t, err)
	}
	{
		begin_json := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "begin",
				"title": "hello",
				"cancellable": true
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(begin_json)))
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

		reportJson := `{
			"cancellable" : true,
			"kind" : "report",
			"message" : "mesg",
			"percentage" : 15
		}`
		reportData, err := json.Marshal(reportParam)
		require.NoError(t, err)
		require.JSONEq(t, reportJson, string(reportData))

		err = json.Unmarshal([]byte(reportJson), &reportParam)
		require.NoError(t, err)
	}
	{
		report_json := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "report",
				"message": "hello",
				"cancellable": true
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(report_json)))
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

		endJson := `{
			"kind" : "end",
			"message" : "mesg"
		}`
		endData, err := json.Marshal(endParam)
		require.NoError(t, err)
		require.JSONEq(t, endJson, string(endData))

		err = json.Unmarshal([]byte(endJson), &endParam)
		require.NoError(t, err)
	}
	{
		end_json := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "end",
				"message": "bye"
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(end_json)))
		require.NoError(t, err)
		require.IsType(t, &ProgressParams{}, r)
		pp := r.(*ProgressParams)
		end := pp.TryToDecodeWellKnownValues()
		require.IsType(t, WorkDoneProgressEnd{}, end)
		require.Equal(t, "bye", end.(WorkDoneProgressEnd).Message)
	}
}
