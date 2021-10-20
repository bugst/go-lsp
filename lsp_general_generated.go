package lsp

import (
	"errors"

	"go.bug.st/json"
)

// WorkspaceFolderArrayOrNull is an intersection type of WorkspaceFolderArray
type WorkspaceFolderArrayOrNull []byte

// MarshalJSON implements json.Marshaler interface
func (m WorkspaceFolderArrayOrNull) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *WorkspaceFolderArrayOrNull) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: WorkspaceFolderArray")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *WorkspaceFolderArrayOrNull) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}

	var r0 WorkspaceFolderArray
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	return nil
}
