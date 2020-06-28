package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/common"
)

type KeyEncoder struct {
	Child  common.Encoder
	Packet KeyPacket
}

type Encoder struct {
	Parent *KeyEncoder
	key    string
}

func (enc *Encoder) Encode(v interface{}) error {
	return enc.Parent.Encode(enc.key, v)
}

func (enc *KeyEncoder) NewEncoderForKey(key string) common.Encoder {
	return &Encoder{enc, key}
}

func (enc *KeyEncoder) Encode(key string, v interface{}) error {
	err := enc.Packet.KeyMarshal(key, v)
	if err != nil {
		return err
	}
	return enc.Child.Encode(enc.Packet)
}
