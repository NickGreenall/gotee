package mock

//MockCoder is a mock encoder/decoder. To use, initialise
//with an array of interfaces. Each successive call Encode
//will store the value in Calls, or return an error value
//stored in call. Decode does the opposite.
type MockCoder struct {
	index int
	//Calls is an array of interfaces which will store
	//the encoded values or is used to set the decode
	//values. If a call is an error, this will return
	//instead.
	Calls []interface{}
}

type MockUnmarshaler interface {
	MockUnmarshal(v interface{}) error
}

type MockMarshaler interface {
	MockMarshal() (interface{}, error)
}

//NewMockCoder returns a new MockCoder initialesed with
//the calls given.
func NewMockCoder(calls ...interface{}) *MockCoder {
	return &MockCoder{
		0,
		calls,
	}
}
