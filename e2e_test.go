package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/NickGreenall/gotee"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"io"
	"net"
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
			t.Errorf(
				"Expected data: %v, received data: %v",
				inData[i],
				outValue.Interface(),
			)
		} else {
			t.Logf("Received value: %v", inData[i])
		}
	}

	wg.Wait()
}

func TestNoServer(t *testing.T) {
	inText := "num: 123, word: foo\n" +
		"num: 673, word: bar\n" +
		"num: 2, word: dog12\n"

	outText := "word: foo - num: 123\n" +
		"word: bar - num: 673\n" +
		"word: dog12 - num: 2\n"
	regex := `num:\s(?P<num>\d+),\sword:\s(?P<word>\w+)`

	// Plumbing
	rdr, wtr := io.Pipe()
	inBuf := bytes.NewBufferString(inText)
	outBuf := new(bytes.Buffer)

	//Producer consumer setup
	prod := main.InitProducer(wtr)
	cons := main.InitConsumer(rdr, outBuf)

	atmsr, err := atomiser.NewAtomiser(regex, prod.AtomEnc)
	HandleErr(t, err)

	wg := sync.WaitGroup{}

	// Server simulation
	wg.Add(1)
	go func() {
		err := cons.Consume()
		HandleErr(t, err)
		rdr.Close()
		wg.Done()
	}()

	// Client simulation
	prod.SetTemplate("word: {{.word}} - num: {{.num}}\n")
	err = main.Source(inBuf, atmsr)
	HandleErr(t, err)
	wtr.Close()

	wg.Wait()

	// Assert expected output
	out := outBuf.String()
	if out != outText {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestClientServer(t *testing.T) {
	inText := "num: 123, word: foo\n" +
		"num: 673, word: bar\n" +
		"num: 2, word: dog12\n"

	outText := "word: foo - num: 123\n" +
		"word: bar - num: 673\n" +
		"word: dog12 - num: 2\n"
	regex := `num:\s(?P<num>\d+),\sword:\s(?P<word>\w+)`

	// Plumbing
	inBuf := bytes.NewBufferString(inText)
	outBuf := new(bytes.Buffer)

	// Spawn server
	ln, err := net.Listen("unix", "./test.sock")
	HandleErr(t, err)

	wg := new(sync.WaitGroup)
	go main.Sniff(ln, wg, outBuf)

	//Setup client
	conn, err := main.InitConn("./test.sock")
	HandleErr(t, err)
	prod := main.InitProducer(conn)
	atmsr, err := atomiser.NewAtomiser(regex, prod.AtomEnc)
	HandleErr(t, err)

	// Client simulation
	prod.SetTemplate("word: {{.word}} - num: {{.num}}\n")
	err = main.Source(inBuf, atmsr)
	HandleErr(t, err)
	conn.Close()

	// Wait for server to finish cleaning up
	wg.Wait()
	ln.Close()

	// Assert expected output
	out := outBuf.String()
	if out != outText {
		t.Errorf("Unexpected output: %s", out)
	}
}
