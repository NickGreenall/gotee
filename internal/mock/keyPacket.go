package mock

//MockPacket is for testing KeyEncoder/Decoder
type MockPacket struct {
	Key  string
	Data interface{}
}

//MockPacketError is a mock error.
type MockPacketError struct {
	msg string
}

//Error returns the mock error message.
func (e *MockPacketError) Error() string {
	return e.msg
}

//MockUnmarshal is a mock version of unmarshal.
func (pack MockPacket) MockUnmarshal(v interface{}) error {
	target, ok := v.(*MockPacket)
	if !ok {
		return nil
	}
	*target = pack
	return nil
}

//MockMarshal is a mock version of marshal.
func (pack *MockPacket) MockMarshal() (interface{}, error) {
	return *pack, nil
}

//KeyMarshal sets data to v.
func (pack *MockPacket) KeyMarshal(key string, v interface{}) error {
	pack.Key = key
	pack.Data = v
	return nil
}

//GetKey returns the stored key in the packet.
func (pack *MockPacket) GetKey() string {
	return pack.Key
}

//KeyUnmarshal returns the value stored in data dependent of the
//type of v.
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
