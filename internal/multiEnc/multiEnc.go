package multiEnc

import (
	"github.com/NickGreenall/gotee/internal/common"
	"sync"
)

type MultiEncoder struct {
	Key  string
	enc  chan MultiAtom
	err  chan error
	done chan struct{}
}

type MultiAtom struct {
	Key  string
	Data interface{}
}

type MultiEncodeErr struct {
	msg string
}

func (e *MultiEncodeErr) Error() string {
	return e.msg
}

func NewMultiEncoder(key string, done chan struct{}) *MultiEncoder {
	return &MultiEncoder{
		key,
		make(chan MultiAtom),
		make(chan error),
		done,
	}
}

func (e *MultiEncoder) Close() error {
	select {
	case <-e.err:
		return &MultiEncodeErr{"Channels have already been closed"}
	default:
		close(e.enc)
		<-e.err
		return nil
	}
}

func (e *MultiEncoder) Encode(data interface{}) error {
	select {
	case <-e.err:
		return &MultiEncodeErr{"Channels have already been closed"}
	case <-e.done:
		e.Close()
		return &MultiEncodeErr{"Done has already been closed"}
	default:
		atom := MultiAtom{
			e.Key,
			data,
		}
		e.enc <- atom
		err := <-e.err
		return err
	}
}

func Join(enc common.Encoder, done chan struct{}, encoders ...*MultiEncoder) {
	var wg sync.WaitGroup

	outEnc := make(chan interface{})
	outErr := make(chan error)

	fwd := func(e *MultiEncoder) {
		defer wg.Done()
		defer close(e.err)
		for data := range e.enc {
			outEnc <- data
			err := <-outErr
			e.err <- err
		}
	}

	go func() {
		for data := range outEnc {
			err := enc.Encode(data)
			outErr <- err
		}
		close(outErr)
	}()

	wg.Add(len(encoders))
	for _, enc := range encoders {
		go fwd(enc)
	}

	go func() {
		wg.Wait()
		close(outEnc)
	}()
}
