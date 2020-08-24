package muxWriter

import (
	"io"
	"sync"
)

//Writer is io.writer which can be used to write to the
//mux.
type Writer struct {
	done chan bool
	wrt  chan []byte
	rtn  chan writeRtn
}

type writeRtn struct {
	n   int
	err error
}

//Mux is a multiplexer for a single io.Writer. Child Writers
//can be creater which can write to the writer concurrently.
type Mux struct {
	wrtr io.Writer
	wg   *sync.WaitGroup
	done chan bool
	wrt  chan []byte
	rtn  chan writeRtn
}

func (m *Mux) write() {
	for data := range m.wrt {
		n, err := m.wrtr.Write(data)
		m.rtn <- writeRtn{n, err}
	}
}

//NewMux returns a Mux object for wrtr.
func NewMux(wrtr io.Writer) *Mux {
	m := &Mux{
		wrtr,
		new(sync.WaitGroup),
		make(chan bool),
		make(chan []byte),
		make(chan writeRtn),
	}
	go m.write()
	return m
}

//Close closes the mux. This should be called to ensure closure
//of the mux and underlying channels/gorutines.
func (m *Mux) Close() {
	close(m.done)
	m.wg.Wait()
	close(m.wrt)
}

//Write implements io.Writer write method. Writes data to the
//mux and returns written number of bytes/error.
func (w *Writer) Write(data []byte) (n int, err error) {
	select {
	case <-w.done:
		return 0, MuxClosed
	default:
		w.wrt <- data
		rtn := <-w.rtn
		return rtn.n, rtn.err
	}
}

func (w *Writer) forward(m *Mux) {
	for data := range w.wrt {
		m.wrt <- data
		rtn := <-m.rtn
		w.rtn <- rtn
	}
	m.wg.Done()
}

func (w *Writer) cleanup() {
	<-w.done
	close(w.wrt)
}

//NewWriter creates a new writer for the mux. This writer
//can then be used for concurrent writes.
func (m *Mux) NewWriter() *Writer {
	w := &Writer{
		m.done,
		make(chan []byte),
		make(chan writeRtn),
	}
	m.wg.Add(1)
	go w.forward(m)
	go w.cleanup()
	return w
}
