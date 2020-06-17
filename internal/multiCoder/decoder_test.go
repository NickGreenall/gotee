package multiCoder

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"io"
	"sync"
	"testing"
)

type UnmarshalError struct {
	msg string
}

func (e *UnmarshalError) Error() string {
	return e.msg
}

func (atom MultiAtom) MockUnmarshal(v interface{}) error {
	target, ok := v.(*MultiAtom)
	if !ok {
		return &UnmarshalError{"Could not unmarshal into data"}
	}
	*target = atom
	return nil
}

func (atom MultiAtom) GetKey() string {
	return atom.Key
}

func (atom MultiAtom) UnmarshalMulti(v interface{}) error {
	err, ok := atom.Data.(error)
	if ok {
		return err
	}
	s, ok := atom.Data.(string)
	if !ok {
		return &UnmarshalError{"Data not a string"}
	}
	sp, ok := v.(*string)
	if !ok {
		return &UnmarshalError{"Value not a string"}
	}
	*sp = s
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

func checkDecVals(t *testing.T, dec *MultiDecoder, expectedVal []string) {
	for _, expected := range expectedVal {
		val := new(string)
		*val = ""
		err := dec.Decode(val)
		if err != nil {
			t.Errorf("%v: Unexpected error: %v", dec.Key, err)
		} else if *val != expected {
			t.Errorf("%v: Unexpected value: %v, expected: %v", dec.Key, val, expected)
		} else {
			t.Logf("%v: Decoded expected val: %v", dec.Key, *val)
		}
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
		MultiAtom{"A", "3"},
		//MultiAtom{"A", "4"},
	)
	un := &MultiAtom{}
	endA := NewMultiDecoder("A")
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	errChan := JoinDecoders(srcDec, un, done, endA)

	go func() {
		for err := range errChan {
			t.Errorf("Unexpected error: %v", err)
			if err == io.EOF {
				return
			}
		}
	}()
	wg.Add(1)
	go func() {
		checkDecVals(t, endA, expectedVal)
		wg.Done()
	}()
	wg.Wait()
	close(done)
}

func TestMultiDecode(t *testing.T) {
	expectedVal := map[string][]string{
		"A": {"1", "3"},
		"B": {"2", "4"},
	}
	srcDec := mock.NewMockCoder(
		MultiAtom{"A", "1"},
		MultiAtom{"B", "2"},
		MultiAtom{"A", "3"},
		MultiAtom{"B", "4"},
		MultiAtom{"B", "5"},
	)
	un := &MultiAtom{}
	endA := NewMultiDecoder("A")
	endB := NewMultiDecoder("B")
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	defer close(done)

	errChan := JoinDecoders(srcDec, un, done, endA, endB)

	go func(t *testing.T) {
		for err := range errChan {
			t.Errorf("Unexpected error: %v", err)
		}
	}(t)

	wg.Add(2)
	go func() {
		checkDecVals(t, endA, expectedVal["A"])
		wg.Done()
	}()
	go func() {
		checkDecVals(t, endB, expectedVal["B"])
		wg.Done()
	}()
	wg.Wait()
}
