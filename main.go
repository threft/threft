package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/threft/threft/tidm"
	"os"
	"os/exec"
	"path"
	"strings"
)

var options struct {
	Debugging  bool     `short:"d" long:"debug" description:"Enable logging of debug messages to StdOut"`
	InputFiles []string `short:"i" long:"input" description:"Input folders/files"`
	Generator  string   `short:"g" long:"gen" description:"Generator to use (for example: go, html), can include arguments for generator"`
	OutputDir  string   `short:"o" long:"output" description:"Folder to generate code to"`
	DumpTIDM   bool     `long:"dump-tidm" description:"Dumps TIDM structure to ./tidm_dump"`
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

	// check for unexpected arguments
	if len(args) > 0 {
		fmt.Println("Unknown argument '%s'.\n", args[0])
		return
	}

	// hardcode debugging enable
	fmt.Println("Debug mode enabled, hardcoded in code.")
	options.Debugging = true

	var outputDir string
	if options.OutputDir[0] == '/' {
		outputDir = options.OutputDir
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Unable to get wd: %s\n", err)
			return
		}
		outputDir = path.Join(wd, options.OutputDir)
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
		lastPathSeperator := strings.LastIndex(filename, string(os.PathSeparator))
		documentNameString := filename[lastPathSeperator+1:]
		err = t.AddDocument(tidm.DocumentName(documentNameString), file)
		if err != nil {
			fmt.Printf("Error adding document to TIDM: %s\n", err)
			return
		}
		file.Close()
	}

	// parse complete TIDM structure (each document, each target, each namespace)
	perr := t.Parse()
	if perr != nil {
		fmt.Printf("\nError in %s:%d\n \t%s\n", perr.DocLine.DocumentName, perr.DocLine.Line, perr.Message)
	}

	// do a TIDM dump if requested by user
	if options.DumpTIDM {
		err = dumpTIDM(t)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// get generator fields (possibly options)
	genFields := strings.Fields(options.Generator)
	if len(genFields) == 0 {
		fmt.Println("No generator given. Can not continue. Use -g to generate code.")
		return
	}

	// prepare generator command
	genCmd := exec.Command("threft-gen-"+genFields[0], genFields[1:]...)
	genCmd.Dir = outputDir
	genCmd.Stderr = os.Stderr
	genCmd.Stdout = os.Stdout

	// get stdinPipe to send json when process has started
	stdinPipe, err := genCmd.StdinPipe()
	if err != nil {
		fmt.Printf("Error getting stdin pipe: %s\n", err)
		return
	}

	// start generator
	err = genCmd.Start()
	if err != nil {
		fmt.Printf("Error on starting generator: %s\n", err)
		return
	}

	// write tidm-json to generator
	err = t.WriteTo(stdinPipe)
	if err != nil {
		fmt.Printf("Error writing data to generator: %s\n", err)
		return
	}

	// close the stdinPipe
	err = stdinPipe.Close()
	if err != nil {
		fmt.Printf("Error closing stdin pipe: %s\n", err)
	}

	// wait for generator to exit
	err = genCmd.Wait()
	if err != nil {
		fmt.Printf("Error while running generator: %s\n", err)
		return
	}

	fmt.Println("All done.")
}
