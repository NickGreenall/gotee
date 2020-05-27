package comm

const (
	WriteCall string = "Write"
)

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

type Atom struct {
	Source string
	Call   string
	Data   interface{}
}

type AtomEncoder struct {
	Name string
	enc  Encoder
}

type AtomDecoder struct {
	Name string
	Dec  Decoder
}

func NewAtomEncoder(name string, enc Encoder) *AtomEncoder {
	return &AtomEncoder{
		name,
		enc,
	}
}
