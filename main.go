package main

import (
	"flag"
	"fmt"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"io"
	"log"
	"os"
	"syscall"
)

const (
	DEFAULT_FORMAT = "JSON"
)

var (
	format = flag.String("f", DEFAULT_FORMAT, "Set output format, can be custom template")
)

func usage() {
	fmt.Printf("%s:\n", os.Args[0])
	fmt.Print(
		"Tees the standard input to a foreground print server.\n" +
			"\n" +
			"Usage:\n" +
			"  gotee [-f string] \"REGEXP\"\n" +
			"\n" +
			"  REGEXP - Can be any valid GO regexp. Named groups can be accessed in format templates.\n" +
			"\n",
	)
	flag.PrintDefaults()
}

func main() {
	var inStrm io.Reader

	// Parse flags and get args
	flag.Usage = usage
	flag.Parse()
	pattern := flag.Arg(0)

	// Setup socket based on process group
	pgrp := syscall.Getpgrp()
	sock := fmt.Sprintf("./.%d.gotee.sock", pgrp)

	if AmForeground() {
		// Spawn server if in foreground
		srv, err := NewServer("unix", sock)
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
	conn, err := InitConn(sock)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Setup the producer side
	prod := InitProducer(conn)
	switch *format {
	case "JSON":
		prod.SetJson()
	default:
		// TODO comeup with more explicit way to handle newlines...
		prod.SetTemplate(*format + "\n")
	}

	atmsr, err := atomiser.NewAtomiser(pattern, prod.AtomEnc)
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
