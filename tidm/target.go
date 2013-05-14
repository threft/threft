package tidm

import ()

// TargetName for a namespace (language or docstyle).
type TargetName string

const (
	TargetNameDefault = TargetName("*")
	TargetNameHtml    = TargetName("html")
	TargetNameGo      = TargetName("go")
)

// Target defines a set of namespaces within a target.
// A target can be a language
type Target struct {
	Name       TargetName
	Namespaces map[NamespaceName]*Namespace
}

func newTarget(name TargetName) *Target {
	// create new Target
	target := &Target{
		Name:       name,
		Namespaces: make(map[NamespaceName]*Namespace),
	}

	// Return created target
	return target
}
