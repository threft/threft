package tidm

// Reference references a Document and Identifier within the TIDM
type Reference struct {
	DocumentName   DocumentName
	IdentifierName IdentifierName
}

// ConstReference references a Const within the TIDM
type ConstReference Reference

// TypedefReference references a Typedef within the TIDM
type TypedefReference Reference

// EnumReference references an Enum within the TIDM
type EnumReference Reference

// SenumReference references a Senum within the TIDM
type SenumReference Reference

// StructReference references a Struct within the TIDM
type StructReference Reference

// ExceptionReference references an Exception within the TIDM
type ExceptionReference Reference

// ServiceReference references a Service within the TIDM
type ServiceReference Reference
