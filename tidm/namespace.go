package tidm

import (
	"fmt"
)

// Name for a namespace.
type NamespaceName string

// For each target+namespace combination, a namespace is created. Within each namespace, all identifiers must be unique.
// There is a default namespace created for each document.
type Namespace struct {
	T         *TIDM
	Target    *Target            // The target this namespace belongs to
	Name      NamespaceName      // The name of this namespace
	Documents map[*Document]bool // Boolean value indicates: true=namespace by header. false=namespace by document name.

	Constants  map[*Identifier]*Constant
	Typedefs   map[*Identifier]*Typedef
	Enums      map[*Identifier]*Enums
	Senums     map[*Identifier]*Senum
	Structs    map[*Identifier]*Struct
	Exceptions map[*Identifier]*Exception
	Services   map[*Identifier]*Service
}

func (target *Target) createNamespace(namespaceName NamespaceName) (*Namespace, error) {
	_, exists := target.Namespaces[namespaceName]
	if exists {
		return nil, fmt.Errorf("Namespace '%s' exists already on target '%s'", target.Name, namespaceName)
	}

	// No existing namespace was found, creating a new one.
	newNamespace := &Namespace{
		T:         target.T,
		Target:    target,
		Name:      namespaceName,
		Documents: make(map[*Document]bool),

		// Definitions
		Constants:  make(map[*Identifier]*Constant),
		Typedefs:   make(map[*Identifier]*Typedef),
		Enums:      make(map[*Identifier]*Enums),
		Senums:     make(map[*Identifier]*Senum),
		Structs:    make(map[*Identifier]*Struct),
		Exceptions: make(map[*Identifier]*Exception),
		Services:   make(map[*Identifier]*Service),
	}
	target.Namespaces[namespaceName] = newNamespace
	return newNamespace, nil
}

// Human readable identifier for this namespace (target + name)
func (ns *Namespace) String() string {
	return string(ns.Target.Name) + "[" + string(ns.Name) + "]"
}
