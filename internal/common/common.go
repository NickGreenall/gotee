package common

type Encoder interface {
	Encode(interface{}) error
}
