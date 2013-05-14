package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/threft/threft-gen-go/gog"
	"github.com/threft/threft/tidm"
)

var options struct {
	Debugging bool `short:"d" long:"debug" `
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
	// Options are all hardcoded for now.
	fmt.Println("Debug mode enabled.")
	options.Debugging = true

	if len(args) == 0 {
		fmt.Println("No input file given.")
		return
	}
	if len(args) > 1 {
		fmt.Println("Can only parse one file at this point in development.")
		return
	}

	// Parse thrift definition
	//++ TODO: Use tidm.NewEmptyTidm() to get a tidm.TIDM. Then invoke .ParseThrift() method on tidm.TIDM with an io.Reader
	T, err := tidm.ParseThrift(args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// right now we directly call the gog library, instead of parsing to tidm-json and invoking seperate binary (threft-gen-go)
	gog.GenerateGo(T)
}
