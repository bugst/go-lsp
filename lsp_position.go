package lsp

import "fmt"

type Range struct {
	// The range's start position.
	Start Position `json:"start,required"`

	// The range's end position.
	End Position `json:"end,required"`
}

var NilRange = Range{}

func (r Range) String() string {
	return fmt.Sprintf("%s-%s", r.Start, r.End)
}

// Overlaps returns true if the Range overlaps with the given Range p
func (r Range) Overlaps(p Range) bool {
	return r.Start.In(p) || r.End.In(p) || p.Start.In(r) || p.End.In(r)
}

type Position struct {
	// Line position in a document (zero-based).
	Line int `json:"line,required"`

	// Character offset on a line in a document (zero-based).
	Character int `json:"character,required"`
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Character)
}

// In returns true if the Position is within the Range r
func (p Position) In(r Range) bool {
	return p.AfterOrEq(r.Start) && p.BeforeOrEq(r.End)
}

// BeforeOrEq returns true if the Position is before or equal the give Position
func (p Position) BeforeOrEq(q Position) bool {
	return p.Line < q.Line || (p.Line == q.Line && p.Character <= q.Character)
}

// AfterOrEq returns true if the Position is after or equal the give Position
func (p Position) AfterOrEq(q Position) bool {
	return p.Line > q.Line || (p.Line == q.Line && p.Character >= q.Character)
}

// Location represents a location inside a resource, such as a line inside a text file.
type Location struct {
	URI DocumentURI `json:"uri,required"`

	Range Range `json:"range,required"`
}

type LocationLink struct {
	// Span of the origin of this link.
	//
	// Used as the underlined span for mouse interaction. Defaults to the word
	// range at the mouse position.
	OriginSelectionRange *Range `json:"originSelectionRange,omitempty"`

	// The target resource identifier of this link.
	TargetUri DocumentURI `json:"targetUri,required"`

	// The full target range of this link. If the target for example is a symbol
	// then target range is the range enclosing this symbol not including
	// leading/trailing whitespace but everything else like comments. This
	// information is typically used to highlight the range in the editor.
	TargetRange Range `json:"targetRange,required"`

	// The range that should be selected and revealed when this link is being
	// followed, e.g the name of a function. Must be contained by the the
	// `targetRange`. See also `DocumentSymbol#range`
	TargetSelectionRange Range `json:"targetSelectionRange,required"`
}

// A document highlight is a range inside a text document which deserves
// special attention. Usually a document highlight is visualized by changing
// the background color of its range.
type DocumentHighlight struct {
	// The range this highlight applies to.
	Range Range `json:"range,required"`

	// The highlight kind, default is DocumentHighlightKind.Text.
	Kind DocumentHighlightKind `json:"kind,omitempty"`
}

// A document highlight kind.
type DocumentHighlightKind int

// A textual occurrence.
const DocumentHighlightKindText DocumentHighlightKind = 1

// Read-access of a symbol, like reading a variable.
const DocumentHighlightKindRead DocumentHighlightKind = 2

// Write-access of a symbol, like writing to a variable.
const DocumentHighlightKindWrite DocumentHighlightKind = 3
