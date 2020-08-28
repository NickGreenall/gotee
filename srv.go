package main

import (
	"context"
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"github.com/NickGreenall/gotee/internal/muxWriter"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// AmForeground returns true if the output is a terminal.
func AmForeground() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

// InitConsumer sets up a consumer on given connection for
// given output using JSON decoder.
func InitConsumer(conn io.Reader, out io.Writer) *Consumer {
	dec := json.NewDecoder(conn)
	cons := new(Consumer)
	cons.Dec = keyEncoding.NewJsonKeyDecoder(dec)
	cons.Out = out
	return cons
}

// Server represents a server instance. Can be used to control and
// terminate client connections.
type Server struct {
	serverContext   context.Context
	listenerContext context.Context
	cancelListen    context.CancelFunc
	ln              net.Listener
	wg              *sync.WaitGroup
}

// NewServer creates a new server instance connected over given network
// and address. The will use the given content to control early cancalation.
func NewServer(serverContext context.Context, network string, address string) (*Server, error) {
	ln, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	listenerContext, cancelListen := context.WithCancel(serverContext)
	srv := &Server{
		serverContext,
		listenerContext,
		cancelListen,
		ln,
		new(sync.WaitGroup),
	}
	return srv, nil
}

// Sniff waits for a new connection spawning Sink for a new connection.
// Sink will write received output to out.
func (srv *Server) Sniff(out io.Writer) {
	mux := muxWriter.NewMux(out)
	go func() {
		<-srv.serverContext.Done()
		mux.Close()
	}()
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			select {
			case <-srv.listenerContext.Done():
			default:
				log.Fatalln(err)
			}
			return
		}
		srv.wg.Add(1)
		sinkWrtr := mux.NewWriter()
		go srv.Sink(conn, sinkWrtr)
	}
}

// Sink takes creates a consumer for given connection connected to output.
// All content on the connection will be consumed by the consumer.
func (srv *Server) Sink(conn io.ReadWriteCloser, out io.Writer) {
	connContex, connCancel := context.WithCancel(srv.serverContext)
	go func() {
		<-connContex.Done()
		conn.Close()
	}()
	conn.Write([]byte{100})
	cons := InitConsumer(conn, out)
	err := cons.Consume()
	if err != nil && connContex.Err() != context.Canceled {
		log.Printf("Unexpected error: %v", err)
	}
	connCancel()
	srv.wg.Done()
}

// Close will politly close the server. It will stop listening and wait
// for connections to close client side. The context used during
// server creation can be used to force close the server.
func (srv *Server) Close() {
	srv.cancelListen()
	srv.ln.Close()
	srv.wg.Wait()
}
