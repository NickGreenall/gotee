package main

import (
	"bufio"
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"github.com/NickGreenall/gotee/internal/producer"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

func InitConn(addr string) (net.Conn, error) {

	for !SockOpen(addr) {
		time.Sleep(1)
	}

	conn, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func InitProducer(out io.Writer) *producer.Producer {
	enc := json.NewEncoder(out)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	prod := producer.NewProducer(keyEnc)

	return prod
}

func Source(in io.ReadCloser, out io.Writer) error {
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

func IntReader(in io.Reader, sig ...os.Signal) io.ReadCloser {
	rdr, wtr := io.Pipe()
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	SpawnCopy(wtr, in, c)
	return rdr
}

func SpawnCopy(wtr io.WriteCloser, in io.Reader, c chan os.Signal) {
	go func() {
		defer close(c)
		io.Copy(wtr, in)
	}()

	go func() {
		<-c
		wtr.Close()
	}()
}
