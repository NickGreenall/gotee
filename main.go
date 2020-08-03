package main

import (
	"github.com/NickGreenall/gotee/internal/atomiser"
	"io"
	"log"
	"os"
)

func main() {
	var inStrm io.Reader

	if AmForeground() {
		srv, err := NewServer("unix", "./test.sock")
		if err != nil {
			log.Fatalln(err)
		}

		defer srv.Close()

		go srv.Sniff(os.Stdout)

		inStrm = os.Stdin
	} else {
		inStrm = io.TeeReader(os.Stdin, os.Stdout)
	}

	conn, err := InitConn("./test.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	prod := InitProducer(conn)
	prod.SetJson()
	//prod.SetTemplate("\033[32mdig: {{.dig}}\033[0m\n")

	atmsr, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, prod.AtomEnc)
	if err != nil {
		log.Fatalln(err)
	}

	srcRdr := InitReader(inStrm, os.Interrupt)
	defer srcRdr.Close()

	Source(srcRdr, atmsr)
	if err != nil {
		log.Fatalln(err)
	}
}
