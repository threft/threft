package tidm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type DocumentParseState int

const (
	DocumentParseState_None = DocumentParseState(iota)
	DocumentParseState_Added
	DocumentParseState_HeaderParsed
	DocumentParseState_BodyParsed
)

const (
	NoDefinitionsInDocument = -1
)

type DocumentName string

type Document struct {
	T                     *TIDM
	Name                  DocumentName              // The name of this document. Is also used as NamespaceName for the default Target. Take TIDM.BasePath + Document.Name + ".thrift" = absolute path.
	NamespaceByTargetName map[TargetName]*Namespace // Contains a list of namespaces for the given document. Once the document is parsed: at least one namespace is in this map, target: default.
	DefaultNamespace      *Namespace                // Direct link to the default target namespace. This namespace is unique for this document and is ONLY described by this document.

	// Source management
	lines                []string // All source lines for this document. Whitespace is trimmed.
	lastParsedLineNumber int      // line number of the last parsed line. Used by NextMeaningfulLine().
}

func (T *TIDM) newDocumentFromFile(filename string) *Document {
	var err error

	docName := DocumentName(strings.Replace(filename[len(T.BasePath):len(filename)-7], "/", ".", -1))
	if len(docName) > T.DocumentNameMaxLength {
		T.DocumentNameMaxLength = len(docName)
	}

	// Create a new document
	doc := &Document{
		T:    T,
		Name: docName,
		NamespaceByTargetName: make(map[TargetName]*Namespace),

		lastParsedLineNumber: -1, // defaults to -1. As the first line is actually line number 0.
	}

	// Get the namespace for default target
	// In this case it is newly created.
	// It should be imposible that a Namespace is already existing for doc.Name in the Target_default
	target, exists := T.Targets[Target_default]
	if !exists {
		target, err = T.createTarget(Target_default)
		if err != nil {
			panic(err)
		}
	}

	ns, err := target.createNamespace(NamespaceName(doc.Name))
	if err != nil {
		panic(err)
	}

	// Directly link to default namespace (as that namespace is only for this document)
	doc.DefaultNamespace = ns

	// Add the default Namespace is to the Document.Namespaces list
	doc.NamespaceByTargetName[Target_default] = ns

	// Connect this Document on the Namespace
	ns.Documents[doc] = true

	// Connect this Document on the TIDM
	T.Documents[doc.Name] = doc

	// Assuming file name is correct and file is existing.
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	sourcereader := bufio.NewReader(file)
	for {
		// Blegh, this for loop feels like a trampoline. The EOF pickup could be nicer..
		// Also: ReadString('\n') isn't really cross-platform, is it?
		line, err := sourcereader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					goto addLine
				}
				break
			}
			fmt.Printf("Error while reading line %d from file %s. %s\n", len(doc.lines)+1, filename, err)
			break
		}
	addLine:
		doc.lines = append(doc.lines, strings.TrimSpace(line))
	}

	return doc
}

func (doc *Document) NextMeaningfulLine() (line string, ok bool) {
	for {
		// Check if last line has been given already
		if len(doc.lines)-1 == doc.lastParsedLineNumber {
			return "", false
		}

		// Fetch next line
		doc.lastParsedLineNumber++
		line = doc.lines[doc.lastParsedLineNumber]

		// Remove eventual comment from line
		pos := strings.Index(line, "#")
		if pos > -1 {
			line = line[:pos]
		}
		pos = strings.Index(line, "//")
		
		if pos > -1 {
			line = line[:pos]
		}

		// Trim space and list seperators from line
		line = strings.TrimSpace(line)
		line = strings.TrimRight(line, ",; ")

		// ignore line if it is empty
		if len(line) == 0 {
			// try next line
			continue
		}
		return line, true
	}
	panic("unreachable")
}

func (doc *Document) parseDocumentHeaders() {
	var err error

	// loop through lines
	for {

		line, ok := doc.NextMeaningfulLine()
		if !ok {
			break // no new lines
		}

		// get fields from line
		fields := strings.Fields(line)

		switch fields[0] {
		case "include":
			// Not supporting cross-document references yet.
			fmt.Printf("Ignoring include statement '%s' in document '%s'\n", line, doc.Name)
			continue // try next line

		case "cpp_include":
			//Not supporting cpp inclusion yet.
			fmt.Printf("Ignoring cpp_include statement '%s' in document '%s'\n", line, doc.Name)
			continue // try next line

		case "namespace":
			if len(fields) != 3 {
				fmt.Println("Invalid namespace header. Expecting 'namespace <target> <name>'.")
				continue // try next line
			}
			tn := TargetName(fields[1])
			nsn := NamespaceName(fields[2])
			target, exists := doc.T.Targets[tn]
			if !exists {
				target, err = doc.T.createTarget(tn)
				if err != nil {
					panic(err)
				}
			}
			namespace, exists := target.Namespaces[nsn]
			if !exists {
				namespace, err = target.createNamespace(nsn)
				if err != nil {
					panic(err)
				}
			}

			// This document describes the namespace by header statement (hence value is true).
			namespace.Documents[doc] = true

			// Connect namespace on the document.
			doc.NamespaceByTargetName[tn] = namespace

			//++ ?? more to do ?

			continue // try next line

		default:
			// Seems we are at the end of the headers (and now at the first definition).
			doc.lastParsedLineNumber--
			return
		}
	}

	// Found end of document (all lines parsed) without finding the start of a definition.
	// No first definition found.
	fmt.Printf("It looks like document '%s' contains no definitions.\n", doc.Name)
}

func (doc *Document) parseDocumentDefinitions() {
	// loop through lines
	for {
		line, ok := doc.NextMeaningfulLine()
		if !ok {
			break // No new line. End of document.
		}

		if strings.HasPrefix(line, "const") { // 'const' FieldType Identifier '=' ConstValue ListSeparator?
			//++
		} else if strings.HasPrefix(line, "typedef") { // 'typedef' DefinitionType Identifier
			//++
		} else if strings.HasPrefix(line, "enum") { // 'enum' Identifier '{' (Identifier ('=' IntConstant)? ListSeparator?)* '}'
			//++
		} else if strings.HasPrefix(line, "senum") { // 'senum' Identifier '{' (Literal ListSeparator?)* '}'
			//++
		} else if strings.HasPrefix(line, "struct") { // 'struct' Identifier 'xsd_all'? '{' Field* '}'
			//++
		} else if strings.HasPrefix(line, "exception") { // 'exception' Identifier '{' Field* '}'
			//++
		} else if strings.HasPrefix(line, "service") { // 'service' Identifier ( 'extends' Identifier )? '{' Function* '}'		}
			//++
		}
		fmt.Printf("definition: %s\n", line)
	}
	return
}
