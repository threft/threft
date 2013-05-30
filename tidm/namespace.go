package tidm

import (
	"fmt"
)

// Name for a Namespace.
type NamespaceName string

// Namespace defines a set of definitions with unique identifiers within a single scope.
type Namespace struct {
	// target for this namespace
	target *Target //++ TODO: remove this? is this ever used?

	Name NamespaceName // the name of this namespace

	// References to definitions this namespace contains
	ConstReferences     map[IdentifierName]*ConstReference
	TypedefReferences   map[IdentifierName]*TypedefReference
	EnumReferences      map[IdentifierName]*EnumReference
	StructReferences    map[IdentifierName]*StructReference
	ExceptionReferences map[IdentifierName]*ExceptionReference
	ServiceReferences   map[IdentifierName]*ServiceReference

	// source & parse management
	identifiers map[IdentifierName]*Identifier // list of identifiers used in this namespace, used to check uniqueness
}

// create a new (empty) namespace for the target
func (target *Target) newNamespace(name NamespaceName) (*Namespace, error) {
	// check for existing namespace
	_, exists := target.Namespaces[name]
	if exists {
		return nil, fmt.Errorf("Namespace '%s' exists already on target '%s'", name, target.Name)
	}

	// create and save new namespace
	newNamespace := &Namespace{
		target:      target,
		identifiers: make(map[IdentifierName]*Identifier),

		Name: name,

		ConstReferences:     make(map[IdentifierName]*ConstReference),
		TypedefReferences:   make(map[IdentifierName]*TypedefReference),
		EnumReferences:      make(map[IdentifierName]*EnumReference),
		StructReferences:    make(map[IdentifierName]*StructReference),
		ExceptionReferences: make(map[IdentifierName]*ExceptionReference),
		ServiceReferences:   make(map[IdentifierName]*ServiceReference),
	}
	target.Namespaces[name] = newNamespace

	// all done
	return newNamespace, nil
}

// Human readable identifier for this namespace (includes target)
func (ns *Namespace) FullName() string {
	return fmt.Sprintf("[target: %s, namespace: %s]", ns.target.Name, ns.Name)
}
