package tidm

// The TIDM is the top-level object for Threft Interface Definition Model.
// It contains documents and namespaces. Each document and namespace contains the complete definition model within that target.
type TIDM struct {
	BasePath  string                     // base folder where the .thrift files (and eventual subdirectories with .thrift files) are located.
	Documents map[DocumentName]*Document // List of all documents that belong to the full TIDM. Bool indicates document parse state
	Targets   map[TargetName]*Target     // List of all targets that belong to the full TIDM. Value contains the namespaces for the target.

	// Some details for pretty printing
	DocumentNameMaxLength int
}

// Creates a new TIDM object
func newTIDM() *TIDM {
	return &TIDM{
		BasePath:  "Unset",
		Documents: make(map[DocumentName]*Document),
		Targets:   make(map[TargetName]*Target),
	}
}
