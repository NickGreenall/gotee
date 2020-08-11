package muxWriter

import (
	"io"
	"sync"
)

type Writer struct {
	done chan bool
	wrt  chan []byte
	rtn  chan writeRtn
}

type writeRtn struct {
	n   int
	err error
}

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

func (m *Mux) Close() {
	close(m.done)
	m.wg.Wait()
	close(m.wrt)
}

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
