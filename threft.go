package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"github.com/threft/threft-gen-go/gog"
	"github.com/threft/threft/tidm"
	"os"
	"strings"
)

var options struct {
	Debugging    bool     `short:"d" long:"debug" description:"Enable logging of debug messages to StdOut"`
	InputFiles   []string `short:"i" long:"input" description:"Input folders/files"`
	Generator    string   `short:"g" long:"gen" description:"Generator to use (for example: go, html), can include arguments for generator"`
	OutputFolder string   `short:"o" long:"output" description:"Folder to generate code to"`
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

	// create slice to store all filenames in..
	filenames := []string{}

	fmt.Println("Searching for thrift files and setting up documents.")
	for _, filefolder := range options.InputFiles {
		if filefolder[0:1] != string(os.PathSeparator) {
			pwd := os.Getenv("PWD")
			filefolder = pwd + string(os.PathSeparator) + filefolder
		}

		fi, err := os.Stat(filefolder)
		if err != nil {
			fmt.Printf("Error getting info on '%s': %s\n", filefolder, err)
			return
		}

		if fi.IsDir() {
			// remove an eventual path seperator on the right
			filefolder = strings.TrimRight(filefolder, string(os.PathSeparator))

			// setup recursive scan method
			var scanDir func(path string)
			scanDir = func(path string) {
				// open given path
				f, err := os.Open(path)
				if err != nil {
					fmt.Printf("Error opening '%s': %s\n", path, err)
					return
				}

				// read fileInfo for all files/folders
				fis, err := f.Readdir(-1)
				if err != nil {
					fmt.Printf("Error reading dir info on '%s': %s\n", path, err)
				}
				// loop through all files/folders
				for _, fi := range fis {
					foundFile := path + string(os.PathSeparator) + fi.Name()
					if fi.IsDir() {
						// disabled, not doing recursive now..
						// // recursive scan dir
						// err := scanDir(foundFile)
						// if err != nil {
						// 	return err
						// }
					} else if strings.HasSuffix(foundFile, ".thrift") {
						// found a .thrift file
						filenames = append(filenames, foundFile)
					}
				}
				return
			}

			// do recursive file find
			scanDir(filefolder)

			// print findings
			fmt.Printf("Found %d files in given path '%s'.\n", len(filenames), filefolder)
			for _, filename := range filenames {
				fmt.Printf("• %s\n", filename)
			}
			fmt.Println("")
		} else {
			// only one file given
			// check if file is thrift file
			if !strings.HasSuffix(filefolder, ".thrift") {
				fmt.Printf("Error: invalid file extension for '%s' (expected .thrift).", filefolder)
				return
			}

			// add filename to list
			filenames = append(filenames, filefolder)
		}
	}

	// create new TIDM
	t := tidm.NewTIDM()

	// create document for each file found
	for _, filename := range filenames {
		// assuming file name is correct and file is existing.
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error opening file. %s\n", err)
			return
		}
		defer file.Close()

		// add document to TIDM
		err = t.AddDocument(tidm.DocumentName(filename), file)
		if err != nil {
			fmt.Printf("Error adding document to TIDM: %s\n", err)
			return
		}
		file.Close()
	}

	// verify complete TIDM structure (each target, each namespace)
	perr := t.Verify()
	if perr != nil {
		spew.Dump(perr)
	}

	// right now we directly call the gog library, instead of parsing to tidm-json and invoking seperate binary (threft-gen-go)
	gog.GenerateGo(t)
}
