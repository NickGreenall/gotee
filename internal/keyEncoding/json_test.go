package keyEncoding

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestObj struct {
	Q, W int
}

func TestJsonUnmarshal(t *testing.T) {
	inputVals := [][]byte{
		[]byte(`{"Key": "A", "Data": 1}`),
		[]byte(`{"Key": "B", "Data": "foo"}`),
		[]byte(`{"Key": "C", "Data": {"Q": 1, "W":2}}`),
	}
	expectedPacks := []JsonKeyPacket{
		JsonKeyPacket{"A", []byte("1")},
		JsonKeyPacket{"B", []byte(`"foo"`)},
		JsonKeyPacket{"C", []byte(`{"Q": 1, "W":2}`)},
	}
	for i, in := range inputVals {
		packet := JsonKeyPacket{}
		json.Unmarshal(in, &packet)
		if packet.Key != expectedPacks[i].Key || string(packet.Data) != string(expectedPacks[i].Data) {
			t.Errorf("Expected value: %s, got: %s", expectedPacks[i], packet)
		} else {
			t.Logf("Decoded value: %s", packet)
		}
	}
}

func TestKeyUnmarshalErr(t *testing.T) {
	packet := JsonKeyPacket{"B", []byte(`"foo"`)}
	badVal := new(int)
	err := packet.KeyUnmarshal(badVal)
	uerr, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		t.Errorf("Unexpected error: %v", err)
	} else {
		t.Logf("Expected error: %v", uerr)
	}
}

func TestKeyUnmarshal(t *testing.T) {
	inputPacks := []JsonKeyPacket{
		JsonKeyPacket{"A", []byte("1")},
		JsonKeyPacket{"B", []byte(`"foo"`)},
		JsonKeyPacket{"C", []byte(`{"Q": 1, "W":2}`)},
	}
	expectedVals := []interface{}{
		1,
		"foo",
		TestObj{1, 2},
	}
	for i, input := range inputPacks {
		expected := expectedVals[i]
		actualPtr := reflect.New(reflect.TypeOf(expected))
		err := input.KeyUnmarshal(actualPtr.Interface())
		actual := reflect.Indirect(actualPtr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual.Interface() != expected {
			t.Errorf("Expected val: %v, got: %v", expected, actual.Interface())
		}
	}
}
