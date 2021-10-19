package lsp
import (
	"go.bug.st/json"
	"errors"
)

// Any is an intersection type of string, int, bool, Array, Object
type Any []byte

// MarshalJSON implements json.Marshaler interface
func (m Any) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *Any) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: string, int, bool, Array, Object")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *Any) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}


	var r0 string
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	var r1 int
	if err := json.Unmarshal(data, &r1); err == nil {
		return r1
	}

	var r2 bool
	if err := json.Unmarshal(data, &r2); err == nil {
		return r2
	}

	var r3 Array
	if err := json.Unmarshal(data, &r3); err == nil {
		return r3
	}

	var r4 Object
	if err := json.Unmarshal(data, &r4); err == nil {
		return r4
	}


	return nil
}


// ResponseResult is an intersection type of string, int, bool, Object
type ResponseResult []byte

// MarshalJSON implements json.Marshaler interface
func (m ResponseResult) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *ResponseResult) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: string, int, bool, Object")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *ResponseResult) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}


	var r0 string
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	var r1 int
	if err := json.Unmarshal(data, &r1); err == nil {
		return r1
	}

	var r2 bool
	if err := json.Unmarshal(data, &r2); err == nil {
		return r2
	}

	var r3 Object
	if err := json.Unmarshal(data, &r3); err == nil {
		return r3
	}


	return nil
}


// IntOrString is an intersection type of string, int
type IntOrString []byte

// MarshalJSON implements json.Marshaler interface
func (m IntOrString) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *IntOrString) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: string, int")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *IntOrString) Value() interface{} {
	data := []byte(*r)


	var r0 string
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	var r1 int
	if err := json.Unmarshal(data, &r1); err == nil {
		return r1
	}


	return nil
}


// IntOrNull is an intersection type of int
type IntOrNull []byte

// MarshalJSON implements json.Marshaler interface
func (m IntOrNull) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *IntOrNull) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: int")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *IntOrNull) Value() interface{} {
	data := []byte(*r)

	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}


	var r0 int
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}


	return nil
}


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


// ArrayOrObject is an intersection type of Array, Object
type ArrayOrObject []byte

// MarshalJSON implements json.Marshaler interface
func (m ArrayOrObject) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *ArrayOrObject) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: Array, Object")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *ArrayOrObject) Value() interface{} {
	data := []byte(*r)


	var r0 Array
	if err := json.Unmarshal(data, &r0); err == nil {
		return r0
	}

	var r1 Object
	if err := json.Unmarshal(data, &r1); err == nil {
		return r1
	}


	return nil
}

