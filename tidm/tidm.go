package tidm

import (
	"encoding/json"
	"errors"
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
		perr = doc.parseDocumentHeaders()
		if perr != nil {
			return perr
		}
		perr = doc.parseDocumentDefinitions()
		if perr != nil {
			return perr
		}
	}

	//++ get list of all targets

	//++ loop through targets and "default"
	//		perr = t.populateTarget(tname)
	//		if perr != nil {
	//			return perr
	//		}
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

func (t *TIDM) WriteTo(w io.Writer) (err error) {
	enc := json.NewEncoder(w)
	err = enc.Encode(t)
	return err
}

func ReadFrom(r io.Reader) (t *TIDM, err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(t)
	if err != nil {
		return nil, err
	}
	t.parsed = true
	return t, nil
}

// create a target and populate it from documents
func (t *TIDM) populateTarget(targetName TargetName) (target *Target, perr *ParseError) {
	var err error

	// create new empty target instance
	target = newTarget(targetName)
	t.Targets[targetName] = target

	// loop through documents
	for _, doc := range t.Documents {
		// find namespace, create one if it does not exist
		namespaceName := doc.NamespaceForTarget[targetName]
		namespace, nsExists := target.Namespaces[namespaceName]
		if !nsExists {
			namespace, err = target.newNamespace(namespaceName)
			if err != nil {
				return nil, &ParseError{
					Type:    ParseErrorTypeUnexpectedError,
					Message: err.Error(),
				}
			}
		}

		//++ TODO: use addExisting methods for all definitions from doc to namespace.
		_ = namespace
	}
	return
}
