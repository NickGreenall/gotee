//mock contains mock decoder/encoders etc.
package mock

import (
	"io"
)

//Decoder mocks out decode. Instead of reading from an
//underlying reader, returns interfaces from the stored
//values. If the stored value is an error, returns that
//instead.
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

//MockDecodeError is a mock error.
type MockDecodeError struct {
	msg string
}

//Error returns the mock error message.
func (e *MockDecodeError) Error() string {
	return e.msg
}
