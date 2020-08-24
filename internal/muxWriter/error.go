//muxWriter defines a writer multiplexer suitable for
//concurrent use. It can be used to write to a single
//writer from multiple gorutines.
package muxWriter

var (
	//MuxClosed will be returned when the mux has
	//been closed and a write is attempted.
	MuxClosed *MuxWrtrErr = &MuxWrtrErr{
		"Mux is closed and is not writable",
	}
)

//MuxWrtrErr is a returned error from a mux child writer.
type MuxWrtrErr struct {
	msg string
}

func (e *MuxWrtrErr) Error() string {
	return e.msg
}
