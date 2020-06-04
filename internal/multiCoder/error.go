package multiCoder

type MultiEncodeErr struct {
	msg string
}

type MultiDecodeErr struct {
	msg string
}

func (e *MultiEncodeErr) Error() string {
	return e.msg
}

func (e *MultiDecodeErr) Error() string {
	return e.msg
}
