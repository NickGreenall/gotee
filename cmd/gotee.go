package main

import "bufio"
import "encoding/json"
import "fmt"
import "github.com/NickGreenall/gotee/internal/atomiser"
import "os"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	a, err := atomiser.NewAtomiser(`(?P<dig>\d+)`, encoder)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}
	for scanner.Scan() {
		b := scanner.Bytes()
		fmt.Printf("Bytes: %s\n", b)
		_, err := a.Write(b)
		if err != nil {
			fmt.Printf("Unexpected error: %v", err)
			return
		}
	}
}
