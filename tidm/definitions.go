package tidm

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
