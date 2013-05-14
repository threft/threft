package tidm

import (
	"fmt"
	"io"
)

// The TIDM is the top-level object for Threft Interface Definition Model.
// It contains documents and targets.
type TIDM struct {
	// open data
	Documents map[DocumentName]*Document `json:"documents"` // List of all documents that belong to the full TIDM. Bool indicates document parse state

	// stats for info and pretty printing
	documentNameMaxLength int // Longest name, for pretty printing

	// private stuff, must be populated
	targetNames map[TargetName]bool
	targets     map[TargetName]*Target // List of all targets that belong to the full TIDM. Value contains the namespaces for the target.
}

// newTIDM sets up a new and empty TIDM
func newTIDM() *TIDM {
	return &TIDM{
		Documents: make(map[DocumentName]*Document),
		targets:   make(map[TargetName]*Target),
	}
}

// Creates a new and emtpy TIDM object
func NewTIDM() *TIDM {
	return newTIDM()
}

// NewTIDMFromJson creates a new TIDM instance and populates it from JSON
// It returns a normal error, not a parse-error because tidm-json should've been checked.
func NewTIDMFromJson(jsonBytes []byte) (*TIDM, error) {
	t := newTIDM()
	return t, nil
}

// AddDocument adds a document to the TIDM docTree
func (t *TIDM) AddDocument(name DocumentName, reader io.Reader) error {
	_, err := t.newDocumentFromReader(name, reader)
	return err
}

// Verify verifies the complete TIDM tree (each target, each namespace)
func (t *TIDM) Verify() (perr *ParseError) {
	//++ get list of all targets

	//++ loop through targets
	//		perr = t.populateTarget(tname)
	//		if perr != nil {
	//			return perr
	//		}
	return
}

// Target() returns a Target for given TargetName
func (t *TIDM) Target(targetName TargetName) (target *Target, err error) {
	var exists bool
	var perr *ParseError

	// get target from targets map
	target, exists = t.targets[targetName]

	// see if target exists, if not, populate it
	if !exists {
		perr = t.populateTarget(targetName)
		if perr != nil {
			return nil, fmt.Errorf("Unexpected parse error, TIDM should've been verified before using Target(). %s", perr.Error())
		}
	}

	// all done
	return target, nil
}

func (t *TIDM) populateTarget(tname TargetName) (perr *ParseError) {
	//++ create target on TIDM

	//++ loop trough documents
	//++	see if namespace for this target exists in *Target, if not create it.
	//++	add items from document to namespace, check for each item if it exists
	return
}
