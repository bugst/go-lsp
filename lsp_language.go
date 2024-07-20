//
// Copyright 2024 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

import (
	"go.bug.st/json"
)

type CompletionParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams

	// The completion context. This is only available if the client specifies
	// to send this using the client capability
	// `completion.contextSupport === true`
	Context *CompletionContext `json:"context,omitempty"`
}

type TextDocumentPositionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The position inside the text document.
	Position Position `json:"position,required"`
}

func (t TextDocumentPositionParams) String() string {
	return t.TextDocument.String() + ":" + t.Position.String()
}

// Contains additional information about the context in which a completion
// request is triggered.
type CompletionContext struct {
	// How the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind,required"`

	// The trigger character (a single character) that has trigger code
	// complete. Is undefined if
	// `triggerKind !== CompletionTriggerKind.TriggerCharacter`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

// How a completion was triggered
type CompletionTriggerKind int

// Completion was triggered by typing an identifier (24x7 code
// complete), manual invocation (e.g Ctrl+Space) or via API.
const CompletionTriggerKindInvoked CompletionTriggerKind = 1

// Completion was triggered by a trigger character specified by
// the `triggerCharacters` properties of the
// `CompletionRegistrationOptions`.
const CompletionTriggerKindTriggerCharacter CompletionTriggerKind = 2

// Completion was re-triggered as the current completion list is incomplete.
const CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3

// Params for the CodeActionRequest
type CodeActionParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The document in which the command was invoked.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The range for which the command was invoked.
	Range Range `json:"range,required"`

	// Context carrying additional information.
	Context CodeActionContext `json:"context,required"`
}

// Contains additional diagnostic information about the context in which
// a code action is run.
type CodeActionContext struct {
	// An array of diagnostics known on the client side overlapping the range
	// provided to the `textDocument/codeAction` request. They are provided so
	// that the server knows which errors are currently presented to the user
	// for the given range. There is no guarantee that these accurately reflect
	// the error state of the resource. The primary parameter
	// to compute code actions is the provided range.
	Diagnostics []Diagnostic `json:"diagnostics,required"`

	// Requested kind of actions to return.
	//
	// Actions not of this kind are filtered out by the client before being
	// shown. So servers can omit computing them.
	Only []CodeActionKind `json:"only,omitempty"`
}

// Structure to capture a description for an error code.
//
// @since 3.16.0
type CodeDescription struct {
	// An URI to open with more information about the diagnostic error.
	Href URI `json:"href,required"`
}

type HoverParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
}

type DocumentSymbolParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`
}

type DocumentFormattingParams struct {
	*WorkDoneProgressParams

	// The document to format.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The format options.
	Options FormattingOptions `json:"options,required"`
}

type FormattingOptions map[string]interface{}

/*
// Value-object describing what options formatting should use.
type FormattingOptions struct {
	// Size of a tab in spaces.
	TabSize int `json:"tabSize,required"`

	// Prefer spaces over tabs.
	InsertSpaces bool `json:"insertSpaces,required"`

	// Trim trailing whitespace on a line.
	//
	// @since 3.15.0
	TrimTrailingWhitespace bool `json:"trimTrailingWhitespace,omitempty"`

	// Insert a newline character at the end of the file if one does not exist.
	//
	// @since 3.15.0
	InsertFinalNewline bool `json:"insertFinalNewline,omitempty"`

	// Trim all newlines after the final newline at the end of the file.
	//
	// @since 3.15.0
	TrimFinalNewlines bool `json:"trimFinalNewlines,omitempty"`

	// Signature for further properties.
	// Key map[string]interface{}
	//[key: string]: boolean | integer | string `json:"[key"`
}
*/

type DocumentRangeFormattingParams struct {
	*WorkDoneProgressParams

	// The document to format.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The range to format
	Range Range `json:"range,required"`

	// The format options
	Options FormattingOptions `json:"options,required"`
}

type SignatureHelpParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams

	// The signature help context. This is only available if the client
	// specifies to send this using the client capability
	// `textDocument.signatureHelp.contextSupport === true`
	//
	// @since 3.15.0
	Context *SignatureHelpContext `json:"context,omitempty"`
}

// Additional information about the context in which a signature help request
// was triggered.
//
// @since 3.15.0
type SignatureHelpContext struct {
	// Action that caused signature help to be triggered.
	TriggerKind SignatureHelpTriggerKind `json:"triggerKind,required"`

	// Character that caused signature help to be triggered.
	//
	// This is undefined when triggerKind !==
	// SignatureHelpTriggerKind.TriggerCharacter
	TriggerCharacter string `json:"triggerCharacter,omitempty"`

	// `true` if signature help was already showing when it was triggered.
	//
	// Retriggers occur when the signature help is already active and can be
	// caused by actions such as typing a trigger character, a cursor move, or
	// document content changes.
	IsRetrigger bool `json:"isRetrigger,required"`

	// The currently active `SignatureHelp`.
	//
	// The `activeSignatureHelp` has its `SignatureHelp.activeSignature` field
	// updated based on the user navigating through available signatures.
	ActiveSignatureHelp *SignatureHelp `json:"activeSignatureHelp,omitempty"`
}

// How a signature help was triggered.
//
// @since 3.15.0
type SignatureHelpTriggerKind int

// Signature help was invoked manually by the user or by a command.
const SignatureHelpTriggerKindInvoked SignatureHelpTriggerKind = 1

// Signature help was triggered by a trigger character.
const SignatureHelpTriggerKindTriggerCharacter SignatureHelpTriggerKind = 2

// Signature help was triggered by the cursor moving or by the document
// content changing.
const SignatureHelpTriggerKindContentChange SignatureHelpTriggerKind = 3

// Signature help represents the signature of something
// callable. There can be multiple signature but only one
// active and only one active parameter.
type SignatureHelp struct {
	// One or more signatures. If no signatures are available the signature help
	// request should return `null`.
	Signatures []SignatureInformation `json:"signatures,required"`

	// The active signature. If omitted or the value lies outside the
	// range of `signatures` the value defaults to zero or is ignore if
	// the `SignatureHelp` as no signatures.
	//
	// Whenever possible implementors should make an active decision about
	// the active signature and shouldn't rely on a default value.
	//
	// In future version of the protocol this property might become
	// mandatory to better express this.
	ActiveSignature *int `json:"activeSignature,omitempty"`

	// The active parameter of the active signature. If omitted or the value
	// lies outside the range of `signatures[activeSignature].parameters`
	// defaults to 0 if the active signature has parameters. If
	// the active signature has no parameters it is ignored.
	// In future version of the protocol this property might become
	// mandatory to better express the active parameter if the
	// active signature does have any.
	ActiveParameter *int `json:"activeParameter,omitempty"`
}

// Represents the signature of something callable. A signature
// can have a label, like a function-name, a doc-comment, and
// a set of parameters.
type SignatureInformation struct {
	// The label of this signature. Will be shown in
	// the UI.
	Label string `json:"label,required"`

	// The human-readable doc-comment of this signature. Will be shown
	// in the UI but can be omitted.
	Documentation json.RawMessage `json:"documentation,omitempty"`

	// The parameters of this signature.
	Parameters []ParameterInformation `json:"parameters,omitempty"`

	// The index of the active parameter.
	//
	// If provided, this is used in place of `SignatureHelp.activeParameter`.
	//
	// @since 3.16.0
	ActiveParameter *int `json:"activeParameter,omitempty"`
}

// Represents a parameter of a callable-signature. A parameter can
// have a label and a doc-comment.
type ParameterInformation struct {
	// The label of this parameter information.
	//
	// Either a string or an inclusive start and exclusive end offsets within
	// its containing signature label. (see SignatureInformation.label). The
	// offsets are based on a UTF-16 string representation as `Position` and
	// `Range` does.
	//
	////Note*: a label of type string should be a substring of its containing
	// signature label. Its intended use case is to highlight the parameter
	// label part in the `SignatureInformation.label`.
	Label json.RawMessage `json:"label,required"`

	// The human-readable doc-comment of this parameter. Will be shown
	// in the UI but can be omitted.
	Documentation json.RawMessage `json:"documentation,omitempty"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

type TypeDefinitionParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

type ImplementationParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

type DocumentHighlightParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

type DeclarationParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

type ReferenceParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
	Context *ReferenceContext `json:"context,required"`
}

type ReferenceContext struct {
	// Include the declaration of the current symbol.
	IncludeDeclaration bool `json:"includeDeclaration,required"`
}

type DocumentOnTypeFormattingParams struct {
	TextDocumentPositionParams

	// The character that has been typed.
	Ch string `json:"ch,required"`

	// The format options.
	Options FormattingOptions `json:"options,required"`
}

