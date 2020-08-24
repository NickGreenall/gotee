package keyEncoding

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/common"
)

//JsonKeyPacket implements the KeyPacket interface for
//encoding/json. Using an object of this type allows
//you to use KeyEncoder/Decoder with json.Encoder/Decoder
type JsonKeyPacket struct {
	//Key stored in packet during marshal/unmarshal
	Key string
	//Data stored in packet. Data is controlled so
	//that key can be unmarshalled/marshalled independently
	//of the data being encoded.
	Data json.RawMessage
}

//KeyMarshal marshals v into the packet under key.
func (pack *JsonKeyPacket) KeyMarshal(key string, v interface{}) (err error) {
	pack.Key = key
	pack.Data, err = json.Marshal(v)
	return err
}

//GetKey gets the key stored in the packet.
func (pack *JsonKeyPacket) GetKey() string {
	return pack.Key
}

//KeyUnmarshal unmarshals data into v.
func (pack *JsonKeyPacket) KeyUnmarshal(v interface{}) error {
	return json.Unmarshal(pack.Data, v)
}

//NewJsonKeyEncoder is a helper wrapper function which constructs
//a KeyEncoder for a given json encoder.
func NewJsonKeyEncoder(enc common.Encoder) *KeyEncoder {
	return &KeyEncoder{enc, new(JsonKeyPacket)}
}

//NewJsonKeyDecoder is a helper wrapper function which constructs
//a KeyDecoder for a given json decoder.
func NewJsonKeyDecoder(dec common.Decoder) *KeyDecoder {
	return &KeyDecoder{dec, new(JsonKeyPacket)}
}
