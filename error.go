package main

type ConsumerError struct {
	msg string
}

func (e *ConsumerError) Error() string {
	return e.msg
}

type InitError struct{}

func (*InitError) Error() string {
	return "Did not receive accept byte"
}

type TimeoutError struct{}

func (*TimeoutError) Error() string {
	return "Timedout while trying to connect."
}
