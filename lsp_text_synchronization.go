//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"fmt"
	"strconv"
)

type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument,required"`
}

type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`
}

type DidChangeTextDocumentParams struct {
	// The document that did change. The version number points
	// to the version after all provided content changes have
	// been applied.
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument,required"`

	// The actual content changes. The content changes describe single state
	// changes to the document. So if there are two content changes c1 (at
	// array index 0) and c2 (at array index 1) for a document in state S then
	// c1 moves the document from S to S' and c2 from S' to S''. So c1 is
	// computed on the state S and c2 is computed on the state S'.
	//
	// To mirror the content of a document using change events use the following
	// approach:
	// - start with the same initial content
	// - apply the 'textDocument/didChange' notifications in the order you
	//   receive them.
	// - apply the `TextDocumentContentChangeEvent`s in a single notification
	//   in the order you receive them.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges,required"`
}

type TextDocumentIdentifier struct {
	// The text document's URI.
	URI DocumentURI `json:"uri,required"`
}

func (t TextDocumentIdentifier) String() string {
	return t.URI.String()
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier

	// The version number of this document.
	//
	// The version number of a document will increase after each change,
	// including undo/redo. The number doesn't need to be consecutive.
	Version int `json:"version,required"`
}

func (v VersionedTextDocumentIdentifier) String() string {
	return fmt.Sprintf("%s@%d", v.TextDocumentIdentifier, v.Version)
}

// An event describing a change to a text document. If range and rangeLength are
// omitted the new text is considered to be the full content of the document.
type TextDocumentContentChangeEvent struct {
	// The range of the document that changed.
	Range *Range `json:"range,omitempty"`

	// The optional length of the range that got replaced.
	//
	// @deprecated use range instead.
	RangeLength *int `json:"rangeLength,omitempty"`

	// The new text for the provided range or the new text of the whole document if range and rangeLength are omitted.
	Text string `json:"text,required"`
}

func (change TextDocumentContentChangeEvent) String() string {
	if change.Range == nil {
		return fmt.Sprintf("FULLTEXT -> %s", strconv.Quote(change.Text))
	}

	l := ""
	if change.RangeLength != nil {
		l = fmt.Sprintf(" (len=%d)", *change.RangeLength)
	}
	return fmt.Sprintf("%s%s -> %s", change.Range, l, strconv.Quote(change.Text))
}

type TextDocumentItem struct {
	// The text document's URI.
	URI DocumentURI `json:"uri,required"`

	// The text document's language identifier.
	LanguageID string `json:"languageId,required"`

	// The version number of this document (it will increase after each
	// change, including undo/redo).
	Version int `json:"version,required"`

	// The content of the opened text document.
	Text string `json:"text,required"`
}

func (t TextDocumentItem) String() string {
	return fmt.Sprintf("%s@%d as '%s'", t.URI, t.Version, t.LanguageID)
}

type DidSaveTextDocumentParams struct {
	// The document that was saved.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// Optional the content when saved. Depends on the includeText value
	// when the save notification was requested.
	Text string `json:"text,omitempty"`
}

type RenameParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams

	// The new name of the symbol. If the given name is not valid the
	// request must return a [ResponseError](#ResponseError) with an
	// appropriate message set.
	NewName string `json:"newName,required"`
}

// The parameters send in a will save text document notification.
type WillSaveTextDocumentParams struct {
	// The document that will be saved.
	RextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The 'TextDocumentSaveReason'.
	Reason TextDocumentSaveReason `json:"reason,required"`
}

// Represents reasons why a text document is saved.
type TextDocumentSaveReason int

// Manually triggered, e.g. by the user pressing save, by starting
// debugging, or by an API call.
const TextDocumentSaveReasonManual TextDocumentSaveReason = 1

// Automatic after a delay.
const TextDocumentSaveReasonAfterDelay TextDocumentSaveReason = 2

// When the editor lost focus.
const TextDocumentSaveReasonFocusOut TextDocumentSaveReason = 3
