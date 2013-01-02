package main

import (
	"fmt"
	"github.com/threft/tidm"
	"os"
)

var (
	optionPrintJson bool
	optionDebugging bool
)

func main() {
	// These options are all hardcoded for now.
	fmt.Println("No output given. Will print tidm-json.")
	optionPrintJson = true
	fmt.Println("Debug mode enabled.")
	optionDebugging = true

	fmt.Println("Pretty json for debugging hardcoded by importing launchpad.net/rjson as json package.")

	if len(os.Args) < 2 {
		fmt.Println("No input file given.")
		return
	}
	T, err := tidm.ParseThrift(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("%#v\n", T)

}
