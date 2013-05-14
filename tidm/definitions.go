package tidm

type Definitions struct {
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
		Constants:  make(map[*Identifier]*Constant),
		Typedefs:   make(map[*Identifier]*Typedef),
		Enums:      make(map[*Identifier]*Enums),
		Senums:     make(map[*Identifier]*Senum),
		Structs:    make(map[*Identifier]*Struct),
		Exceptions: make(map[*Identifier]*Exception),
		Services:   make(map[*Identifier]*Service),
	}
}

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
