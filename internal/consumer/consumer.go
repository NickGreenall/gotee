package consumer

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/common"
	"github.com/NickGreenall/gotee/internal/templateEncoder"
	"io"
	"text/template"
)

type Consumer struct {
	Dec interface {
		Pop() (string, error)
		Decode(interface{}) error
	}
	Out io.Writer
	enc common.Encoder
}

func (c *Consumer) Consume() error {
	for {
		key, err := c.Dec.Pop()
		if err != nil {
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

func (c *Consumer) CreateJsonEnc() {
	c.enc = json.NewEncoder(c.Out)
}

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

func (c *Consumer) HandleAtom() error {
	atom := make(atomiser.AtomData)
	if c.Dec == nil {
		return &ConsumerError{"Decoder not initialised"}
	}
	err := c.Dec.Decode(&atom)
	if err != nil {
		return err
	}
	return c.enc.Encode(atom)
}
