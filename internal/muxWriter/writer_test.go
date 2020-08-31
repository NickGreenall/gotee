package muxWriter

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"
)

func TestSingleWriter(t *testing.T) {
	inputStrings := []string{
		"1\n",
		"2A\n",
		"3FooBar\n",
	}
	expectedCounts := []int{2, 3, 8}
	expectedErr := []error{nil, nil, nil}

	outBuf := new(bytes.Buffer)
	mux := NewMux(context.Background(), outBuf)
	defer mux.Close()
	wrtr := mux.NewWriter()

	for i, input := range inputStrings {
		n, err := wrtr.Write([]byte(input))
		if n != expectedCounts[i] {
			t.Errorf("Expected n=%v, got %v", expectedCounts[i], n)
		}
		if err != expectedErr[i] {
			t.Errorf("Expected err=%v, got %v", expectedErr[i], err)
		}
	}

	for _, input := range inputStrings {
		str, err := outBuf.ReadString(byte('\n'))
		if err != nil {
			t.Error(err)
		}
		if str != input {
			t.Errorf("Unexpected output: %v", str)
		}
	}

}

func TestWriteAfterClose(t *testing.T) {
	outBuf := new(bytes.Buffer)
	mux := NewMux(context.Background(), outBuf)
	wrtr := mux.NewWriter()
	mux.Close()
	n, err := wrtr.Write([]byte("TEST"))
	if n != 0 {
		t.Errorf("Expected no bytes to be written, instead %v", n)
	}
	if err != MuxClosed {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestErrorPropagation(t *testing.T) {
	outRdr, outWrtr := io.Pipe()
	mux := NewMux(context.Background(), outWrtr)
	wrtr := mux.NewWriter()
	defer mux.Close()
	outRdr.Close()
	n, err := wrtr.Write([]byte("TEST"))
	if n != 0 {
		t.Errorf("Expected no bytes to be written, instead %v", n)
	}
	if err != io.ErrClosedPipe {
		t.Errorf("Unexpected error: %v", err)
	}
}
func TestMultipleConcurrentWriters(t *testing.T) {
	inputStringsA := []string{
		"A:1\n",
		"A:2A\n",
		"A:3FooBar\n",
	}
	expectedCountsA := []int{4, 5, 10}
	expectedErrA := []error{nil, nil, nil}

	inputStringsB := []string{
		"B:1\n",
		"B:523\n",
		"B:Bar\n",
	}
	expectedCountsB := []int{4, 6, 6}
	expectedErrB := []error{nil, nil, nil}

	outBuf := new(bytes.Buffer)
	mux := NewMux(context.Background(), outBuf)
	defer mux.Close()
	wrtr := mux.NewWriter()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		for i, input := range inputStringsA {
			n, err := wrtr.Write([]byte(input))
			if n != expectedCountsA[i] {
				t.Errorf("Expected n=%v, got %v", expectedCountsA[i], n)
			}
			if err != expectedErrA[i] {
				t.Errorf("Expected err=%v, got %v", expectedErrA[i], err)
			}
		}
		wg.Done()
	}()

	go func() {
		for i, input := range inputStringsB {
			n, err := wrtr.Write([]byte(input))
			if n != expectedCountsB[i] {
				t.Errorf("Expected n=%v, got %v", expectedCountsB[i], n)
			}
			if err != expectedErrB[i] {
				t.Errorf("Expected err=%v, got %v", expectedErrB[i], err)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	iA, iB := 0, 0
	for {
		str, err := outBuf.ReadString(byte('\n'))
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}
		if str[0] == 'A' {
			if str != inputStringsA[iA] {
				t.Errorf("Unexpected output: %v", str)
			}
			iA += 1
		} else if str[0] == 'B' {
			if str != inputStringsB[iB] {
				t.Errorf("Unexpected output: %v", str)
			}
			iB++
		} else {
			t.Errorf("Unexpected start character: %v", str[0])
		}
	}

	if iA != len(inputStringsA) {
		t.Errorf("Unexpected number of A strings: %v", iA)
	}
	if iB != len(inputStringsB) {
		t.Errorf("Unexpected number of B strings: %v", iB)
	}
}
