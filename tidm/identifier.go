package tidm

type IdentifierName string

type Identifier struct {
	DocLine *DocLine
	Name    IdentifierName
}
