package keyEncoding

type KeyPacket interface {
	KeyMarshal(key string, v interface{}) error
	GetKey() string
	KeyUnmarshal(v interface{}) error
}

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}
