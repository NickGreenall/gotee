package mock

func (e *MockCoder) Encode(v interface{}) error {
	var err error
	err, _ = e.Calls[e.index].(error)
	if err == nil {
		u, ok := v.(MockMarshaler)
		if ok {
			e.Calls[e.index], err = u.MockMarshal()
		} else {
			err = &MockEncoderError{
				"Could not marshal",
			}
		}
	}
	e.index++
	return err
}

type MockEncoderError struct {
	msg string
}

func (e *MockEncoderError) Error() string {
	return e.msg
}
