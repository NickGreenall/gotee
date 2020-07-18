package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
)

const ClientTestStream string = `{"key":"json"}{"key":"atom", "data": {"A": "1"}}`

func TestBackground(t *testing.T) {
	// Tests run in a redirected environement
	if AmForeground() {
		t.Fatal("I expected to be in the background")
	}
}

func TestSockNotOpen(t *testing.T) {
	if SockOpen("does_not.exist") {
		t.Fatal("Non existent socket exists...")
	}
}

func TestSockOpen(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "temp")
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}

	if !SockOpen(tmpFile.Name()) {
		t.Error("Expected file to exist")
	}

	err = tmpFile.Close()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	os.Remove(tmpFile.Name())
}

func MockConn() io.Reader {
	return bytes.NewBufferString(`{"key":"json"}{"key":"atom", "data": {"A": "1"}}`)
}

func TestSink(t *testing.T) {
	conn := bytes.NewBufferString(ClientTestStream)
	outBuf := new(bytes.Buffer)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	Sink(conn, wg, outBuf)
	out := outBuf.String()
	if out != "{\"A\":\"1\"}\n" {
		t.Fatalf("Unexpected output: %s", out)
	}
}

func TestSniff(t *testing.T) {
	ln, err := net.Listen("unix", "./test.sock")
	rdr, wtr := io.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	wg := new(sync.WaitGroup)

	go Sniff(ln, wg, wtr)

	conn, err := net.Dial("unix", "./test.sock")
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte(ClientTestStream))
	if err != nil {
		t.Error(err)
	}

	out := make([]byte, 10)

	_, err = rdr.Read(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "{\"A\":\"1\"}\n" {
		t.Errorf("Unexpected output: %s", out)
	}

	err = conn.Close()
	if err != nil {
		t.Error(err)
	}

	wg.Wait()
	ln.Close()

}
