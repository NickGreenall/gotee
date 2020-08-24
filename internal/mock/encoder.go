package mock

//Encode is a mock encode call. Instead of writing to an underlying
//writer, writes to the MockCoders underlying interface array, unless
//this contains an error, in which case error is returned instead.
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

//MockEncoderError is a mock error.
type MockEncoderError struct {
	msg string
}

//Error returns the mock encoder error.
func (e *MockEncoderError) Error() string {
	return e.msg
}
