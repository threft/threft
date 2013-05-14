package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"github.com/threft/threft-gen-go/gog"
	"github.com/threft/threft/tidm"
)

var options struct {
	Debugging    bool   `short:"d" long:"debug" description:"Enable logging of debug messages to StdOut"`
	InputFiles   string `short:"i" long:"input" description:"Input folders/files"`
	Generator    string `short:"g" long:"gen" description:"Generator to use (for example: go, html), can include arguments for generator"`
	OutputFolder string `short:"o" long:"output" description:"Folder to generate code to"`
}

func main() {
	args, err := flags.Parse(&options)
	if err != nil {
		flagError := err.(*flags.Error)
		if flagError.Type == flags.ErrHelp {
			return
		}
		if flagError.Type == flags.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options.")
			return
		}
		fmt.Printf("Error parsing flags: %s\n", err)
		return
	}
	// options are all hardcoded for now.
	fmt.Println("Debug mode enabled.")
	options.Debugging = true

	// check for unexpected arguments
	if len(args) > 0 {
		fmt.Println("Unknown argument '%s'.\n", args[0])
		return
	}

	// create new TIDM
	t := tidm.NewTIDM()

	// assuming file name is correct and file is existing.
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file. %s\n", err)
		return
	}

	// add document to TIDM
	err = t.AddDocument(filename, file)
	if err != nil {
		fmt.Printf("Error adding document to TIDM: %s\n", err)
		return
	}

	// verify complete TIDM structure (each target, each namespace)
	perr := t.Verify()
	if perr != nil {
		spew.Dump(perr)
	}

	// right now we directly call the gog library, instead of parsing to tidm-json and invoking seperate binary (threft-gen-go)
	gog.GenerateGo(t)
}
