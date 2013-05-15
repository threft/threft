package tidm

import (
	"errors"
)

const (
	ErrIdentifierExists = errors.New("Identifiers exists already")
)

type Definitions struct {
	// list of identifiers used in this set of Definitions
	IdentifierStrings map[string]bool

	// Actual definitions
	Constants  map[*Identifier]*Constant
	Typedefs   map[*Identifier]*Typedef
	Enums      map[*Identifier]*Enums
	Senums     map[*Identifier]*Senum
	Structs    map[*Identifier]*Struct
	Exceptions map[*Identifier]*Exception
	Services   map[*Identifier]*Service
}

func newDefinitions() *Definitions {
	return &Definitions{
		IdentifierStrings: make(map[string]bool),

		Constants:  make(map[*Identifier]*Constant),
		Typedefs:   make(map[*Identifier]*Typedef),
		Enums:      make(map[*Identifier]*Enums),
		Senums:     make(map[*Identifier]*Senum),
		Structs:    make(map[*Identifier]*Struct),
		Exceptions: make(map[*Identifier]*Exception),
		Services:   make(map[*Identifier]*Service),
	}
}

//++ TODO: Add AddConstant, AddTypedef, etc. methods that check if identifier is unique and then add the type and its identifier

//++ TODO
type Constant struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Typedef struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Enums struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Senum struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Struct struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Exception struct {
	DocLine DocLine
	foo     string
	bar     int
}

//++ TODO
type Service struct {
	DocLine DocLine
	foo     string
	bar     int
}
