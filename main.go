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
		// Spawn server if in foreground
		srv, err := NewServer("unix", "./test.sock")
		if err != nil {
			log.Fatalln(err)
		}
		defer srv.Close()

		go srv.Sniff(os.Stdout)

		inStrm = os.Stdin
	} else {
		// Otherwise assume connected to pipe and forward on
		inStrm = io.TeeReader(os.Stdin, os.Stdout)
	}

	// Connect to the foreground server
	conn, err := InitConn("./test.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Setup the producer side
	prod := InitProducer(conn)
	prod.SetJson()
	//prod.SetTemplate("\033[32mdig: {{.dig}}\033[0m\n")

	atmsr, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, prod.AtomEnc)
	if err != nil {
		log.Fatalln(err)
	}

	// Used to close the connection on interrupt
	srcRdr := InitReader(inStrm, os.Interrupt)
	defer srcRdr.Close()

	for {
		// Source parses line by line and forwards
		err = Source(srcRdr, atmsr)
		_, ok := err.(*atomiser.AtomiserError)
		if !ok {
			break
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
}
