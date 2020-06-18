package keyEncoding

import (
	"encoding/json"
)

type JsonKeyPacket struct {
	Key  string
	Data json.RawMessage
}

func (pack *JsonKeyPacket) KeyMarshal(key string, v interface{}) (err error) {
	pack.Key = key
	pack.Data, err = json.Marshal(v)
	return err
}

func (pack *JsonKeyPacket) GetKey() string {
	return pack.Key
}

func (pack *JsonKeyPacket) KeyUnmarshal(v interface{}) error {
	return json.Unmarshal(pack.Data, v)
}
