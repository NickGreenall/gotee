package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

// SockOpen checks if a socket exists in the filesystem namespace
func SockOpen(addr string) bool {
	_, err := os.Stat(addr)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// InitConn opens a connection to the GoTee server on a unix socket
// given by addr. It will poll till the socket address is available,
// then waits for an accept byte from the server (this is to ensure
// connection has been established server side).
func InitConn(ctx context.Context, addr string) (net.Conn, error) {

	for !SockOpen(addr) {
		select {
		case <-ctx.Done():
			return nil, &TimeoutError{}
		default:
			time.Sleep(1)
		}
	}

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", addr)
	if err != nil {
		return nil, err
	}

	accept := make([]byte, 1)
	_, err = conn.Read(accept)
	if err != nil {
		return nil, err
	}

	if accept[0] != 100 {
		return nil, &InitError{}
	}

	return conn, nil
}

// InitProducer creates a new producer object initialised with a JSON
// key encoder on out.
func InitProducer(out io.Writer) *Producer {
	enc := json.NewEncoder(out)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	prod := NewProducer(keyEnc)

	return prod
}

// Source scans lines and writes each line to out.
func Source(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		b := scanner.Bytes()
		_, err := out.Write(b)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

// InitReader creates a reader object from in. This will
// have content copied from in, but will be closed when
// signals ("sig") are received.
func InitReader(in io.Reader, sig ...os.Signal) io.ReadCloser {
	rdr, wtr := io.Pipe()
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	SpawnCopy(wtr, in, c)
	return rdr
}

// SpawnCopy copies the content from wtr to in, closing wtr if
// a signal is received on c.
func SpawnCopy(wtr *io.PipeWriter, in io.Reader, c chan os.Signal) {
	go func() {
		defer close(c)
		_, err := io.Copy(wtr, in)
		if err != io.ErrClosedPipe && err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		<-c
		wtr.Close()
	}()
}
