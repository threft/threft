package tidm

import (
	"fmt"
)

// Name for a Namespace.
type NamespaceName string

// Namespace defines a set of definitions with unique identifiers within a single scope.
type Namespace struct {
	target *Target

	Name        NamespaceName // the name of this namespace
	Definitions *Definitions  // definitions within this namespace
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
		Name:        name,
		Definitions: newDefinitions(),
	}
	target.Namespaces[name] = newNamespace

	// all done
	return newNamespace, nil
}

// Human readable identifier for this namespace
func (ns *Namespace) String() string {
	return string(ns.target.Name) + "[" + string(ns.Name) + "]"
}
