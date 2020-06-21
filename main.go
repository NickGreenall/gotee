package main

import "bufio"
import "encoding/json"
import "fmt"
import "github.com/NickGreenall/gotee/internal/atomiser"
import "github.com/NickGreenall/gotee/internal/keyEncoding"
import "io"
import "os"

func main() {
	rdr, wtr := io.Pipe()
	scanner := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(wtr)
	keyEnc := keyEncoding.NewJsonKeyEncoder(enc)
	atomEnc := keyEnc.NewEncoderForKey("atom")
	dec := json.NewDecoder(rdr)
	keyDec := keyEncoding.NewJsonKeyDecoder(dec)
	a, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, atomEnc)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}
	go func() {
		atom := make(map[string][]byte)
		for {
			key, err := keyDec.Pop()
			if err != nil {
				fmt.Printf("Unexpected error: %v", err)
				return
			}
			//fmt.Println(key)
			err = keyDec.Decode(&atom)
			if err != nil {
				fmt.Printf("Unexpected error: %v", err)
				return
			}
			dig, _ := atom["dig"]
			fmt.Printf("%v: %s\n", key, dig)
		}
	}()
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
