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

// DocumentName represents any documentname, this can be anything (filename, random string, "stdin", etc.)
type DocumentName string

// DocLine is a reference to a thrift idl document and the source line number
type DocLine struct {
	DocumentName DocumentName
	Line         int
}

func (dl DocLine) String() string {
	return fmt.Sprintf("%s:%d", dl.DocumentName, dl.Line+1)
}

// Document can be any thrift definition language source (filename, StdIn, etc.)
// A document is created by this package. See TIDM.AddDocument()
type Document struct {
	t *TIDM

	Name               DocumentName                 // The name of this document.
	NamespaceForTarget map[TargetName]NamespaceName // List of namespaces (per target) this document describes

	// Definitions in this document
	Consts     map[IdentifierName]*Const
	Typedefs   map[IdentifierName]*Typedef
	Enums      map[IdentifierName]*Enums
	Senums     map[IdentifierName]*Senum
	Structs    map[IdentifierName]*Struct
	Exceptions map[IdentifierName]*Exception
	Services   map[IdentifierName]*Service

	// source & parse management
	identifiers          map[IdentifierName]*Identifier // list of identifiers used in this document, used to check uniqueness
	lines                []string                       // All source lines for this document.
	lastParsedLineNumber int                            // line number of the last parsed line. Used by nextMeaningfulLine().
}

func (t *TIDM) newDocument(name DocumentName) (*Document, error) {
	// check if docname is unique
	if _, exists := t.Documents[name]; exists {
		return nil, ErrDocumentWithNameExists
	}

	// create and save new doc
	doc := &Document{
		t: t,

		Name:               name,
		NamespaceForTarget: make(map[TargetName]NamespaceName),

		Consts:     make(map[IdentifierName]*Const),
		Typedefs:   make(map[IdentifierName]*Typedef),
		Enums:      make(map[IdentifierName]*Enums),
		Senums:     make(map[IdentifierName]*Senum),
		Structs:    make(map[IdentifierName]*Struct),
		Exceptions: make(map[IdentifierName]*Exception),
		Services:   make(map[IdentifierName]*Service),

		identifiers:          make(map[IdentifierName]*Identifier),
		lastParsedLineNumber: -1,
	}
	doc.NamespaceForTarget[TargetNameDefault] = NamespaceName(strings.Replace(string(name), ".thrift", "", -1))
	t.Documents[name] = doc

	// set max doc filename length
	if len(name) > t.documentNameMaxLength {
		t.documentNameMaxLength = len(name)
	}

	// all done
	return doc, nil
}

func (t *TIDM) newDocumentFromReader(name DocumentName, sourceInput io.Reader) (*Document, error) {
	// create an empty document
	doc, err := t.newDocument(name)
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
		doc.lines = append(doc.lines, line)
	}

	// all done
	return doc, nil
}

// nextMeaningfulLine gives the next line that is not empty nor a comment
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

// parseDocumentHeaders parses document headers
func (doc *Document) parseDocumentHeaders() (perr *ParseError) {
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
			fmt.Printf("Ignoring include statement at %s:%d\n", doc.Name, doc.lastParsedLineNumber+1)
			continue

		case "cpp_include":
			// not supporting cpp inclusion (yet?).
			fmt.Printf("Ignoring cpp_include statement at %s:%d\n", doc.Name, doc.lastParsedLineNumber+1)
			continue

		case "namespace":
			// invalid namespace header. notify user, then continue.
			if len(fields) != 3 {
				fmt.Println("Invalid namespace header. Expecting 'namespace <target> <name>'.")
				continue
			}

			// add target/namespace to document
			targetName := TargetName(fields[1])
			namespaceName := NamespaceName(fields[2])
			doc.NamespaceForTarget[targetName] = namespaceName

			// done, next line!
			continue

		default:
			// it seems that we arived at the end of the headers and now at the first definition
			// decrement lastParsedLineNumber, this line should be parsed by parseDocumentDefinitions()
			doc.lastParsedLineNumber--
			// done parsing headers successfully
			return nil
		}
	}

	// end of document
	return nil
}

// parseDocumentDefinitions parses document definitions
func (doc *Document) parseDocumentDefinitions() (perr *ParseError) {
	var countDefinitions int

	// loop through lines
	for {
		line := doc.nextMeaningfulLine()
		if len(line) == 0 {
			// done
			break
		}
		countDefinitions++

		currentDocLine := &DocLine{
			DocumentName: doc.Name,
			Line:         doc.lastParsedLineNumber,
		}
		words := strings.Fields(line)

		switch words[0] {
		case "const": // 'const' FieldType Identifier '=' ConstValue ListSeparator?
			if len(words) != 5 {
				return &ParseError{
					Type:    ParseErrorTypeInvalidConstDefinition,
					Message: "Invalid const definition.",
					DocLine: currentDocLine,
				}
			}

			//++ TODO: check field type

			//++ TODO: regexp identifier

			// check if identifier is unique
			if i, exists := doc.identifiers[IdentifierName(words[2])]; exists {
				return &ParseError{
					Type:    ParseErrorTypeDuplicateIdentifier,
					Message: fmt.Sprintf("The given identifier has been declared before in this document. Previous declaration at %s", i.DocLine),
					DocLine: currentDocLine,
				}
			}
			// check that third word is an equal sign
			if words[3] != "=" {
				return &ParseError{
					Type:    ParseErrorTypeInvalidConstDefinition,
					Message: "Invalid const definition. Expecting '='.",
					DocLine: currentDocLine,
				}
			}
			// create constant instance
			c := &Const{
				Type: FieldType(words[1]),
				Identifier: &Identifier{
					Name:    IdentifierName(words[2]),
					DocLine: currentDocLine,
				},
				Value: words[4],
			}
			// save identifier and constant
			doc.identifiers[c.Identifier.Name] = c.Identifier
			doc.Consts[c.Identifier.Name] = c

		case "typedef": // 'typedef' DefinitionType Identifier
			return &ParseError{
				Type:    ParseErrorNotSupported,
				Message: "Error: typedef is not supported right now.",
				DocLine: currentDocLine,
			}
		case "enum": // 'enum' Identifier '{' (Identifier ('=' IntConstant)? ListSeparator?)* '}'
			//++
		case "senum": // 'senum' Identifier '{' (Literal ListSeparator?)* '}'
			return &ParseError{
				Type:    ParseErrorNotSupported,
				Message: "Error: senum is not supported right now.",
				DocLine: currentDocLine,
			}
		case "struct": // 'struct' Identifier 'xsd_all'? '{' Field* '}'
			//++
		case "exception": // 'exception' Identifier '{' Field* '}'
			//++
		case "service": // 'service' Identifier ( 'extends' Identifier )? '{' Function* '}'		}
			//++
		}
		// fmt.Printf("definition: %s\n", line)
	}

	// it is required that the document contained definitions, otherwise return an error
	if countDefinitions == 0 {
		return &ParseError{
			Type:    ParseErrorTypeNoDefinitionsFound,
			Message: fmt.Sprintf("No definitions found for document '%s'. This is an error.", doc.Name),
		}
	}

	// all done
	return nil
}
