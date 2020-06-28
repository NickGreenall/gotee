package consumer

type ConsumerError struct {
	msg string
}

func (e *ConsumerError) Error() string {
	return e.msg
}
