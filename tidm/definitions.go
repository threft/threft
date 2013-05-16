package tidm

// Definitions is a set of definitions with a unique identifier
type Definitions struct {
	// list of identifiers used in this set of Definitions
	identifiers map[IdentifierName]*Identifier

	// Actual definitions
	Consts     map[IdentifierName]*Const
	Typedefs   map[IdentifierName]*Typedef
	Enums      map[IdentifierName]*Enums
	Senums     map[IdentifierName]*Senum
	Structs    map[IdentifierName]*Struct
	Exceptions map[IdentifierName]*Exception
	Services   map[IdentifierName]*Service
}

func newDefinitions() *Definitions {
	return &Definitions{
		identifiers: make(map[IdentifierName]*Identifier),

		Consts:     make(map[IdentifierName]*Const),
		Typedefs:   make(map[IdentifierName]*Typedef),
		Enums:      make(map[IdentifierName]*Enums),
		Senums:     make(map[IdentifierName]*Senum),
		Structs:    make(map[IdentifierName]*Struct),
		Exceptions: make(map[IdentifierName]*Exception),
		Services:   make(map[IdentifierName]*Service),
	}
}

//++ TODO: Add AddConst, AddTypedef, etc. methods that check if identifier is unique and then add the type and its identifier

type FieldType string //++ TODO

type Const struct {
	Type       FieldType
	Identifier *Identifier
	Value      interface{}
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
