//keyEncoding contains the encoder/decoder multiplexers
//which multiplex by key.
package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/common"
)

//KeyDecoder is a decoder which enables multiplexion by
//key. To use, create a new KeyEncoder and set the Child
//decoder (e.g. json.Decoder) and the Packet to the
//appropriate KeyPacket type (e.g. JsonKeyPacket)
type KeyDecoder struct {
	// Child is the decoder which is used to decode
	// values from.
	Child common.Decoder
	// Packet is an intermediate structure used to
	// decode into. See the KeyPacket interface.
	Packet KeyPacket
}

//Pop performs a child decode and returns the key
//from the decoded packet. As it actually performs
//docoding here, an error can possibly be returned
func (dec *KeyDecoder) Pop() (string, error) {
	err := dec.Child.Decode(dec.Packet)
	if err != nil {
		return "", err
	}
	return dec.Packet.GetKey(), nil
}

//Decode should only be called after a Pop. This will
//try to unmarshal into v, the value stored in the
//key packet.
func (dec *KeyDecoder) Decode(v interface{}) error {
	return dec.Packet.KeyUnmarshal(v)
}
