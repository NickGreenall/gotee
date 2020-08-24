package keyEncoding

//KeyPacket defines an abstract interface which allows
//new encoders to be used with KeyEncoder/Decoder
type KeyPacket interface {
	//KeyMarshal marshals value v into the packet
	//under key.
	KeyMarshal(key string, v interface{}) error
	//GetKey gets the key stored in this packet.
	GetKey() string
	//KeyUnmarshal unmarshals the data stored in
	//the packet into interface v.
	KeyUnmarshal(v interface{}) error
}
