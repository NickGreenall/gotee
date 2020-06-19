package keyEncoding

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"reflect"
	"testing"
)

func TestGoodDecode(t *testing.T) {
	inputVals := []interface{}{
		mock.MockPacket{"A", 1},
		mock.MockPacket{"B", "foo"},
		mock.MockPacket{"C", 3},
	}
	expectedKeys := []string{"A", "B", "C"}
	expectedVals := []interface{}{1, "foo", 3}

	dec := mock.NewMockCoder(inputVals...)
	kDec := &KeyDecoder{dec, new(mock.MockPacket)}

	for i, val := range expectedVals {
		actual := reflect.New(reflect.TypeOf(val))
		key, err := kDec.Pop()
		if err != nil {
			t.Errorf("Unexpected error during pop: %v", err)
		}
		if key != expectedKeys[i] {
			t.Errorf("Expected key: %v, got: %v", expectedKeys[i], key)
		} else {
			t.Logf("Actual key: %v", key)
		}
		err = kDec.Decode(actual.Interface())
		if err != nil {
			t.Errorf("Unexpected error during decode: %v", err)
		}
		actualDecoded := reflect.Indirect(actual).Interface()
		if actualDecoded != val {
			t.Errorf("Expected value: %v, got: %v", val, actualDecoded)
		} else {
			t.Logf("Actual value: %v", actualDecoded)
		}
	}
}

func TestDecodeWithErrors(t *testing.T) {
	inputVals := []interface{}{
		mock.MockPacket{"A", 1},
		mock.MockPacket{"B", "foo"},
		new(ExpectedError),
	}
	expectedKeys := []string{"A", "B", "C"}
	expectedVals := []interface{}{1, 2.3, 3}
	expectedErrs := []error{nil, new(mock.MockPacketError), new(ExpectedError)}

	dec := mock.NewMockCoder(inputVals...)
	kDec := &KeyDecoder{dec, new(mock.MockPacket)}

	for i, val := range expectedVals {
		actual := reflect.New(reflect.TypeOf(val))
		key, err := kDec.Pop()
		if err != nil {
			if reflect.TypeOf(err) != reflect.TypeOf(expectedErrs[i]) {
				t.Fatalf("Unexpected error during pop: %v", err)
			} else {
				t.Logf("Expected error: %v", err)
			}
		} else {
			if key != expectedKeys[i] {
				t.Errorf("Expected key: %v, got: %v", expectedKeys[i], key)
			} else {
				t.Logf("Actual key: %v", key)
			}
			err = kDec.Decode(actual.Interface())
			if err != nil {
				if reflect.TypeOf(err) != reflect.TypeOf(expectedErrs[i]) {
					t.Fatalf("Unexpected error during decode: %v", err)
				} else {
					t.Logf("Expected error: %v", err)
				}
			} else {
				actualDecoded := reflect.Indirect(actual).Interface()
				if actualDecoded != val {
					t.Errorf("Expected value: %v, got: %v", val, actualDecoded)
				} else {
					t.Logf("Actual value: %v", actualDecoded)
				}

			}
		}
	}
}
