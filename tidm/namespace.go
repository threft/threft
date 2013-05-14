package tidm

import (
	"fmt"
)

// Name for a namespace.
type NamespaceName string

// For each target+namespace combination, a namespace is created. Within each namespace, all identifiers must be unique.
type Namespace struct {
	target *Target

	Name        NamespaceName // the name of this namespace
	Definitions *Definitions  // definitions within this namespace

	identifierStrings map[string]bool // Identifiers used in this namespace
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

		identifierStrings: make(map[string]bool),
	}
	target.Namespaces[name] = newNamespace

	// all done
	return newNamespace, nil
}

// Human readable identifier for this namespace
func (ns *Namespace) String() string {
	return string(ns.target.Name) + "[" + string(ns.Name) + "]"
}
