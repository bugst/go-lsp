package lsp

type ShowMessageParams struct {
	// The message type. See {@link MessageType}.
	Type MessageType `json:"type,required"`

	// The actual message.
	Message string `json:"message,required"`
}

type MessageType int

// An error message.
const MessageTypeError MessageType = 1

// A warning message.
const MessageTypeWarning MessageType = 2

// An information message.
const MessageTypeInfo MessageType = 3

// A log message.
const MessageTypeLog MessageType = 4

type ShowMessageRequestParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type,required"`

	// The actual message
	Message string `json:"message,required"`

	// The message action items to present.
	Actions []MessageActionItem `json:"actions,omitempty"`
}

type MessageActionItem struct {
	// A short title like 'Retry', 'Open Log' etc.
	Title string `json:"title,required"`
}

// Params to show a document.
//
// @since 3.16.0
type ShowDocumentParams struct {
	// The document uri to show.
	URI URI `json:"uri,omitempty"`

	// Indicates to show the resource in an external program.
	// To show for example `https://code.visualstudio.com/`
	// in the default WEB browser set `external` to `true`.
	External bool `json:"external,omitempty"`

	// An optional property to indicate whether the editor
	// showing the document should take focus or not.
	// Clients might ignore this property if an external
	// program is started.
	TakeFocus bool `json:"takeFocus,omitempty"`

	// An optional selection range if the document is a text
	// document. Clients might ignore the property if an
	// external program is started or the file is not a text
	// file.
	Selection Range `json:"selection,omitempty"`
}

// The result of an show document request.
//
// @since 3.16.0
type ShowDocumentResult struct {
	// A boolean indicating if the show was successful.
	Success bool `json:"success,required"`
}

type LogMessageParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type,required"`

	// The actual message
	Message string `json:"message,required"`
}
