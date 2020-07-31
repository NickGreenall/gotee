package main

import (
	"encoding/json"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"github.com/NickGreenall/gotee/internal/mock"
	"github.com/NickGreenall/gotee/internal/templateEncoder"
	"reflect"
	"testing"
)

func TestInitJson(t *testing.T) {
	inputVals := []interface{}{
		mock.MockPacket{"json", nil},
	}
	dec := mock.NewMockCoder(inputVals...)
	kDec := &keyEncoding.KeyDecoder{dec, new(mock.MockPacket)}
	expectedType := reflect.TypeOf(new(json.Encoder))

	consumer := &Consumer{
		kDec,
		nil,
		nil,
	}

	err := consumer.Consume()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	actualType := reflect.TypeOf(consumer.enc)
	if actualType != expectedType {
		t.Errorf("Unexpected encoder type: %s", actualType)
	}
}

func TestInitTemplate(t *testing.T) {
	inputVals := []interface{}{
		mock.MockPacket{"template", "{{.match}}"},
	}
	dec := mock.NewMockCoder(inputVals...)
	kDec := &keyEncoding.KeyDecoder{dec, new(mock.MockPacket)}
	expectedType := reflect.TypeOf(new(templateEncoder.TemplateEncoder))

	consumer := &Consumer{
		kDec,
		nil,
		nil,
	}

	err := consumer.Consume()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	actualType := reflect.TypeOf(consumer.enc)
	if actualType != expectedType {
		t.Errorf("Unexpected encoder type: %s", actualType)
	}
}

func TestAtomReceive(t *testing.T) {
	inputData := make(atomiser.AtomData)
	inputData["match"] = "foo bar"
	inputVals := []interface{}{
		mock.MockPacket{"atom", inputData},
	}
	dec := mock.NewMockCoder(inputVals...)
	kDec := &keyEncoding.KeyDecoder{dec, new(mock.MockPacket)}
	enc := mock.NewMockCoder(nil)

	consumer := &Consumer{
		kDec,
		nil,
		enc,
	}
	err := consumer.Consume()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(inputData, enc.Calls[0]) {
		t.Errorf("Data wasn't successfully encoded")
	}
}
