package multiCoder

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"testing"
)

func TestSingleDecode(t *testing.T) {
	// Stubed out to test New mock decoder
	dec := mock.NewMockCoder(
		mock.MockVal{"A"},
		mock.MockVal{"B"},
	)
	var bar mock.MockVal
	dec.Decode(&bar)
	if bar.Val != "A" {
		t.Fatal("failed")
	}
}
