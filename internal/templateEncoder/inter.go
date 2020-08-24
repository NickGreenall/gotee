package templateEncoder

import (
	"io"
)

//Template is a simplified interface for golang stdlib template
//objects.
type Template interface {
	Execute(io.Writer, interface{}) error
}
