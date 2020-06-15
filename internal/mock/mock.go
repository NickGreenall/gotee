package mock

type MockCoder struct {
	index int
	Calls []interface{}
}

type MockVal struct {
	Val string
}

func (m *MockVal) Unmarshal(v interface{}) error {
	mV, ok := v.(MockVal)
	if !ok {
		return &MockCoderError{}
	}
	*m = mV
	return nil
}

type MockUnmarshaler interface {
	MockUnmarshal(v interface{}) error
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
	err, _ := e.Calls[e.index].(error)
	if err == nil {
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
}

func (e *MockCoderError) Error() string {
	return "Error"
}
