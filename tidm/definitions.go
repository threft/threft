package tidm

import (
	"errors"
)

var (
	ErrDulicateIdentifier = errors.New("Duplicate identifier")
)

// Definitions is a set of definitions with a unique identifier
type Definitions struct {
	// list of identifiers used in this set of Definitions
	identifiers map[IdentifierName]*Identifier

	// Actual definitions
	Constants  map[IdentifierName]*Constant  `json:"constants"`
	Typedefs   map[IdentifierName]*Typedef   `json:"typedefs"`
	Enums      map[IdentifierName]*Enums     `json:"enums"`
	Senums     map[IdentifierName]*Senum     `json:"senums"`
	Structs    map[IdentifierName]*Struct    `json:"structs"`
	Exceptions map[IdentifierName]*Exception `json:"exceptions"`
	Services   map[IdentifierName]*Service   `json:"services"`
}

func newDefinitions() *Definitions {
	return &Definitions{
		identifiers: make(map[IdentifierName]*Identifier),

		Constants:  make(map[IdentifierName]*Constant),
		Typedefs:   make(map[IdentifierName]*Typedef),
		Enums:      make(map[IdentifierName]*Enums),
		Senums:     make(map[IdentifierName]*Senum),
		Structs:    make(map[IdentifierName]*Struct),
		Exceptions: make(map[IdentifierName]*Exception),
		Services:   make(map[IdentifierName]*Service),
	}
}

//++ TODO: Add AddConstant, AddTypedef, etc. methods that check if identifier is unique and then add the type and its identifier

//++ TODO
type Constant struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Typedef struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Enums struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Senum struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Struct struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Exception struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}

//++ TODO
type Service struct {
	Identifier *Identifier

	//++ TODO: fields have their own type (structs) with data and DocLine to identify the field-specific doc and line

	Foo string
	Bar int
}
