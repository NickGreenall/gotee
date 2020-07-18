package main

import (
	"github.com/NickGreenall/gotee/internal/atomiser"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	var inStrm io.Reader

	if AmForeground() {
		ln, err := net.Listen("unix", "./test.sock")
		if err != nil {
			log.Fatalln(err)
		}

		defer ln.Close()
		wg := new(sync.WaitGroup)
		defer wg.Wait()

		go Sniff(ln, wg, os.Stdout)

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

	srcRdr := IntReader(inStrm, os.Interrupt)
	defer srcRdr.Close()

	Source(srcRdr, atmsr)
	if err != nil {
		log.Fatalln(err)
	}
}
