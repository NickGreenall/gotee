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
	serverContext   context.Context
	listenerContext context.Context
	cancelListen    context.CancelFunc
	ln              net.Listener
	wg              *sync.WaitGroup
}

func NewServer(serverContext context.Context, network string, address string) (*Server, error) {
	listenerContext, cancelListen := context.WithCancel(serverContext)
	ln, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		serverContext,
		listenerContext,
		cancelListen,
		ln,
		new(sync.WaitGroup),
	}
	return srv, nil
}

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

func (srv *Server) Close() {
	srv.cancelListen()
	srv.ln.Close()
	srv.wg.Wait()
}
