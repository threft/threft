package tidm

import (
	"fmt"
	"regexp"
	"strings"
)

type ParseErrorType int

const (
	ParseErrorTypeUnexpectedKeyword = ParseErrorType(iota)
	ParseErrorTypeAlreadyParsed
	ParseErrorTypeNoDefinitionsFound
	ParseErrorTypeInvalidConstDefinition
	ParseErrorTypeDuplicateIdentifier
	ParseErrorTypeUnexpectedError
	ParseErrorNotSupported
)

// ParseError contains information about a parse error
// ParseError implements the go-builtin error interface
type ParseError struct {
	Type    ParseErrorType // Type of error
	Message string         // Error message
	DocLine *DocLine       // DocLine where the problem has ocurred
}

// Error method to implement the go-builtin error interface
func (pe *ParseError) Error() string {
	return pe.Message
}

var (
	regexpMatchIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_\-\.]*$`)

	// TODO(GeertJohan): Add escape functionality to double-quoted string
	// TODO(GeertJohan): remove support for string literal with single quote.
	regexpMatchStringLiteral = regexp.MustCompile(`^(?:"[^"]*") | (?:'[^']*')$`)
)

// nextMeaningfulLine gives the next line that is not empty nor a comment
func (doc *Document) nextMeaningfulLine() string {
	for {
		// check if the complete doc has been parsed
		if len(doc.lines)-1 == doc.lastParsedLineNumber {
			return ""
		}

		// fetch next line
		doc.lastParsedLineNumber++
		line := doc.lines[doc.lastParsedLineNumber]

		// fast path for empty line
		if len(line) == 0 {
			continue
		}

		// remove comments from line
		pos = strings.Index(line, "//")
		if pos > -1 {
			line = line[:pos]
		}

		// trim space and list seperators from line
		line = strings.TrimSpace(line)
		// disabled now.. .threft file just shouldn't contain these character at the end of the line..
		// line = strings.TrimRight(line, ",; ")
		// line = strings.TrimSpace(line)

		// try next line if this one is empty after removing
		if len(line) == 0 {
			continue
		}
		return line
	}
}

// parseDocumentHeaders parses document headers
func (doc *Document) parseDocumentHeaders() *ParseError {
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
func (doc *Document) parseDocumentDefinitions() *ParseError {
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
			return &ParseError{
				Type:    ParseErrorNotSupported,
				Message: "Error: typedef is not supported right now.",
				DocLine: currentDocLine,
			}
		case "struct": // 'struct' Identifier 'xsd_all'? '{' Field* '}'
			//++
		case "exception": // 'exception' Identifier '{' Field* '}'
			//++
		case "service": // 'service' Identifier ( 'extends' Identifier )? '{' Function* '}'		}
			//++
		default:
			return &ParseError{
				Type:    ParseErrorTypeUnexpectedKeyword,
				Message: fmt.Sprintf("Error: keyword '%s' is not valid.", words[0]),
				DocLine: currentDocLine,
			}
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
