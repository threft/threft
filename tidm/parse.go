package tidm

type ParseErrorType int

const (
	ParseErrorTypeUnexpectedKeyword = ParseErrorType(iota)
	ParseErrorTypeAlreadyParsed
	ParseErrorTypeNoDefinitionsFound
	ParseErrorTypeUnexpectedError
)

// ParseError contains information about a parse error
// ParseError implements the go-builtin error interface
type ParseError struct {
	Type       ParseErrorType
	SourceName string
	SourceLine int
	Message    string
}

// Error method to implement the go-builtin error interface
func (pe *ParseError) Error() string {
	return pe.Message
}
