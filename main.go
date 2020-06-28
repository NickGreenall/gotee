package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/consumer"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"github.com/NickGreenall/gotee/internal/producer"
	"io"
	"os"
)

func main() {
	rdr, wtr := io.Pipe()
	scanner := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(wtr)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	prod := producer.NewProducer(keyEnc)

	dec := json.NewDecoder(rdr)
	keyDec := keyEncoding.NewJsonKeyDecoder(dec)
	cons := new(consumer.Consumer)
	cons.Dec = keyDec
	cons.Out = os.Stdout

	a, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, prod.AtomEnc)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}

	go cons.Consume()

	prod.SetJson()
	//prod.SetTemplate("\033[32mdig: {{.dig}}\033[0m\n")
	for scanner.Scan() {
		b := scanner.Bytes()
		//fmt.Printf("Bytes: %s\n", b)
		_, err := a.Write(b)
		if err != nil {
			fmt.Printf("Unexpected error: %v", err)
			return
		}
	}
}