// Represents a collection of [completion items](#CompletionItem) to be
// presented in the editor.
type CompletionList struct {
	// This list is not complete. Further typing should result in recomputing
	// this list.
	//
	// Recomputed lists have all their items replaced (not appended) in the
	// incomplete completion sessions.
	IsIncomplete bool `json:"isIncomplete,required"`

	// The completion items.
	Items []CompletionItem `json:"items,required"`
}

type CompletionItem struct {
	// The label of this completion item.
	//
	// The label property is also by default the text that
	// is inserted when selecting this completion.
	//
	// If label details are provided the label itself should
	// be an unqualified name of the completion item.
	Label string `json:"label,required"`

	// Additional details for the label
	//
	// @since 3.17.0 - proposed state
	LabelDetails *CompletionItemLabelDetails `json:"labelDetails,omitempty"`

	// The kind of this completion item. Based of the kind
	// an icon is chosen by the editor. The standardized set
	// of available values is defined in `CompletionItemKind`.
	Kind CompletionItemKind `json:"kind,omitempty"`

	// Tags for this completion item.
	//
	// @since 3.15.0
	Tags []CompletionItemTag `json:"tags,omitempty"`

	// A human-readable string with additional information
	// about this item, like type or symbol information.
	Detail string `json:"detail,omitempty"`

	// A human-readable string that represents a doc-comment.
	Documentation json.RawMessage `json:"documentation,omitempty"`

	// Indicates if this item is deprecated.
	//
	// @deprecated Use `tags` instead if supported.
	Deprecated bool `json:"deprecated,omitempty"`

	// Select this item when showing.
	//
	// *Note* that only one completion item can be selected and that the
	// tool / client decides which item that is. The rule is that the *first*
	// item of those that match best is selected.
	Preselect bool `json:"preselect,omitempty"`

	// A string that should be used when comparing this item
	// with other items. When `falsy` the label is used
	// as the sort text for this item.
	SortText string `json:"sortText,omitempty"`

	// A string that should be used when filtering a set of
	// completion items. When `falsy` the label is used as the
	// filter text for this item.
	FilterText string `json:"filterText,omitempty"`

	// A string that should be inserted into a document when selecting
	// this completion. When `falsy` the label is used as the insert text
	// for this item.
	//
	// The `insertText` is subject to interpretation by the client side.
	// Some tools might not take the string literally. For example
	// VS Code when code complete is requested in this example
	// `con<cursor position>` and a completion item with an `insertText` of
	// `console` is provided it will only insert `sole`. Therefore it is
	// recommended to use `textEdit` instead since it avoids additional client
	// side interpretation.
	InsertText string `json:"insertText,omitempty"`

	// The format of the insert text. The format applies to both the
	// `insertText` property and the `newText` property of a provided
	// `textEdit`. If omitted defaults to `InsertTextFormat.PlainText`.
	InsertTextFormat InsertTextFormat `json:"insertTextFormat,omitempty"`

	// How whitespace and indentation is handled during completion
	// item insertion. If not provided the client's default value depends on
	// the `textDocument.completion.insertTextMode` client capability.
	//
	// @since 3.16.0
	// @since 3.17.0 - support for `textDocument.completion.insertTextMode`
	InsertTextMode InsertTextMode `json:"insertTextMode,omitempty"`

	// An edit which is applied to a document when selecting this completion.
	// When an edit is provided the value of `insertText` is ignored.
	//
	// *Note:* The range of the edit must be a single line range and it must
	// contain the position at which completion has been requested.
	//
	// Most editors support two different operations when accepting a completion
	// item. One is to insert a completion text and the other is to replace an
	// existing text with a completion text. Since this can usually not be
	// predetermined by a server it can report both ranges. Clients need to
	// signal support for `InsertReplaceEdit`s via the
	// `textDocument.completion.completionItem.insertReplaceSupport` client
	// capability property.
	//
	// *Note 1:* The text edit's range as well as both ranges from an insert
	// replace edit must be a [single line] and they must contain the position
	// at which completion has been requested.
	// *Note 2:* If an `InsertReplaceEdit` is returned the edit's insert range
	// must be a prefix of the edit's replace range, that means it must be
	// contained and starting at the same position.
	//
	// @since 3.16.0 additional type `InsertReplaceEdit`
	// TODO: TextEdit?: TextEdit | InsertReplaceEdit
	TextEdit *TextEdit `json:"textEdit,omitempty"`

	// An optional array of additional text edits that are applied when
	// selecting this completion. Edits must not overlap (including the same
	// insert position) with the main edit nor with themselves.
	//
	// Additional text edits should be used to change text unrelated to the
	// current cursor position (for example adding an import statement at the
	// top of the file if the completion item will insert an unqualified type).
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`

	// An optional set of characters that when pressed while this completion is
	// active will accept it first and then type that character. *Note* that all
	// commit characters should have `length=1` and that superfluous characters
	// will be ignored.
	CommitCharacters []string `json:"commitCharacters,omitempty"`

	// An optional command that is executed *after* inserting this completion.
	// *Note* that additional modifications to the current document should be
	// described with the additionalTextEdits-property.
	Command *Command `json:"command,omitempty"`

	// A data entry field that is preserved on a completion item between
	// a completion and a completion resolve request.
	Data json.RawMessage `json:"data,omitempty"`
}

// Additional details for a completion item label.
//
// @since 3.17.0 - proposed state
type CompletionItemLabelDetails struct {
	// An optional string which is rendered less prominently directly after
	// {@link CompletionItemLabel.label label}, without any spacing. Should be
	// used for function signatures or type annotations.
	Detail string `json:"detail,omitempty"`

	// An optional string which is rendered less prominently after
	// {@link CompletionItemLabel.detail}. Should be used for fully qualified
	// names or file path.
	Description string `json:"description,omitempty"`
}

// Defines whether the insert text in a completion item should be interpreted as
// plain text or a snippet.
type InsertTextFormat int

// The primary text to be inserted is treated as a plain string.
const InsertTextFormatPlainText InsertTextFormat = 1

// The primary text to be inserted is treated as a snippet.
//
// A snippet can define tab stops and placeholders with `$1`, `$2`
// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
// the end of the snippet. Placeholders with equal identifiers are linked,
// that is typing in one will update others too.
const InsertTextFormatSnippet InsertTextFormat = 2

type TextEdit struct {
	// The range of the text document to be manipulated. To insert
	// text into a document create a range where start === end.
	Range Range `json:"range,required"`

	// The string to be inserted. For delete operations use an
	// empty string.
	NewText string `json:"newText,required"`
}

// A special text edit to provide an insert and a replace operation.
//
// @since 3.16.0
type InsertReplaceEdit struct {
	// The string to be inserted.
	NewText string `json:"newText,required"`

	// The range if the insert is requested
	Insert Range `json:"insert,required"`

	// The range if the replace is requested.
	Replace Range `json:"replace,required"`
}

type Command struct {
	// Title of the command, like `save`.
	Title string `json:"title,required"`

	// The identifier of the actual command handler.
	Command string `json:"command,required"`

	// Arguments that the command handler should be
	// invoked with.
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}

// A code action represents a change that can be performed in code, e.g. to fix
// a problem or to refactor code.
//
// A CodeAction must set either `edit` and/or a `command`. If both are supplied,
// the `edit` is applied first, then the `command` is executed.
type CodeAction struct {

	// A short, human-readable, title for this code action.
	Title string `json:"title,required"`

	// The kind of the code action.
	//
	// Used to filter code actions.
	Kind CodeActionKind `json:"kind,omitempty"`

	// The diagnostics that this code action resolves.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`

	// Marks this as a preferred action. Preferred actions are used by the
	// `auto fix` command and can be targeted by keybindings.
	//
	// A quick fix should be marked preferred if it properly addresses the
	// underlying error. A refactoring should be marked preferred if it is the
	// most reasonable choice of actions to take.
	//
	// @since 3.15.0
	IsPreferred bool `json:"isPreferred,omitempty"`

	// Marks that the code action cannot currently be applied.
	//
	// Clients should follow the following guidelines regarding disabled code
	// actions:
	//
	// - Disabled code actions are not shown in automatic lightbulbs code
	//   action menus.
	//
	// - Disabled actions are shown as faded out in the code action menu when
	//   the user request a more specific type of code action, such as
	//   refactorings.
	//
	// - If the user has a keybinding that auto applies a code action and only
	//   a disabled code actions are returned, the client should show the user
	//   an error message with `reason` in the editor.
	//
	// @since 3.16.0
	Disabled *struct {

		// Human readable description of why the code action is currently
		// disabled.
		//
		// This is displayed in the code actions UI.
		Reason string `json:"reason,required"`
	} `json:"disabled,omitempty"`

	// The workspace edit this code action performs.
	Edit *WorkspaceEdit `json:"edit,omitempty"`

	// A command this code action executes. If a code action
	// provides an edit and a command, first the edit is
	// executed and then the command.
	Command *Command `json:"command,omitempty"`

	// A data entry field that is preserved on a code action between
	// a `textDocument/codeAction` and a `codeAction/resolve` request.
	//
	// @since 3.16.0
	Data json.RawMessage `json:"data,omitempty"`
}

