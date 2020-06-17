package mock

import (
	"io"
)

type MockCoder struct {
	index int
	Calls []interface{}
}

type MockVal struct {
	Val string
}

type MockUnmarshaler interface {
	MockUnmarshal(v interface{}) error
}

func (m *MockVal) Unmarshal(v interface{}) error {
	mV, ok := v.(MockVal)
	if !ok {
		return &MockCoderError{"Mock unmarshal error"}
	}
	*m = mV
	return nil
}
func NewMockCoder(calls ...interface{}) *MockCoder {
	return &MockCoder{
		0,
		calls,
	}
}

func (e *MockCoder) Encode(v interface{}) error {
	call, _ := e.Calls[e.index].(error)
	if call == nil {
		e.Calls[e.index] = v
	}
	e.index++
	return call
}

func (e *MockCoder) Decode(v interface{}) error {
	if e.index >= len(e.Calls) {
		return io.EOF
	}
	err, ok := e.Calls[e.index].(error)
	if !ok {
		u, ok := e.Calls[e.index].(MockUnmarshaler)
		if !ok {
			return &MockDecodeError{"Not valid unmarshaler"}
		}
		err = u.MockUnmarshal(v)
	}
	e.index++
	return err
}

type MockDecodeError struct {
	msg string
}

func (e *MockDecodeError) Error() string {
	return e.msg
}

type MockCoderError struct {
	msg string
}

func (e *MockCoderError) Error() string {
	return e.msg
}
