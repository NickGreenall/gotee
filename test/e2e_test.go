package gotee_test

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"io"
	"reflect"
	"sync"
	"testing"
)

func NewConnectedPair() (*keyEncoding.KeyDecoder, *keyEncoding.KeyEncoder) {
	rdr, wtr := io.Pipe()

	jsonEnc := json.NewEncoder(wtr)
	jsonDec := json.NewDecoder(rdr)

	keyEnc := keyEncoding.NewJsonKeyEncoder(jsonEnc)
	keyDec := keyEncoding.NewJsonKeyDecoder(jsonDec)

	return keyDec, keyEnc
}

func HandleErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestEncodeDecode(t *testing.T) {
	inKeys := []string{"A", "A", "B", "A", "B", "B"}
	inData := []interface{}{1, "2", 3, "4", 5.6, "6"}
	dec, enc := NewConnectedPair()
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for i, inKey := range inKeys {
			err := enc.Encode(inKey, inData[i])
			HandleErr(t, err)
			t.Logf("Encoded value: %v using key: %v", inData[i], inKey)
		}
		wg.Done()
	}()

	for i, inKey := range inKeys {
		outKey, err := dec.Pop()
		HandleErr(t, err)
		if inKey != outKey {
			t.Errorf("Expected key: %v, received key: %v", inKey, outKey)
		} else {
			t.Logf("Received key: %v", inKey)
		}
		outPtr := reflect.New(reflect.TypeOf(inData[i]))
		outValue := reflect.Indirect(outPtr)
		err = dec.Decode(outPtr.Interface())
		HandleErr(t, err)
		if outValue.Interface() != inData[i] {
			t.Errorf("Expected data: %v, received data: %v", inData[i], outValue.Interface())
		} else {
			t.Logf("Received value: %v", inData[i])
		}
	}

	wg.Wait()
}

func TestProducerConsumer(t *testing.T) {
	// TODO
}

func TestAtomiseDecode(t *testing.T) {
	// TODO build on producer/consumer
	inText := []string{
		"num: 123, word: foo",
		"num: 673, word: bar",
		"num: 2, word: dog12",
	}
	outVals := []atomiser.AtomData{
		atomiser.AtomData{
			"match": "num: 123, word: foo",
			"num":   "123", "word": "foo",
		},
		atomiser.AtomData{
			"match": "num: 673, word: bar",
			"num":   "673", "word": "bar",
		},
		atomiser.AtomData{
			"match": "num: 2, word: dog12",
			"num":   "2", "word": "dog12",
		},
	}

	regex := `num:\s(?P<num>\d+),\sword:\s(?P<word>\w+)`
	dec, enc := NewConnectedPair()
	atmsr, err := atomiser.NewAtomiser(regex, enc.NewEncoderForKey("regex"))
	HandleErr(t, err)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for _, line := range inText {
			_, err := atmsr.Write([]byte(line))
			HandleErr(t, err)
			t.Logf("Encoded line: %v", line)
		}
		wg.Done()
	}()

	for _, outVal := range outVals {
		atom := make(atomiser.AtomData)
		outKey, err := dec.Pop()
		HandleErr(t, err)
		if outKey != "regex" {
			t.Errorf("Expected key: regex, received key: %v", outKey)
		}
		err = dec.Decode(&atom)
		HandleErr(t, err)
		if !reflect.DeepEqual(atom, outVal) {
			t.Errorf("Expected data: %s, received data: %s", outVal, atom)
		} else {
			t.Logf("Received value: %s", atom)
		}
	}
	wg.Wait()
}

func TestIntSource(t *testing.T) {
	// TODO
}

func TestClientServer(t *testing.T) {
	// TODO
}

func TestAtomiserClientServer(t *testing.T) {
	// TODO
}
