// Gotee is an application which can be used to tee content to a server
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DEFAULT_FORMAT = "{{.match}}"
	DEFAULT_REGEXP = ".*"
)

var (
	format    = flag.String("f", DEFAULT_FORMAT, "Set output format, can be custom template")
	pattern   = flag.String("r", DEFAULT_REGEXP, "Sets the regexp to parse Stdin and send to print server")
	runServer = flag.Bool("S", false, "Run server. If not present, will run if foreground process (tail of process group)")
	trunc     = flag.Bool("t", false, "Truncate rather than tee. Stdin will not be forwarded to output and does not run server if foreground.")
	bkGnd     = flag.Bool("b", false, "Background mode, don't try to spindup server if not in foreground")
)

func usage() {
	// TODO look into finding a package to handle the CLI help message.
	fmt.Printf("%s:\n", os.Args[0])
	fmt.Print(
		"Tees the standard input to a foreground print server.\n" +
			"\n" +
			"Usage:\n" +
			"  gotee [-S] [-t] [-b] [-f string] [-r string] [SOCKET]\n" +
			"\n" +
			"  SOCKET - Socket to send parsed standard input.\n" +
			"\n" +
			"  By default a server will be started if the process is in the foreground.\n" +
			"  If no socked is provided, a socket will be autogenrated based on the process group.\n" +
			"\n",
	)
	flag.PrintDefaults()
}

func main() {
	var inStrm io.Reader = os.Stdin
	var srv *Server

	// Parse flags and get args
	flag.Usage = usage
	flag.Parse()
	sock := flag.Arg(0)

	appCtx, appCancel := context.WithCancel(context.Background())

	if sock == "" {
		// Setup socket based on process group
		pgrp := syscall.Getpgrp()
		sock = fmt.Sprintf("./.%d.gotee.sock", pgrp)
	}

	if *runServer || (!*bkGnd && !*trunc && AmForeground()) {
		var err error
		// Spawn server if in foreground
		srv, err = NewServer(appCtx, "unix", sock)
		if err != nil {
			log.Fatalln(err)
		}

		go srv.Sniff(os.Stdout)
	} else if !*trunc {
		// Otherwise tee
		inStrm = io.TeeReader(os.Stdin, os.Stdout)
	}

	// Connect to the foreground server
	clientCtx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	conn, err := InitConn(clientCtx, sock)
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

	atmsr, err := atomiser.NewAtomiser(*pattern, prod.AtomEnc)
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

	if srv != nil {
		log.Println("Closing server, press <C-c> to force kill")
		c := make(chan os.Signal)
		go func() {
			<-c
			appCancel()
		}()
		signal.Notify(c, os.Interrupt)
		srv.Close()
		close(c)
	}
}
