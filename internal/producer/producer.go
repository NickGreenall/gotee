package producer

import (
	"github.com/NickGreenall/gotee/internal/common"
)

type KeyEncoder interface {
	Encode(string, interface{}) error
	NewEncoderForKey(string) common.Encoder
}

type Producer struct {
	AtomEnc common.Encoder
	enc     KeyEncoder
}

func NewProducer(keyEncoder KeyEncoder) *Producer {
	return &Producer{
		keyEncoder.NewEncoderForKey("atom"),
		keyEncoder,
	}
}

func (producer *Producer) SetJson() error {
	return producer.enc.Encode("json", nil)
}

func (producer *Producer) SetTemplate(template string) error {
	return producer.enc.Encode("template", template)
}
