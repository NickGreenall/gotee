package multiCoder

import (
	"github.com/NickGreenall/gotee/internal/common"
)

type MultiDecoder struct {
	Key  string
	dec  chan interface{}
	done chan struct{}
}

func NewMultiDecoder(key string) *MultiDecoder {
	return &MultiDecoder{
		key,
		nil,
		nil,
	}
}

func (d *MultiDecoder) Decode() interface{} {
	select {
	case <-d.done:
		return nil
	default:
		return <-d.dec
	}
}

func JoinDecoders(
	dec common.Decoder, done chan struct{}, decoders ...*MultiDecoder,
) chan error {
	outErr := make(chan error)
	outMap := make(map[string]chan interface{})

	for _, decoder := range decoders {
		decoder.dec = make(chan interface{})
		outMap[decoder.Key] = decoder.dec
	}

	go func() {
		defer func() {
			close(outErr)
			for _, dec := range decoders {
				close(dec.dec)
			}
		}()
		atom := &MultiAtom{}
		for {
			select {
			case <-done:
				break
			default:
				atom.Key = ""
				err := dec.Decode(atom)
				if err != nil {
					outErr <- err
				} else {
					outChan, ok := outMap[atom.Key]
					if ok {
						outChan <- atom.Data
					} else {
						outErr <- &MultiDecodeErr{
							"Atom for key which doesn't exist",
						}
					}
				}
			}
		}
	}()
	return outErr
}
