package templateEncoder

import (
	"io"
)

type TemplateEncoder struct {
	Tmplt Template
	Wrtr  io.Writer
}

func (e *TemplateEncoder) Encode(v interface{}) error {
	return e.Tmplt.Execute(e.Wrtr, v)
}
