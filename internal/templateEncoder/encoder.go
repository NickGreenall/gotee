//templateEncoder implements an encoder which uses a
//given template to encode values.
package templateEncoder

import (
	"io"
)

//TemplateEncoder encodes values using Tmplt into Wrtr.
type TemplateEncoder struct {
	//Tmplt is used to define how values are written to Wrtr.
	Tmplt Template
	//Wrtr is the output writer.
	Wrtr io.Writer
}

//Encode runs v through the template, writing the output to
//Wrtr.
func (e *TemplateEncoder) Encode(v interface{}) error {
	return e.Tmplt.Execute(e.Wrtr, v)
}
