package main

import (
	"github.com/NickGreenall/gotee/internal/common"
)

// KeyEncoder is a simpliffied interface for KeyEncoder like
// Objects.
type KeyEncoder interface {
	Encode(string, interface{}) error
	NewEncoderForKey(string) common.Encoder
}

// Producer is the client side of the application. There should
// be only once producer for each client and connection.
type Producer struct {
	// Encoder to connect the atom to. Encodes using the atom
	// key/command
	AtomEnc common.Encoder
	// enc is the KeyEncoder which the producer is attached to.
	enc KeyEncoder
}

// NewProducer creates a new producer object connected to the given encoder.
func NewProducer(keyEncoder KeyEncoder) *Producer {
	return &Producer{
		keyEncoder.NewEncoderForKey("atom"),
		keyEncoder,
	}
}

// SetJson sets the output for the corresponding consumer to JSON.
func (producer *Producer) SetJson() error {
	return producer.enc.Encode("json", nil)
}

// SetTemplate sets the output to be rendered by given template.
func (producer *Producer) SetTemplate(template string) error {
	return producer.enc.Encode("template", template)
}
