package main

import (
	"encoding/json"
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

func InitConsumer(conn io.Reader, out io.Writer) *Consumer {
	dec := json.NewDecoder(conn)
	cons := new(Consumer)
	cons.Dec = keyEncoding.NewJsonKeyDecoder(dec)
	cons.Out = out
	return cons
}

type Server struct {
	done chan bool
	ln   net.Listener
	wg   *sync.WaitGroup
}

func NewServer(network string, address string) (*Server, error) {
	// TODO setup output multiplexer
	ln, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		make(chan bool),
		ln,
		new(sync.WaitGroup),
	}
	return srv, nil
}

func (srv *Server) Sniff(out io.Writer) {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			select {
			case <-srv.done:
			default:
				// TODO Setup error returning rather than logging.
				log.Fatalln(err)
			}
			return
		}
		srv.wg.Add(1)
		go srv.Sink(conn, out)
	}
}

func (srv *Server) Sink(conn io.ReadWriter, out io.Writer) {
	conn.Write([]byte{100})
	cons := InitConsumer(conn, out)
	err := cons.Consume()
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
	srv.wg.Done()
}

func (srv *Server) Close() {
	close(srv.done)
	srv.ln.Close()
	srv.wg.Wait()
}
