package multiCoder

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"testing"
)

func (atom MultiAtom) MockUnmarshal(v interface{}) error {
	target, ok := v.(*MultiAtom)
	if !ok {
		return &mock.MockCoderError{}
	}
	*target = atom
	return nil
}

func TestMockDecode(t *testing.T) {
	// Stubed out to test New mock decoder
	dec := mock.NewMockCoder(
		MultiAtom{"A", "A"},
	)
	var bar MultiAtom
	err := dec.Decode(&bar)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if bar.Key != "A" {
		t.Errorf("Unexpected key: %v", bar.Key)
	}
	if bar.Data != "A" {
		t.Errorf("Unexpected data: %v", bar.Data)
	}
}

func TestSingleDecode(t *testing.T) {
	expectedVal := []string{
		"1",
		"2",
	}
	srcDec := mock.NewMockCoder(
		MultiAtom{"A", "1"},
		MultiAtom{"A", "2"},
	)
	endA := NewMultiDecoder("A")
	done := make(chan struct{})
	defer close(done)

	errChan := JoinDecoders(srcDec, done, endA)

	go func(t *testing.T) {
		select {
		case err := <-errChan:
			t.Fatalf("Unexpected error: %v", err)
		case <-done:
		}
	}(t)

	for _, expected := range expectedVal {
		val := endA.Decode()
		if val != expected {
			t.Fatalf("Unexpected value: %v, expected: %v", val, expected)
		}
	}
}
