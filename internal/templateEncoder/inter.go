package templateEncoder

import (
	"io"
)

type Template interface {
	Execute(io.Writer, interface{}) error
}
