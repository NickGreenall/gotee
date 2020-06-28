package atomiser

import (
	"github.com/NickGreenall/gotee/internal/mock"
	"reflect"
	"testing"
)

func TestAtomiserNewError(t *testing.T) {
	a, err := NewAtomiser(`(?PBAD)`, nil)
	if a != nil {
		t.Error("Expected new atomiser to be nil")
	}
	if err == nil {
		t.Fatal("Expected regexp error")
	}
}

func TestAtomiserMatchWrite(t *testing.T) {
	mockInputs := []string{
		`match`,
		`123 abc`,
	}
	mockPats := []string{
		`\w+`,
		`(?P<a>\d+)\s(?P<b>\w+)`,
	}
	expectedData := []*AtomData{
		{
			"match": "match",
		},
		{
			"match": "123 abc",
			"a":     "123",
			"b":     "abc",
		},
	}
	for i, input := range mockInputs {
		t.Logf(`input %d: %v, pat: "%v"`, i, input, mockPats[i])
		mockEnc := mock.NewMockCoder(nil)
		atmsr, err := NewAtomiser(
			mockPats[i],
			mockEnc,
		)
		if err != nil {
			t.Fatalf("Could not create test pattern: %v\n", err)
		}
		n, err := atmsr.Write([]byte(input))
		if err != nil {
			t.Errorf(`Unexpected error: %v`, err)
		}
		if n != len(input) {
			t.Errorf(
				"Returned write count does not match length of input. Actual: %v, Expected: %v\n",
				n,
				len(input),
			)
		}
		actual, ok := mockEnc.Calls[0].(AtomData)
		if !ok {
			t.Error(`Value written isn't correct data type`)
		}
		if reflect.DeepEqual(actual, expectedData[i]) {
			t.Errorf(
				"Encoded value doesn't match expected\n\tActual: %v\n\tExpected: %v\n",
				actual,
				expectedData[i],
			)
		}
		t.Logf("input %d: Pass", i)
	}
}

func TestAtomiserNoMatch(t *testing.T) {
	mockInput := []byte(`not a match`)
	mockPat := `MATCH`
	atmsr, err := NewAtomiser(
		mockPat,
		nil,
	)
	if err != nil {
		t.Errorf(`Unexpected error: %v`, err)
	}
	_, err = atmsr.Write(mockInput)
	if err == nil {
		t.Error(`Expected an error, received nil`)
	}
	e, ok := err.(*AtomiserError)
	if !ok {
		t.Error(`Expected an atomiser error\n`)
	}
	if e.Error() != "Not a match" {
		t.Errorf(`Unexpected error message, received: %v`, e.Error())
	}
}

func TestAtomiserEncError(t *testing.T) {
	mockInput := []byte(`match`)
	mockPat := `match`
	mockEnc := mock.NewMockCoder(&mock.MockEncoderError{})
	atmsr, err := NewAtomiser(
		mockPat,
		mockEnc,
	)
	if err != nil {
		t.Errorf(`Unexpected error: %v`, err)
	}
	_, err = atmsr.Write(mockInput)
	if err == nil {
		t.Error(`Expected an error, received nil`)
	}
	_, ok := err.(*AtomiserError)
	if ok {
		t.Error(`Unexpected an atomiser error\n`)
	}
	_, ok = err.(*mock.MockEncoderError)
	if !ok {
		t.Error(`Expected an encoder error\n`)
	}
}