type WorkspaceEdit struct {
	// Holds changes to existing resources.

	// changes?: { [uri: DocumentUri]: TextEdit[]; };
	Changes map[DocumentURI][]TextEdit `json:"changes,omitempty"`

	// Depending on the client capability
	// `workspace.workspaceEdit.resourceOperations` document changes are either
	// an array of `TextDocumentEdit`s to express changes to n different text
	// documents where each text document edit addresses a specific version of
	// a text document. Or it can contain above `TextDocumentEdit`s mixed with
	// create, rename and delete file / folder operations.
	//
	// Whether a client supports versioned document edits is expressed via
	// `workspace.workspaceEdit.documentChanges` client capability.
	//
	// If a client neither supports `documentChanges` nor
	// `workspace.workspaceEdit.resourceOperations` then only plain `TextEdit`s
	// using the `changes` property are supported.

	// documentChanges?: (
	// 	TextDocumentEdit[] | (TextDocumentEdit | CreateFile | RenameFile | DeleteFile)[]
	// );

	// A map of change annotations that can be referenced in
	// `AnnotatedTextEdit`s or create, rename and delete file / folder
	// operations.
	//
	// Whether clients honor this property depends on the client capability
	// `workspace.changeAnnotationSupport`.
	//
	// @since 3.16.0
	//
	// changeAnnotations?: {
	// 	[id: string /* ChangeAnnotationIdentifier */]: ChangeAnnotation;
	// };
	ChangeAnnotations map[string]ChangeAnnotation `json:"changeAnnotations,omitempty"`
}

