package comm

type MockEncoder struct {
	val interface{}
	err error
}

func (e *MockEncoder) Encode(v interface{}) error {
	if e.err != nil {
		return e.err
	} else {
		e.val = v
		return nil
	}
}

type MockEncodeError struct {
}

func (e *MockEncodeError) Error() string {
	return "Error"
}
