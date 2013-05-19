package tidm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotParsedYet = errors.New("Cannot get a Target from an unparsed TIDM.")
)

// The TIDM is the top-level object for Threft Interface Definition Model.
// It contains documents and targets.
type TIDM struct {
	// exported fields, to be marshalled to tidm-json.
	Documents map[DocumentName]*Document // List of all documents that belong to the full TIDM. Bool indicates document parse state
	Targets   map[TargetName]*Target     // List of all targets that belong to the full TIDM. Value contains the namespaces for the target.

	// stats for info and pretty printing
	documentNameMaxLength int // Longest name, for pretty printing

	// private stuff, must be populated
	parsed bool // true when TIDM was parsed
}

// newTIDM sets up a new and empty TIDM
func newTIDM() *TIDM {
	return &TIDM{
		Documents: make(map[DocumentName]*Document),
		Targets:   make(map[TargetName]*Target),
	}
}

// Creates a new and emtpy TIDM object
func NewTIDM() *TIDM {
	return newTIDM()
}

// Document returns a *Document for given Reference, or an error when Document cannot be found.
func (t *TIDM) Document(ref Reference) (doc *Document, err error) {
	var exists bool
	doc, exists = t.Documents[ref.DocumentName]
	if !exists {
		return nil, errors.New("Document for given Reference does not exist.")
	}
	return doc, nil
}

// Const returns a *Const for given ConstReference, or an error when Const cannot be found.
func (t *TIDM) Const(ref ConstReference) (con *Const, err error) {
	var doc *Document
	doc, err = t.Document(Reference(ref))
	if err != nil {
		return nil, err
	}

	var exists bool
	con, exists = doc.Consts[ref.IdentifierName]
	if !exists {
		return nil, errors.New("Const for given ConstReference does not exist.")
	}
	return con, nil
}

// write tidm-json to given writer
func (t *TIDM) WriteTo(w io.Writer) (err error) {
	enc := json.NewEncoder(w)
	err = enc.Encode(t)
	return err
}

// read tidm-json from given reader
func ReadFrom(r io.Reader) (t *TIDM, err error) {
	t = newTIDM()
	dec := json.NewDecoder(r)
	err = dec.Decode(t)
	if err != nil {
		return nil, err
	}
	t.parsed = true
	return t, nil
}

// AddDocument adds a document to the TIDM docTree
// The given reader can be closed directly after this call returns
func (t *TIDM) AddDocument(name DocumentName, reader io.Reader) error {
	if t.parsed {
		return &ParseError{
			Type:    ParseErrorTypeAlreadyParsed,
			Message: "Cannot add a document after the TIDM has been parsed.",
		}
	}

	_, err := t.newDocumentFromReader(name, reader)
	return err
}

// Parse parses and verifies the complete TIDM tree (each document, each target, each namespace)
func (t *TIDM) Parse() (perr *ParseError) {
	if t.parsed {
		return &ParseError{
			Type:    ParseErrorTypeAlreadyParsed,
			Message: "Cannot parse an already parsed TIDM structure.",
		}
	}
	t.parsed = true

	// parse all documents
	for _, doc := range t.Documents {
		// parse headers
		perr = doc.parseDocumentHeaders()
		if perr != nil {
			return perr
		}

		// parse definitions
		perr = doc.parseDocumentDefinitions()
		if perr != nil {
			return perr
		}

		// add defined Targets to TIDM Targets map
		for targetName, _ := range doc.NamespaceForTarget {
			if _, exists := t.Targets[targetName]; !exists {
				target := newTarget(targetName)
				t.Targets[targetName] = target
			}
		}
	}

	// loop through targets and populate them with the parsed data
	for targetName, _ := range t.Targets {
		perr = t.populateTarget(targetName)
		if perr != nil {
			return perr
		}
	}

	// all done
	return
}

// Target() returns a Target for given TargetName
// If given TargetName does not exist, the default Target is returned.
func (t *TIDM) Target(targetName TargetName) (target *Target, err error) {
	if !t.parsed {
		return nil, ErrNotParsedYet
	}

	// get target from targets map
	var exists bool
	target, exists = t.Targets[targetName]

	// get default target if there is no target for given TargetName
	if !exists {
		target = t.Targets[TargetNameDefault]
	}

	// all done
	return target, nil
}

// add definitions from TIDM.Documents to the right namespace for given target
func (t *TIDM) populateTarget(targetName TargetName) (perr *ParseError) {
	var err error

	// get target from TIDM.Targets list
	target := t.Targets[targetName]

	// loop through documents
	for _, doc := range t.Documents {

		// find namespace for this target/document, create one if it does not exist
		namespaceName := doc.NamespaceForTarget[targetName]
		if len(namespaceName) == 0 {
			namespaceName = doc.NamespaceForTarget[TargetNameDefault]
		}
		namespace, nsExists := target.Namespaces[namespaceName]
		if !nsExists {
			namespace, err = target.newNamespace(namespaceName)
			if err != nil {
				return &ParseError{
					Type:    ParseErrorTypeUnexpectedError,
					Message: err.Error(),
				}
			}
		}

		// check if identifiers from this doc can 'fit' in target namespace
		for _, newIdentifier := range doc.identifiers {
			if existingIdentifier, exists := namespace.identifiers[newIdentifier.Name]; exists {
				return &ParseError{
					Type:    ParseErrorTypeDuplicateIdentifier,
					Message: fmt.Sprintf("The identifier '%s' is not unique for %s. Previous declaration at %s:%d", existingIdentifier.Name, namespace.FullName(), existingIdentifier.DocLine.DocumentName, existingIdentifier.DocLine.Line),
					DocLine: newIdentifier.DocLine,
				}
			}
		}

		// add const definitions to target namespace
		for _, c := range doc.Consts {
			namespace.identifiers[c.Identifier.Name] = c.Identifier
			namespace.ConstReferences[c.Identifier.Name] = &ConstReference{doc.Name, c.Identifier.Name}
		}

		//++ add other definitions to target namespace
	}
	return nil
}
