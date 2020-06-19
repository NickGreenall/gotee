package keyEncoding

type KeyPacket interface {
	KeyMarshal(key string, v interface{}) error
	GetKey() string
	KeyUnmarshal(v interface{}) error
}
