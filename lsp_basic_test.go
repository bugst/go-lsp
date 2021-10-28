package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bug.st/json"
)

func TestDecodingProgressReports(t *testing.T) {
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
		report_json := `{
			"token": "aaaa10292992929287172",
			"value": {
				"kind": "end",
				"message": "bye"
			}
		}`
		r, err := DecodeServerNotificationParams("$/progress", json.RawMessage([]byte(report_json)))
		require.NoError(t, err)
		require.IsType(t, &ProgressParams{}, r)
		pp := r.(*ProgressParams)
		end := pp.TryToDecodeWellKnownValues()
		require.IsType(t, WorkDoneProgressEnd{}, end)
		require.Equal(t, "bye", end.(WorkDoneProgressEnd).Message)
	}
}
