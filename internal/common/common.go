package comm

type Encoder interface {
	Encode(interface{}) error
}
