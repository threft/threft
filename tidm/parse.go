package tidm

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
