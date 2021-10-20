package lsp

import (
	"errors"

	"go.bug.st/json"
)

// StringOrNull is an intersection type of string
type StringOrNull []byte

// MarshalJSON implements json.Marshaler interface
func (m StringOrNull) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *StringOrNull) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: string")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *StringOrNull) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}

	var r0 string
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	return nil
}

// DocumentURIOrNull is an intersection type of DocumentURI
type DocumentURIOrNull []byte

// MarshalJSON implements json.Marshaler interface
func (m DocumentURIOrNull) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *DocumentURIOrNull) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: DocumentURI")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *DocumentURIOrNull) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}

	var r0 DocumentURI
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	return nil
}
