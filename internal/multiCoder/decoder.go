package multiCoder

import (
	"github.com/NickGreenall/gotee/internal/common"
)

type MultiDecoder struct {
	Key  string
	dec  chan interface{}
	err  chan error
	done chan struct{}
}

type Unmarshaler interface {
	UnmarshalMulti(v interface{}) error
	GetKey() string
}

func NewMultiDecoder(key string) *MultiDecoder {
	return &MultiDecoder{
		key,
		nil,
		nil,
		nil,
	}
}

func (d *MultiDecoder) Decode(v interface{}) error {
	select {
	case <-d.done:
		return &MultiDecodeErr{"Decoding cancelled"}
	case d.dec <- v:
		return <-d.err
	}
}

func JoinDecoders(
	dec common.Decoder,
	unmarshaler Unmarshaler,
	done chan struct{},
	decoders ...*MultiDecoder,
) chan error {
	outErr := make(chan error)
	decMap := make(map[string]*MultiDecoder)

	for _, decoder := range decoders {
		decoder.dec = make(chan interface{})
		decoder.err = make(chan error)
		decoder.done = done
		decMap[decoder.Key] = decoder
	}

	go func() {
		defer func() {
			close(outErr)
		}()
		var decoder *MultiDecoder
		var ok bool
		for {
			err := dec.Decode(unmarshaler)
			if err == nil {
				key := unmarshaler.GetKey()
				decoder, ok = decMap[key]
				if !ok {
					err = &MultiDecodeErr{
						"Recieved key which isn't connected",
					}
				}
			}
			if err != nil {
				outErr <- err
			} else {
				select {
				case v := <-decoder.dec:
					err = unmarshaler.UnmarshalMulti(v)
					decoder.err <- err
				case <-done:
					return
				}
			}
		}
	}()
	return outErr
}
