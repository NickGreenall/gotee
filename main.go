package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NickGreenall/gotee/internal/atomiser"
	"github.com/NickGreenall/gotee/internal/keyEncoding"
	"github.com/NickGreenall/gotee/internal/producer"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	var inStrm io.Reader
	if AmForeground() {
		err := SpawnSniffer("unix", "./test.sock")
		if err != nil {
			fmt.Printf("Unexpected error: %v", err)
			return
		}
		fmt.Println("Foreground")
		inStrm = os.Stdin
	} else {
		inStrm = io.TeeReader(os.Stdin, os.Stdout)
	}
	for !SockOpen("./test.sock") {
		time.Sleep(1)
	}
	conn, err := net.DialTimeout("unix", "./test.sock", 60*time.Second)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}

	scanner := bufio.NewScanner(inStrm)
	enc := json.NewEncoder(conn)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	prod := producer.NewProducer(keyEnc)
	a, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, prod.AtomEnc)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}

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
