package mock

type MockCoder struct {
	index int
	Calls []interface{}
}

type MockUnmarshaler interface {
	MockUnmarshal(v interface{}) error
}

type MockMarshaler interface {
	MockMarshal() (interface{}, error)
}

func NewMockCoder(calls ...interface{}) *MockCoder {
	return &MockCoder{
		0,
		calls,
	}
}
