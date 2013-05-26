package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/threft/threft/tidm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var options struct {
	Debugging  bool     `short:"d" long:"debug" description:"Enable logging of debug messages to StdOut"`
	InputFiles []string `short:"i" long:"input" description:"Input folders/files"`
	Generator  string   `short:"g" long:"gen" description:"Generator to use (for example: go, html), can include arguments for generator"`
	OutputDir  string   `short:"o" long:"output" description:"Folder to generate code to"`
	DumpTIDM   bool     `long:"dump-tidm" description:"Dumps TIDM structure to ./tidm_dump"`
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}

func scanDir(path string) (filenames []string) {
	// TODO if turning this function back into a recursive function,
	// either pass an initial filenames slice as an argument, or
	// append() all resulting slices into one.

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
		return
	}
	// loop through all files/folders
	for _, fi := range fis {
		foundFile := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			// disabled, not doing recursive now..
			// // recursive scan dir
			// err := scanDir(foundFile)
			// if err != nil {
			//	return err
			// }
		} else if strings.HasSuffix(foundFile, ".thrift") {
			// found a .thrift file
			filenames = append(filenames, foundFile)
		}
	}
	return
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
			os.Exit(1)
		}
		fmt.Printf("Error parsing flags: %s\n", err)
		os.Exit(1)
	}

	// check for unexpected arguments
	if len(args) > 0 {
		fmt.Printf("Unknown argument '%s'.\n", args[0])
		os.Exit(1)
	}

	// hardcode debugging enable
	fmt.Println("Debug mode enabled, hardcoded in code.")
	options.Debugging = true

	outputDir, err := filepath.Abs(options.OutputDir)
	if err != nil {
		exitWithError("Error getting absolute path for '%s': %s\n", options.OutputDir, err)
	}

	// create slice to store all filenames in..
	filenames := []string{}

	fmt.Println("Searching for thrift files and setting up documents.")
	for _, filefolder := range options.InputFiles {
		filefolder, err = filepath.Abs(filefolder)
		if err != nil {
			exitWithError("Error getting absolute path for '%s': %s\n", filefolder, err)
		}

		fi, err := os.Stat(filefolder)
		if err != nil {
			exitWithError("Error getting info on '%s': %s\n", filefolder, err)
		}

		if fi.IsDir() {
			// do recursive file find
			filenames = scanDir(filefolder)

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
				exitWithError("Error: invalid file extension for '%s' (expected .thrift).\n", filefolder)
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
			exitWithError("Error opening file. %s\n", err)
		}
		defer file.Close()

		// add document to TIDM
		lastPathSeperator := strings.LastIndex(filename, string(os.PathSeparator))
		documentNameString := filename[lastPathSeperator+1:]
		err = t.AddDocument(tidm.DocumentName(documentNameString), file)
		if err != nil {
			exitWithError("Error adding document to TIDM: %s\n", err)
		}
		file.Close()
	}

	// parse complete TIDM structure (each document, each target, each namespace)
	perr := t.Parse()
	if perr != nil {
		exitWithError("\nError at %s\n \t%s\n", perr.DocLine, perr.Message)
	}

	// do a TIDM dump if requested by user
	if options.DumpTIDM {
		err = dumpTIDM(t)
		if err != nil {
			// TODO output a formatted message like in all other
			// cases?
			exitWithError("%s\n", err)
		}
	}

	// get generator fields (possibly options)
	genFields := strings.Fields(options.Generator)
	if len(genFields) == 0 {
		exitWithError("%s", "No generator given. Can not continue. Use -g to generate code.\n")
	}

	// prepare generator command
	genCmd := exec.Command("threft-gen-"+genFields[0], genFields[1:]...)
	genCmd.Dir = outputDir
	genCmd.Stderr = os.Stderr
	genCmd.Stdout = os.Stdout

	// get stdinPipe to send json when process has started
	stdinPipe, err := genCmd.StdinPipe()
	if err != nil {
		exitWithError("Error getting stdin pipe: %s\n", err)
	}

	// start generator
	err = genCmd.Start()
	if err != nil {
		exitWithError("Error on starting generator: %s\n", err)
	}

	// write tidm-json to generator
	err = t.EncodeTo(stdinPipe)
	if err != nil {
		exitWithError("Error writing data to generator: %s\n", err)
	}

	// close the stdinPipe
	err = stdinPipe.Close()
	if err != nil {
		fmt.Printf("Error closing stdin pipe: %s\n", err)
	}

	// wait for generator to exit
	err = genCmd.Wait()
	if err != nil {
		exitWithError("Error while running generator: %s\n", err)
	}

	fmt.Println("All done.")
}
