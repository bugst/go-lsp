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