// Additional information that describes document changes.
//
// @since 3.16.0
type ChangeAnnotation struct {
	// A human-readable string describing the actual change. The string
	// is rendered prominent in the user interface.
	Label string `json:"label,required"`

	// A flag which indicates that user confirmation is needed
	// before applying the change.
	NeedsConfirmation bool `json:"needsConfirmation,omitempty"`

	// A human-readable string which is rendered less prominent in
	// the user interface.
	Description string `json:"description,omitempty"`
}

// The result of a hover request.
type Hover struct {
	// The hover's content
	Contents MarkupContent `json:"contents,required"`

	// An optional range is a range inside a text document
	// that is used to visualize a hover, e.g. by changing the background color.
	Range *Range `json:"range,omitempty"`
}

// MarkedString can be used to render human readable text. It is either a
// markdown string or a code-block that provides a language and a code snippet.
// The language identifier is semantically equal to the optional language
// identifier in fenced code blocks in GitHub issues.
//
// The pair of a language and a value is an equivalent to markdown:
// ```${language}
// ${value}
// ```
//
// Note that markdown strings will be sanitized - that means html will be
// escaped.
//
// @deprecated use MarkupContent instead.
type MarkedString struct {
	Language string `json:"language,required"`
	Value    string `json:"value,required"`
}

// type MarkedString = string | { language: string; value: string };
func (ms *MarkedString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*ms = MarkedString{Value: s}
		return nil
	}

	type __ MarkedString // avoid loops
	var res __
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*ms = MarkedString(res)
	return nil
}

func (ms MarkedString) MarshalJSON() ([]byte, error) {
	if ms.Language == "" {
		return json.Marshal(ms.Value)
	}
	type __ MarkedString // avoid loops
	return json.Marshal(__(ms))
}

