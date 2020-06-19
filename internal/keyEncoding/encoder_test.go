package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"testing"
)

type ExpectedError int

func (e *ExpectedError) Error() string {
	return "Expected error"
}

func TestGoodEncode(t *testing.T) {
	inputKeys := []string{
		"A",
		"B",
		"C",
	}
	inputVals := []interface{}{
		1,
		"foo",
		2,
	}

	enc := mock.NewMockCoder(nil, nil, nil)
	kEnc := &KeyEncoder{enc, new(mock.MockPacket)}

	for i, key := range inputKeys {
		err := kEnc.Encode(key, inputVals[i])
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	for i, call := range enc.Calls {
		pack, ok := call.(mock.MockPacket)
		if !ok {
			t.Errorf("Not a valid packet for call index: %v", i)
		}
		if pack.Key != inputKeys[i] {
			t.Errorf(
				"Expected key: %v, Actual key: %v",
				inputKeys[i],
				pack.Key,
			)
		}
		if pack.Data != inputVals[i] {
			t.Errorf(
				"Expected value: %v, Actual value: %v",
				inputVals[i],
				pack.Data,
			)
		}
		if pack.Data == inputVals[i] && pack.Key == inputKeys[i] {
			t.Logf("Received input %v, %v", pack.Key, pack.Data)
		}
	}
}

func TestEncodingError(t *testing.T) {
	inputKeys := []string{
		"A",
		"B",
	}
	inputVals := []interface{}{
		1,
		"foo",
	}
	expectedErrs := []interface{}{
		nil,
		new(ExpectedError),
	}
	enc := mock.NewMockCoder(expectedErrs...)
	kEnc := &KeyEncoder{enc, new(mock.MockPacket)}

	for i, key := range inputKeys {
		err := kEnc.Encode(key, inputVals[i])
		expectedErr, _ := expectedErrs[i].(error)
		if err != expectedErr {
			t.Errorf("Unexpected error: %v", err)
		} else {
			t.Logf("Received expected error: %v", err)
		}
	}

	for i, call := range enc.Calls {
		if expectedErrs[i] == nil {
			pack, ok := call.(mock.MockPacket)
			if !ok {
				t.Errorf("Not a valid packet for call index: %v", i)
			}
			if pack.Key != inputKeys[i] {
				t.Errorf(
					"Expected key: %v, Actual key: %v",
					inputKeys[i],
					pack.Key,
				)
			}
			if pack.Data != inputVals[i] {
				t.Errorf(
					"Expected value: %v, Actual value: %v",
					inputVals[i],
					pack.Data,
				)
			}
			if pack.Data == inputVals[i] && pack.Key == inputKeys[i] {
				t.Logf("Received input %v, %v", pack.Key, pack.Data)
			}
		}
	}
}
