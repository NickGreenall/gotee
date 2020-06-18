package mock

import (
	"io"
)

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
