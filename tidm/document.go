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
