package comm

import (
	"regexp"
)

type AtomData map[string][]byte

type Atomiser struct {
	Parser *regexp.Regexp
	enc    Encoder
}

type AtomiserError struct {
	msg string
}

func (a *AtomiserError) Error() string {
	return a.msg
}

func NewAtomiser(pattern string, enc Encoder) (*Atomiser, error) {
	parser, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	} else {
		return &Atomiser{parser, enc}, nil
	}
}

func (a *Atomiser) Write(b []byte) (int, error) {
	if a.Parser.Match(b) {
		data := make(AtomData)
		grpNames := a.Parser.SubexpNames()
		grps := a.Parser.FindSubmatch(b)
		for i, grp := range grps {
			if i == 0 {
				data["match"] = grp
			} else {
				data[grpNames[i]] = grp
			}
		}
		err := a.enc.Encode(data)
		if err != nil {
			return 0, err
		}
		return len(b), nil
	} else {
		return 0, &AtomiserError{"Not a match"}
	}
}
