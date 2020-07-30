package main

import (
	"bytes"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

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

func TestInitConn(t *testing.T) {
	ln, err := net.Listen("unix", "./test.sock")
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	defer ln.Close()

	conn, err := InitConn("./test.sock")
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	conn.Close()
}

func TestInitProducer(t *testing.T) {
	outBuf := new(bytes.Buffer)
	prod := InitProducer(outBuf)

	mockAtom := make(atomiser.AtomData)
	mockAtom["match"] = "A: 123"
	mockAtom["dig"] = "123"

	err := prod.AtomEnc.Encode(mockAtom)
	if err != nil {
		t.Error(err)
	}

	out := outBuf.String()
	if out != "{\"Key\":\"atom\",\"Data\":{\"dig\":\"123\",\"match\":\"A: 123\"}}\n" {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestSource(t *testing.T) {
	const testStr = "1\n2\n3\n"
	inBuf := bytes.NewBufferString(testStr)
	outBuf := new(bytes.Buffer)
	err := Source(inBuf, outBuf)
	if err != nil {
		t.Error(err)
	}
	out := outBuf.String()
	if out != "123" {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestInitReader(t *testing.T) {
	const testStr = "1\n2\n3\n"
	inBuf := bytes.NewBufferString(testStr)
	outBuf := new(bytes.Buffer)

	outRdr := InitReader(inBuf, os.Interrupt)
	io.Copy(outBuf, outRdr)

	out := outBuf.String()
	if testStr != out {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestSpawnCopyNoInt(t *testing.T) {
	const testStr = "1\n2\n3\n"
	outBuf := new(bytes.Buffer)
	inBuf := bytes.NewBufferString(testStr)
	rdr, wtr := io.Pipe()
	c := make(chan os.Signal, 1)

	SpawnCopy(wtr, inBuf, c)
	io.Copy(outBuf, rdr)

	out := outBuf.String()
	if testStr != out {
		t.Errorf("Unexpected output: %s", out)
	}
}

func TestSpawnCopyInt(t *testing.T) {
	const testStr = "123"
	outBuf := new(bytes.Buffer)
	inRdr, inWtr := io.Pipe()
	rdr, wtr := io.Pipe()
	c := make(chan os.Signal, 1)
	SpawnCopy(wtr, inRdr, c)

	_, err := inWtr.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	io.CopyN(outBuf, rdr, 3)

	c <- os.Interrupt

	_, err = wtr.Write([]byte("test"))
	if err != io.ErrClosedPipe {
		t.Errorf("Unexpected Error: %s", err)
	}

	out := outBuf.String()
	if testStr != out {
		t.Errorf("Unexpected output: %s", out)
	}
}
