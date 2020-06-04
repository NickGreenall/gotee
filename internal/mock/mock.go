package mock

type MockEncoder struct {
	index int
	Calls []interface{}
}

func NewMockEncoder(calls ...interface{}) *MockEncoder {
	return &MockEncoder{
		0,
		calls,
	}
}

func (e *MockEncoder) Encode(v interface{}) error {
	call, _ := e.Calls[e.index].(error)
	if call == nil {
		e.Calls[e.index] = v
	}
	e.index++
	return call
}

type MockEncodeError struct {
}

func (e *MockEncodeError) Error() string {
	return "Error"
}
