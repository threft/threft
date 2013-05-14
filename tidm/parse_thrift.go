package tidm

import (
	"fmt"
	"os"
	"strings"
)

func ParseThrift(filefolder string) (*TIDM, error) {
	T := newTIDM()

	fmt.Println("== Searching for thrift files and setting up documents.")
	if filefolder[0:1] != string(os.PathSeparator) {
		pwd := os.Getenv("PWD")
		filefolder = pwd + string(os.PathSeparator) + filefolder
	}

	filenames := make([]string, 0, 1)

	fi, err := os.Stat(filefolder)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		// BasePath is given folder
		T.BasePath = filefolder

		// Remove an eventual path seperator on the right
		filefolder = strings.TrimRight(filefolder, string(os.PathSeparator))

		// Setup recursive scan method
		var scanDir func(name string) error
		scanDir = func(name string) error {
			f, err := os.Open(name)
			if err != nil {
				return err
			}

			// Read fileInfo for all files/folders
			fis, err := f.Readdir(-1)
			if err != nil {
				return err
			}
			// Loop through all files/folders
			for _, fi := range fis {
				foundFile := name + string(os.PathSeparator) + fi.Name()
				if fi.IsDir() {
					// recursive scan dir
					err := scanDir(foundFile)
					if err != nil {
						return err
					}
				} else if strings.HasSuffix(foundFile, ".thrift") {
					// Found a .thrift file.
					filenames = append(filenames, foundFile)
				}
			}
			return nil
		}
		// Do recursive file find
		err := scanDir(filefolder)
		if err != nil {
			return nil, err
		}
	} else {
		// Only one file given.
		// Check if file is thrift file. Error if not.
		if !strings.HasSuffix(filefolder, ".thrift") {
			return nil, fmt.Errorf("Invalid file extension for '%s' (*.thrift expected).", filefolder)
		}

		// Add filename to list
		filenames = append(filenames, filefolder)

		// Set the T.BasePath
		if ind := strings.LastIndex(filefolder, string(os.PathSeparator)); ind > 0 {
			T.BasePath = filefolder[:ind+1]
		} else {
			panic("unexpected behaviour. filefolder should have absolute path and therefore at least one os.PathSeperator")
		}
	}

	// Craete document for each file found.
	for _, filename := range filenames {
		// Create a new Document on TIDM
		T.newDocumentFromFile(filename)
	}

	fmt.Println("---- Report: --")
	{
		fmt.Printf("Found %d documents in base path '%s':\n", len(filenames), T.BasePath)
		for docName, _ := range T.Documents {
			fmt.Printf("• %s\n", docName)
		}
	}
	fmt.Println("----- Done ------\n")

	/**
	 * Parse headers
	 */
	fmt.Println("== Parsing document headers.")
	for _, doc := range T.Documents {
		fmt.Printf("= Document '%s': \n", doc.Name)
		doc.parseDocumentHeaders()
	}
	fmt.Println("---- Report: --")
	fmt.Println("----- Done ------\n")

	/**
	 * Connect documents/namespaces
	 */
	fmt.Println("== Connecting documents and namespaces.")
	// Loop through all docs and make sure that every document has a namespace in each target.
	for _, tar := range T.Targets {
		for _, doc := range T.Documents {
			// Find if this document already has a namespace for this target.
			if _, namespaceExists := doc.NamespaceByTargetName[tar.Name]; namespaceExists {
				continue // Check next target
			}

			// Document has no namespace in this Target
			// Create NamespaceName for this document
			properNamespaceName := NamespaceName(doc.Name)
			// Find if there is already a Namespace with this NamespaceName in this Target
			properNamespace, properNamespaceExists := tar.Namespaces[properNamespaceName]
			if !properNamespaceExists {
				// Create a namespace for this document in this target
				properNamespace, err = tar.createNamespace(properNamespaceName)
				if err != nil {
					panic(err)
				}
			}

			// This document describes the namespace by document name (hence value is false).
			properNamespace.Documents[doc] = false
			doc.NamespaceByTargetName[tar.Name] = properNamespace
		}
	}
	fmt.Println("---- Report: --")
	{
		rowFormat := fmt.Sprintf("%%-%ds %%s\n", T.DocumentNameMaxLength)
		for tarName, _ := range T.Targets {
			fmt.Printf("Target '%s': \n", tarName)
			fmt.Printf(rowFormat, "document:", "namespace:")
			for _, doc := range T.Documents {
				ns, _ := doc.NamespaceByTargetName[tarName]
				fmt.Printf(rowFormat, doc.Name, ns.Name)
			}
			fmt.Printf(rowFormat, "---", "---")
		}
	}
	fmt.Println("----- Done ------\n")

	/**
	 * Parse document definitions
	 */
	fmt.Println("== Parsing document definitions.")
	// parse all definitions
	for _, doc := range T.Documents {
		fmt.Printf("= Document '%s': \n", doc.Name)
		doc.parseDocumentDefinitions()
	}
	fmt.Println("---- Report: --")
	fmt.Println("----- Done ------\n")

	// TIDM is created successfully
	return T, nil
}
