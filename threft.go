package main

import (
	"fmt"
	"github.com/threft/threft-gen-go/gen-go"
	"github.com/threft/tidm"
	"os"
)

var (
	optionDebugging bool
)

func main() {
	// Options are all hardcoded for now.
	fmt.Println("Debug mode enabled.")
	optionDebugging = true

	if len(os.Args) < 2 {
		fmt.Println("No input file given.")
		return
	}
	T, err := tidm.ParseThrift(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	gen_go.GenerateGo(T)

}
