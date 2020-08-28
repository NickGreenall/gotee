package main

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/common"
	"github.com/NickGreenall/gotee/internal/templateEncoder"
	"io"
	"text/template"
)

// Consumer represents the GoTee server logic for a connection.
// Each connection will have a Producer (client side) and a Consumer
// which interperets commands from the Producer.
// Currently the consumer supports 3 producer commands:
// 	- atom: writes an atom of content parsed by a client to out.
// 	- json: sets the output to use a json encoder.
// 	- template: sets the output to use given template.
//
// Either template or json must be called before atom.
type Consumer struct {
	// Dec is an object which supports the KeyDecoder interface
	// and is used to receive commands.
	Dec interface {
		Pop() (string, error)
		Decode(interface{}) error
	}
	// Out is writer to write received output to.
	Out io.Writer
	// enc is the encoder user to write output content. It is
	// attached to Out
	enc common.Encoder
}

// Consume is the main server loop which consumes commands over a connection.
func (c *Consumer) Consume() error {
	for {
		key, err := c.Dec.Pop()
		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}
		switch key {
		case "atom":
			err = c.HandleAtom()
			if err != nil {
				return err
			}
		case "template":
			err = c.CreateTemplateEnc()
			if err != nil {
				return err
			}
		case "json":
			c.CreateJsonEnc()
		default:
			return &ConsumerError{
				"Unexpected key",
			}
		}
	}
}

// CreateJsonEnc sets the consumer output encoder to JSON and
// attaches it to Out. Atoms will then be displayed on the ouptu
// in JSON format.
func (c *Consumer) CreateJsonEnc() {
	c.enc = json.NewEncoder(c.Out)
}

// CreateTemplateEnc sets up the encoder to render using the
// sent template. If the sent template is invalid or can't be
// decoded, will return an error.
func (c *Consumer) CreateTemplateEnc() error {
	tmpltStr := new(string)
	err := c.Dec.Decode(tmpltStr)
	if err != nil {
		return err
	}
	template, err := template.New("out").Parse(*tmpltStr)
	if err != nil {
		return err
	}
	c.enc = &templateEncoder.TemplateEncoder{
		template,
		c.Out,
	}
	return nil
}

// HandleAtom decodes and renders the sent atom. Will return an
// error if either the atom can't be decoded or the output hasn't
// been setup correctly.
func (c *Consumer) HandleAtom() error {
	atom := make(atomiser.AtomData)
	if c.Dec == nil {
		return &ConsumerError{"Decoder not initialised"}
	}
	err := c.Dec.Decode(&atom)
	if err != nil {
		return err
	}
	if c.enc == nil {
		return &ConsumerError{"Output encoder not intialised"}
	}
	return c.enc.Encode(atom)
}
