package main

import (
	"fmt"
	"github.com/threft/tidm"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input file given.")
		return
	}
	T, err := tidm.ParseThrift(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("%#v", T)
}
