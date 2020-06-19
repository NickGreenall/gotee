package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/common"
)

type KeyDecoder struct {
	Child  common.Decoder
	Packet KeyPacket
}

func (dec *KeyDecoder) Pop() (string, error) {
	err := dec.Child.Decode(dec.Packet)
	if err != nil {
		return "", err
	}
	return dec.Packet.GetKey(), nil
}

func (dec *KeyDecoder) Decode(v interface{}) error {
	return dec.Packet.KeyUnmarshal(v)
}
