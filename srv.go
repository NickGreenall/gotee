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
)

func AmForeground() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

func SockOpen(addr string) bool {
	_, err := os.Stat(addr)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Sink(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	cons := new(consumer.Consumer)
	cons.Dec = keyEncoding.NewJsonKeyDecoder(dec)
	cons.Out = os.Stdout
	err := cons.Consume()
	if err != io.EOF {
		log.Printf("Unexpected error: %v", err)
	}
}

func SpawnSniffer(network, addr string) error {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	go Sniff(ln)
	return nil
}

func Sniff(ln net.Listener) {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go Sink(conn)
	}
}
