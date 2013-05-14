package tidm

import (
	"fmt"
)

// TargetName for a namespace (language or docstyle).
type TargetName string

const Target_default = TargetName("default")

type Target struct {
	t *TIDM

	Name       TargetName
	Namespaces map[NamespaceName]*Namespace
}

func (T *TIDM) createTarget(targetName TargetName) (*Target, error) {
	_, exists := T.Targets[targetName]
	if exists {
		return nil, fmt.Errorf("There is already a target '%s' name on TIDM", targetName)
	}

	// Create Target object
	target := &Target{
		T:          T,
		Name:       targetName,
		Namespaces: make(map[NamespaceName]*Namespace),
	}

	// Store in TIDM Targets map
	T.Targets[targetName] = target

	// Return created target
	return target, nil
}
