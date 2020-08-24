package main

import (
	"bytes"
	"context"
	"io"
	"net"
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

func MockConn() io.Reader {
	return bytes.NewBufferString(`{"key":"json"}{"key":"atom", "data": {"A": "1"}}`)
}

func TestSink(t *testing.T) {
	outBuf := new(bytes.Buffer)

	srv := new(Server)
	ln, err := net.Listen("unix", "test.sock")
	if err != nil {
		t.Fatal(err)
	}
	srv.serverContext = context.Background()
	srv.ln = ln
	srv.wg = new(sync.WaitGroup)
	srv.wg.Add(1)
	go func() {
		conn, err := srv.ln.Accept()
		if err != nil {
			t.Error(err)
		}
		srv.Sink(conn, outBuf)
	}()

	conn, err := net.Dial("unix", "test.sock")
	if err != nil {
		t.Error(err)
	}
	accept := make([]byte, 1)
	_, err = conn.Read(accept)
	if err != nil {
		t.Error(err)
	}
	if accept[0] != 100 {
		t.Errorf("unexpected accept byte: %v", accept)
	}
	_, err = conn.Write([]byte(ClientTestStream))
	if err != nil {
		t.Error(err)
	}
	conn.Close()

	srv.wg.Wait()
	srv.ln.Close()
	out := outBuf.String()
	if out != "{\"A\":\"1\"}\n" {
		t.Fatalf("Unexpected output: %s", out)
	}
}

func TestSniff(t *testing.T) {
	srv, err := NewServer(context.Background(), "unix", "./test.sock")
	if err != nil {
		t.Fatal(err)
	}
	rdr, wtr := io.Pipe()

	go srv.Sniff(wtr)

	conn, err := net.Dial("unix", "./test.sock")
	if err != nil {
		t.Error(err)
	}
	accept := make([]byte, 1)
	_, err = conn.Read(accept)
	if err != nil {
		t.Error(err)
	}
	if accept[0] != 100 {
		t.Errorf("Unexpected accept value: %v", accept)
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

	srv.Close()
}
