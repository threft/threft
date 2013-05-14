package tidm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrDocumentWithNameExists = errors.New("A document with given name exists already in this TIDM")
)

//++ TODO: public? seriously?
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

// Document can be any 
type Document struct {
	t *TIDM

	Name               DocumentName                 `json:"name"`               // The name of this document.
	IdentifierStrings  map[string]bool              `json:"identifierStrings"`  // All identifiers in this document
	Definitions        Definitions                  `json:"definitions"`        // List of definitions in this document
	NamespaceForTarget map[TargetName]NamespaceName `json:"namespaceForTarget"` // List of namespaces (per target) this document describes

	// source management
	lines                []string // All source lines for this document. Whitespace is trimmed.
	lastParsedLineNumber int      // line number of the last parsed line. Used by nextMeaningfulLine().
}

// DocLine is a reference to a source document and the code line
type DocLine struct {
	Document *Document
	Line     int
}

func (t *TIDM) newDocument(name DocumentName) (*Document, error) {
	// check if docname is unique
	if _, exists := t.Documents[name]; exists {
		return nil, ErrDocumentWithNameExists
	}

	// create and save new doc
	doc := &Document{
		t:                    t,
		IdentifierStrings:    make(map[string]bool),
		Definitions:          newDefinitions(),
		NamespaceForTarget:   make(map[TargetName]NamespaceName),
		lastParsedLineNumber: -1,
	}
	t.Documents[name] = doc

	// set max doc filename length 
	if len(docName) > t.documentNameMaxLength {
		t.documentNameMaxLength = len(docName)
	}

	// all done
	return doc, nil
}

func (t *TIDM) newDocumentFromReader(name DocumentName, sourceInput io.Reader) (*Document, error) {
	// create an empty document
	doc, err := t.newDocument()
	if err != nil {
		return nil, err
	}

	// read and store lines
	sourcereader := bufio.NewReader(sourceInput)
	for {
		//++ TODO: Blegh, this for loop feels like a trampoline. The EOF pickup could be nicer..
		//++ TODO: Also: ReadString('\n') isn't really cross-platform, is it?
		line, err := sourcereader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					goto addLine
				}
				break
			}
			fmt.Printf("Error while reading line %d from sourceInput %s. %s\n", len(doc.lines)+1, name, err)
			break
		}
	addLine:
		doc.lines = append(doc.lines, strings.TrimSpace(line))
	}

	// all done
	return doc, nil
}

func (doc *Document) nextMeaningfulLine() (line string) {
	for {
		// check if the complete doc has been parsed
		if len(doc.lines)-1 == doc.lastParsedLineNumber {
			return ""
		}

		// fetch next line
		doc.lastParsedLineNumber++
		line = doc.lines[doc.lastParsedLineNumber]

		// remove comments from line
		pos := strings.Index(line, "#")
		if pos > -1 {
			line = line[:pos]
		}
		pos = strings.Index(line, "//")
		if pos > -1 {
			line = line[:pos]
		}

		// trim space and list seperators from line
		line = strings.TrimSpace(line)
		line = strings.TrimRight(line, ",; ")
		line = strings.TrimSpace(line)

		// try next line if this one is empty
		if len(line) == 0 {
			continue
		}
		return line
	}
}

func (doc *Document) parseDocumentHeaders() {
	var err error

	// loop through lines
	for {

		line := doc.nextMeaningfulLine()
		if len(line) == 0 {
			break // no new lines
		}

		// get fields from line
		fields := strings.Fields(line)

		// switch on keyword
		switch fields[0] {
		case "include":
			// not supporting cross-document references (yet?).
			fmt.Printf("Ignoring include statement '%s' in document '%s'\n", line, doc.Name)
			continue

		case "cpp_include":
			// not supporting cpp inclusion (yet?).
			fmt.Printf("Ignoring cpp_include statement '%s' in document '%s'\n", line, doc.Name)
			continue

		case "namespace":
			// invalid namespace header. notify user, then continue.
			if len(fields) != 3 {
				fmt.Println("Invalid namespace header. Expecting 'namespace <target> <name>'.")
				continue
			}

			//++ TODO: regexp over targetName and namespaceName
			tName := TargetName(fields[1])
			nName := NamespaceName(fields[2])

			continue

		default:
			// it seems that we arived at the end of the headers and now at the first definition.
			// decrement lastParsedLineNumber, this line should be parsed by parseDocumentDefinitions()
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
		line, ok := doc.nextMeaningfulLine()
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
