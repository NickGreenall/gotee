package main

import "bufio"
import "encoding/json"
import "fmt"
import "github.com/NickGreenall/gotee/internal/atomiser"
import "os"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	a, err := atomiser.NewAtomiser(".*", encoder)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
		return
	}
	for scanner.Scan() {
		_, err := a.Write(scanner.Bytes())
		if err != nil {
			fmt.Printf("Unexpected error: %v", err)
			return
		}
	}
}
