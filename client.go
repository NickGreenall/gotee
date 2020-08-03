package main

import (
	"bufio"
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func SockOpen(addr string) bool {
	_, err := os.Stat(addr)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func InitConn(addr string) (net.Conn, error) {

	for !SockOpen(addr) {
		time.Sleep(1)
	}

	conn, err := net.Dial("unix", addr)
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

func InitProducer(out io.Writer) *Producer {
	enc := json.NewEncoder(out)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	prod := NewProducer(keyEnc)

	return prod
}

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

func InitReader(in io.Reader, sig ...os.Signal) io.ReadCloser {
	rdr, wtr := io.Pipe()
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	SpawnCopy(wtr, in, c)
	return rdr
}

func SpawnCopy(wtr *io.PipeWriter, in io.Reader, c chan os.Signal) {
	go func() {
		defer close(c)
		_, err := io.Copy(wtr, in)
		if err != io.ErrClosedPipe && err != nil {
			// TODO better solution
			log.Print(err)
		}
	}()

	go func() {
		<-c
		wtr.Close()
	}()
}
