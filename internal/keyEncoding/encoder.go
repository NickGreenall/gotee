package keyEncoding

type KeyEncoder struct {
	Child  Encoder
	Packet KeyPacket
}

func (enc *KeyEncoder) Encode(key string, v interface{}) error {
	err := enc.Packet.KeyMarshal(key, v)
	if err != nil {
		return err
	}
	return enc.Child.Encode(enc.Packet)
}
