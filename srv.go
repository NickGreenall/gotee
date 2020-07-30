package main

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/consumer"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func AmForeground() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

func Sink(conn io.Reader, wg *sync.WaitGroup, out io.Writer) {
	dec := json.NewDecoder(conn)
	cons := new(consumer.Consumer)
	cons.Dec = keyEncoding.NewJsonKeyDecoder(dec)
	cons.Out = out
	err := cons.Consume()
	if err != io.EOF {
		log.Printf("Unexpected error: %v", err)
	}
	wg.Done()
}

func Sniff(ln net.Listener, wg *sync.WaitGroup, out io.Writer) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok {
				s := opErr.Unwrap().Error()
				if s == "use of closed network connection" {
					return
				}
			}
			log.Fatalln(err)
		}
		wg.Add(1)
		go Sink(conn, wg, out)
	}
}