// A `MarkupContent` literal represents a string value which content is
// interpreted base on its kind flag. Currently the protocol supports
// `plaintext` and `markdown` as markup kinds.
//
// If the kind is `markdown` then the value can contain fenced code blocks like
// in GitHub issues.
//
// Here is an example how such a string can be constructed using
// JavaScript / TypeScript:
// ```typescript
// let markdown: MarkdownContent = {
// 	kind: MarkupKind.Markdown,
// 	value: [
// 		'# Header',
// 		'Some text',
// 		'```typescript',
// 		'someCode();',
// 		'```'
// 	].join('\n')
// };
// ```
//
// Please Note* that clients might sanitize the return markdown. A client could
// decide to remove HTML from the markdown to avoid script execution.
type MarkupContent struct {
	// The type of the Markup
	Kind MarkupKind `json:"kind,required"`

	// The content itself
	Value string `json:"value,required"`
}

// Represents programming constructs like variables, classes, interfaces etc.
// that appear in a document. Document symbols can be hierarchical and they
// have two ranges: one that encloses its definition and one that points to its
// most interesting range, e.g. the range of an identifier.
type DocumentSymbol struct {

	// The name of this symbol. Will be displayed in the user interface and
	// therefore must not be an empty string or a string only consisting of
	// white spaces.
	Name string `json:"name,required"`

	// More detail for this symbol, e.g the signature of a function.
	Detail string `json:"detail,omitempty"`

	// The kind of this symbol.
	Kind SymbolKind `json:"kind,required"`

	// Tags for this document symbol.
	//
	// @since 3.16.0
	Tags []SymbolTag `json:"tags,omitempty"`

	// Indicates if this symbol is deprecated.
	//
	// @deprecated Use tags instead
	Deprecated bool `json:"deprecated,omitempty"`

	// The range enclosing this symbol not including leading/trailing whitespace
	// but everything else like comments. This information is typically used to
	// determine if the clients cursor is inside the symbol to reveal in the
	// symbol in the UI.
	Range Range `json:"range,required"`

	// The range that should be selected and revealed when this symbol is being
	// picked, e.g. the name of a function. Must be contained by the `range`.
	SelectionRange Range `json:"selectionRange,required"`

	// Children of this symbol, e.g. properties of a class.
	Children []DocumentSymbol `json:"children,omitempty"`
}

// Represents information about programming constructs like variables, classes,
// interfaces etc.
type SymbolInformation struct {
	// The name of this symbol.
	Name string `json:"name,required"`

	// The kind of this symbol.
	Kind SymbolKind `json:"kind,required"`

	// Tags for this symbol.
	//
	// @since 3.16.0
	Tags []SymbolTag `json:"tags,omitempty"`

	// Indicates if this symbol is deprecated.
	//
	// @deprecated Use tags instead
	Deprecated bool `json:"deprecated,omitempty"`

	// The location of this symbol. The location's range is used by a tool
	// to reveal the location in the editor. If the symbol is selected in the
	// tool the range's start information is used to position the cursor. So
	// the range usually spans more then the actual symbol's name and does
	// normally include things like visibility modifiers.
	//
	// The range doesn't have to denote a node range in the sense of a abstract
	// syntax tree. It can therefore not be used to re-construct a hierarchy of
	// the symbols.
	Location Location `json:"location,required"`

	// The name of the symbol containing this symbol. This information is for
	// user interface purposes (e.g. to render a qualifier in the user interface
	// if necessary). It can't be used to re-infer a hierarchy for the document
	// symbols.
	ContainerName string `json:"containerName,omitempty"`
}

// The parameters of a Workspace Symbol Request.
type WorkspaceSymbolParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// A query string to filter symbols by. Clients may send an empty
	// string here to request all symbols.
	Query string `json:"query,required"`
}

// The parameters sent in notifications/requests for user-initiated creation
// of files.
//
// @since 3.16.0
type CreateFilesParams struct {
	// An array of all files/folders created in this operation.
	Files []FileCreate `json:"files,required"`
}

// Represents information on a file/folder create.
//
// @since 3.16.0
type FileCreate struct {
	// A file:// URI for the location of the file/folder being created.
	URI string `json:"uri,required"`
}

// The parameters sent in notifications/requests for user-initiated renames
// of files.
//
// @since 3.16.0
type RenameFilesParams struct {
	// An array of all files/folders renamed in this operation. When a folder
	// is renamed, only the folder will be included, and not its children.
	Files []FileRename `json:"files,required"`
}

// Represents information on a file/folder rename.
//
// @since 3.16.0
type FileRename struct {
	// A file:// URI for the original location of the file/folder being renamed.
	OldUri string `json:"oldUri,required"`

	// A file:// URI for the new location of the file/folder being renamed.
	NewUri string `json:"newUri,required"`
}

