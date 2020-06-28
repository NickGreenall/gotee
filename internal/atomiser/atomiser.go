// atomiser provides a writer to encode regex patterns.
package atomiser

import (
	"github.com/NickGreenall/gotee/internal/common"
	"regexp"
)

// Atomiser is a  writer which encodes parsed regex groups from written bytes.
// Use NewAtomiser to construct.
type Atomiser struct {
	Parser *regexp.Regexp
	Enc    common.Encoder
}

// AtomiserError is raised during write.
// Currently only occurs when no regex match is found.
type AtomiserError struct {
	msg string
}

func (a *AtomiserError) Error() string {
	return a.msg
}

// NewAtomiser constructs a new atomiser from a regex string and encoder.
// If the regex pattern can't be compiled returns the regex error.
func NewAtomiser(pattern string, enc common.Encoder) (*Atomiser, error) {
	parser, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	} else {
		return &Atomiser{parser, enc}, nil
	}
}

// Write to the underlying encoder.
// Written bytes are parsed by the atomiser regex pattern.
// A AtomData is encoded, with a generic "match" key and
// a key for each subgroup in the regex.
// Returns the number of bytes written (should allways
// be equal to the length of the input) if no error.
// Otherwise returns an encoding error or AtomiserError
func (a *Atomiser) Write(b []byte) (int, error) {
	if a.Parser.Match(b) {
		data := make(AtomData)
		grpNames := a.Parser.SubexpNames()
		grps := a.Parser.FindSubmatch(b)
		for i, grp := range grps {
			if i == 0 {
				data["match"] = string(grp)
			} else {
				data[grpNames[i]] = string(grp)
			}
		}
		err := a.Enc.Encode(data)
		if err != nil {
			return 0, err
		}
		return len(b), nil
	} else {
		return 0, &AtomiserError{"Not a match"}
	}
}
