package muxWriter

var (
	MuxClosed *MuxWrtrErr = &MuxWrtrErr{
		"Mux is closed and is not writable",
	}
)

type MuxWrtrErr struct {
	msg string
}

func (e *MuxWrtrErr) Error() string {
	return e.msg
}