// The parameters sent in notifications/requests for user-initiated deletes
// of files.
//
// @since 3.16.0
type DeleteFilesParams struct {
	// An array of all files/folders deleted in this operation.
	Files []FileDelete `json:"files,required"`
}

// Represents information on a file/folder delete.
//
// @since 3.16.0
type FileDelete struct {
	// A file:// URI for the location of the file/folder being deleted.
	Uri string `json:"uri,required"`
}

type CodeLensParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The document to request code lens for.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`
}

// A code lens represents a command that should be shown along with
// source text, like the number of references, a way to run tests, etc.
//
// A code lens is _unresolved_ when no command is associated to it. For
// performance reasons the creation of a code lens and resolving should be done
// in two stages.
type CodeLens struct {
	// The range in which this code lens is valid. Should only span a single
	// line.
	Range Range `json:"range,required"`

	// The command this code lens represents.
	Command *Command `json:"command,omitempty"`

	// A data entry field that is preserved on a code lens item between
	// a code lens and a code lens resolve request.
	Data json.RawMessage `json:"data,omitempty"`
}

type PrepareRenameParams struct {
	TextDocumentPositionParams
}

type FoldingRangeParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The text document.
	RextDocument TextDocumentIdentifier `json:"textDocument,required"`
}

// Enum of known range kinds
type FoldingRangeKind string

// Folding range for a comment
const FoldingRangeKindComment FoldingRangeKind = "comment"

// Folding range for a imports or includes
const FoldingRangeKindImports FoldingRangeKind = "imports"

// Folding range for a region (e.g. `#region`)
const FoldingRangeKindRegion FoldingRangeKind = "region"

// Represents a folding range. To be valid, start and end line must be bigger
// than zero and smaller than the number of lines in the document. Clients
// are free to ignore invalid ranges.
type FoldingRange struct {

	// The zero-based start line of the range to fold. The folded area starts
	// after the line's last character. To be valid, the end must be zero or
	// larger and smaller than the number of lines in the document.
	StartLine int `json:"startLine,required"`

	// The zero-based character offset from where the folded range starts. If
	// not defined, defaults to the length of the start line.
	StartCharacter *int `json:"startCharacter,omitempty"`

	// The zero-based end line of the range to fold. The folded area ends with
	// the line's last character. To be valid, the end must be zero or larger
	// and smaller than the number of lines in the document.
	EndLine int `json:"endLine,required"`

	// The zero-based character offset before the folded range ends. If not
	// defined, defaults to the length of the end line.
	EndCharacter *int `json:"endCharacter,omitempty"`

	// Describes the kind of the folding range such as `comment` or `region`.
	// The kind is used to categorize folding ranges and used by commands like
	// 'Fold all comments'. See [FoldingRangeKind](#FoldingRangeKind) for an
	// enumeration of standardized kinds.
	Kind string `json:"kind,omitempty"`
}

type SelectionRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams

	// The text document.
	RextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The positions inside the text document.
	Positions []Position `json:"positions,required"`
}

type SelectionRange struct {
	// The [range](#Range) of this selection range.
	Range Range `json:"range,required"`
	// The parent selection range containing this range. Therefore
	// `parent.range` must contain `this.range`.
	Parent *SelectionRange `json:"parent,omitempty"`
}

type CallHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

type CallHierarchyItem struct {
	// The name of this item.
	Name string `json:"name,required"`

	// The kind of this item.
	Kind SymbolKind `json:"kind,required"`

	// Tags for this item.
	Tags []SymbolTag `json:"tags,omitempty"`

	// More detail for this item, e.g. the signature of a function.
	Detail string `json:"detail,omitempty"`

	// The resource identifier of this item.
	URI DocumentURI `json:"uri,required"`

	// The range enclosing this symbol not including leading/trailing whitespace
	// but everything else, e.g. comments and code.
	Range Range `json:"range,required"`

	// The range that should be selected and revealed when this symbol is being
	// picked, e.g. the name of a function. Must be contained by the
	// [`range`](#CallHierarchyItem.range).
	SelectionRange Range `json:"selectionRange,required"`

	// A data entry field that is preserved between a call hierarchy prepare and
	// incoming calls or outgoing calls requests.
	Data json.RawMessage `json:"data,omitempty"`
}

type CallHierarchyIncomingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams

	Item CallHierarchyItem `json:"item,required"`
}

type CallHierarchyIncomingCall struct {

	// The item that makes the call.
	From CallHierarchyItem `json:"from,required"`

	// The ranges at which the calls appear. This is relative to the caller
	// denoted by [`this.from`](#CallHierarchyIncomingCall.from).
	FromRanges []Range `json:"fromRanges,required"`
}

type CallHierarchyOutgoingCallsParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	Item CallHierarchyItem `json:"item,required"`
}

type CallHierarchyOutgoingCall struct {
	// The item that is called.
	Ro CallHierarchyItem `json:"to,required"`

	// The range at which this item is called. This is the range relative to
	// the caller, e.g the item passed to `callHierarchy/outgoingCalls` request.
	FromRanges []Range `json:"fromRanges,required"`
}

type SemanticTokensParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`
}

type SemanticTokens struct {
	// An optional result id. If provided and clients support delta updating
	// the client will include the result id in the next semantic token request.
	// A server can then instead of computing all semantic tokens again simply
	// send a delta.
	ResultID string `json:"resultId,omitempty"`

	// The actual tokens.
	Data []int `json:"data,required"`
}

type SemanticTokensDeltaParams struct {
	*WorkDoneProgressParams
	*PartialResultParams

	// The text document.
	RextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The result id of a previous response. The result Id can either point to
	// a full response or a delta response depending on what was received last.
	PreviousResultID string `json:"previousResultId,required"`
}

type SemanticTokensDelta struct {
	ResultID string `json:"resultId,omitempty"`

	// The semantic token edits to transform a previous result into a new
	// result.
	Edits []SemanticTokensEdit `json:"edits,required"`
}

type SemanticTokensEdit struct {
	// The start offset of the edit.
	Start int `json:"start,required"`

	// The count of elements to remove.
	DeleteCount int `json:"deleteCount,required"`

	// The elements to insert.
	Data []int `json:"data,omitempty"`
}

type SemanticTokensRangeParams struct {
	*WorkDoneProgressParams
	*PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument,required"`

	// The range the semantic tokens are requested for.
	Range Range `json:"range,required"`
}

type LinkedEditingRangeParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
}

type LinkedEditingRanges struct {
	// A list of ranges that can be renamed together. The ranges must have
	// identical length and contain identical text content. The ranges cannot overlap.
	Ranges []Range `json:"ranges,required"`

	// An optional word pattern (regular expression) that describes valid contents for
	// the given ranges. If no pattern is provided, the client configuration's word
	// pattern will be used.
	WordPattern string `json:"wordPattern,omitempty"`
}

type MonikerParams struct {
	TextDocumentPositionParams
	*WorkDoneProgressParams
	*PartialResultParams
}

// Moniker uniqueness level to define scope of the moniker.
type UniquenessLevel string

// The moniker is only unique inside a document
const UniquenessLevelDocument UniquenessLevel = "document"

// The moniker is unique inside a project for which a dump got created
const UniquenessLevelProject UniquenessLevel = "project"

// The moniker is unique inside the group to which a project belongs
const UniquenessLevelGroup UniquenessLevel = "group"

// The moniker is unique inside the moniker scheme.
const UniquenessLevelScheme UniquenessLevel = "scheme"

// The moniker is globally unique
const UniquenessLevelGlobal UniquenessLevel = "global"

// The moniker kind.
type MonikerKind string

// The moniker represent a symbol that is imported into a project
const MonikerKindImport MonikerKind = "import"

// The moniker represents a symbol that is exported from a project
const MonikerKindExport MonikerKind = "export"

// The moniker represents a symbol that is local to a project (e.g. a local
// variable of a function, a class not visible outside the project, ...)
const MonikerKindLocal MonikerKind = "local"

// Moniker definition to match LSIF 0.5 moniker definition.
type Moniker struct {
	// The scheme of the moniker. For example tsc or .Net
	Scheme string `json:"scheme,required"`

	// The identifier of the moniker. The value is opaque in LSIF however
	// schema owners are allowed to define the structure if they want.
	Identifier string `json:"identifier,required"`

	// The scope in which the moniker is unique
	Unique UniquenessLevel `json:"unique,required"`

	// The moniker kind if known.
	Kind MonikerKind `json:"kind,omitempty"`
}
