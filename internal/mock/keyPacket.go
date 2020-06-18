package mock

type MockPacket struct {
	Key  string
	Data interface{}
}

type MockPacketError struct {
	msg string
}

func (e *MockPacketError) Error() string {
	return e.msg
}

func (pack MockPacket) MockUnmarshal(v interface{}) error {
	target, ok := v.(*MockPacket)
	if !ok {
		return nil
	}
	*target = pack
	return nil
}

func (pack *MockPacket) MockMarshal() (interface{}, error) {
	return *pack, nil
}

func (pack *MockPacket) KeyMarshal(key string, v interface{}) error {
	pack.Key = key
	pack.Data = v
	return nil
}

func (pack *MockPacket) GetKey() string {
	return pack.Key
}

func (pack *MockPacket) KeyUnmarshal(v interface{}) error {
	switch umaskd := pack.Data.(type) {
	case string:
		targ, ok := v.(*string)
		if ok {
			*targ = umaskd
			return nil
		}
	case int:
		targ, ok := v.(*int)
		if ok {
			*targ = umaskd
			return nil
		}
	case MockUnmarshaler:
		return umaskd.MockUnmarshal(v)
	default:
	}
	return &MockPacketError{"Could not unmarshal into value"}
}
