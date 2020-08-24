package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/common"
)

//KeyEncoder multiplexes over a child encoder via key.
//To use create a new KeyEncoder and set Child (e.g. to
//json.Encoder) and Packet (e.g. to json.KeyPacket).
type KeyEncoder struct {
	//Child is the Encoder which actually does the encoding.
	Child common.Encoder
	//Packet is used as intermediate structure before encoding.
	Packet KeyPacket
}

//Encoder can be created for a particular key and supports
//the Encoder interface.
type Encoder struct {
	//Parent is the KeyEncoder this encoder will use to multiplex on.
	Parent *KeyEncoder
	key    string
}

//Encode marshels v into some coding under encoder key.
func (enc *Encoder) Encode(v interface{}) error {
	return enc.Parent.Encode(enc.key, v)
}

//NewEncoderForKey returns a new encoder which will will encode under the
//given key. It is a helper function so that the KeyEncoder Encode function
//Doesn't have to be directly called.
func (enc *KeyEncoder) NewEncoderForKey(key string) common.Encoder {
	return &Encoder{enc, key}
}

//Encode marshals value v under the key given. A corresponding Key
//decoder can then Pop and Decode this value on the other end.
func (enc *KeyEncoder) Encode(key string, v interface{}) error {
	err := enc.Packet.KeyMarshal(key, v)
	if err != nil {
		return err
	}
	return enc.Child.Encode(enc.Packet)
}
