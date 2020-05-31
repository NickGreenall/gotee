package comm

import (
	"sync"
	"testing"
)

func checkEncVals(t *testing.T, enc *MockEncoder, atoms ...MultiAtom) {
	for i, atom := range atoms {
		a, ok := enc.calls[i].(MultiAtom)
		if !ok {
			t.Error("Expected multi atom type")
		}
		if a.Key != atom.Key {
			t.Errorf("Unexpected key, actual: %v, expected: %v", a.Key, atom.Key)
		}
		v, ok := a.Data.(string)
		if !ok {
			t.Error("Expected data type")
		}
		expected := atom.Data.(string)
		if v != expected {
			t.Errorf("Unexpected write value, actual: %v, expected: %v", v, expected)
		}
	}
}

func checkEncValsUnordered(t *testing.T, enc *MockEncoder, atoms ...MultiAtom) {
	callMap := make(map[MultiAtom]bool)
	for _, call := range enc.calls {
		callMap[call.(MultiAtom)] = true
	}
	for _, atom := range atoms {
		_, ok := callMap[atom]
		if !ok {
			t.Errorf("Expected atom to be present %v", atom)
		}
	}
}

func TestSingleEncode(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mltEnc := NewMultiEncoder("test", done)
	mockEnc := NewMockEncoder(nil)
	Join(mockEnc, done, mltEnc)
	err := mltEnc.Encode("Test Write")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	checkEncVals(t, mockEnc, MultiAtom{"test", "Test Write"})
}

func TestSeqMultiEncode(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mltEncA := NewMultiEncoder("A", done)
	mltEncB := NewMultiEncoder("B", done)
	mockEnc := NewMockEncoder(nil, nil)
	Join(mockEnc, done, mltEncA, mltEncB)
	err := mltEncA.Encode("A")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = mltEncB.Encode("B")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	checkEncVals(t, mockEnc, MultiAtom{"A", "A"}, MultiAtom{"B", "B"})
}

func TestConMultiEncode(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mltEncA := NewMultiEncoder("A", done)
	mltEncB := NewMultiEncoder("B", done)
	mockEnc := NewMockEncoder(nil, nil)
	Join(mockEnc, done, mltEncA, mltEncB)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := mltEncA.Encode("A")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		wg.Done()
	}()
	go func() {
		err := mltEncB.Encode("B")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
	checkEncValsUnordered(t, mockEnc, MultiAtom{"A", "A"}, MultiAtom{"B", "B"})
}

func TestConMultiEncodeErr(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mltEncA := NewMultiEncoder("A", done)
	mltEncB := NewMultiEncoder("B", done)
	mockEnc := NewMockEncoder(&MultiEncodeErr{}, &MultiEncodeErr{})
	Join(mockEnc, done, mltEncA, mltEncB)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := mltEncA.Encode("A")
		if err == nil {
			t.Errorf("Expected error")
		}
		wg.Done()
	}()
	go func() {
		err := mltEncB.Encode("B")
		if err == nil {
			t.Errorf("Expected error")
		}
		wg.Done()
	}()
	wg.Wait()
}

func TestSingleEncodeClosed(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	mltEnc := NewMultiEncoder("test", done)
	mockEnc := NewMockEncoder(nil)
	Join(mockEnc, done, mltEnc)
	err := mltEnc.Close()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = mltEnc.Encode("Fail encoding")
	_, ok := err.(*MultiEncodeErr)
	if !ok {
		t.Error("Expected closed error")
	}
	t.Logf("Error message: %v", err)
}

func TestSingleEncodeDone(t *testing.T) {
	done := make(chan struct{})
	mltEnc := NewMultiEncoder("test", done)
	mockEnc := NewMockEncoder(nil)
	Join(mockEnc, done, mltEnc)
	close(done)
	err := mltEnc.Encode("Fail encoding")
	_, ok := err.(*MultiEncodeErr)
	if !ok {
		t.Error("Expected closed error")
	}
	t.Logf("Error message: %v", err)
}
