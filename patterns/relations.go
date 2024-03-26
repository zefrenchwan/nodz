package patterns

import (
	"errors"
	"slices"
	"strings"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// FormalRelation is a formal relation between classes.
// For instance, worksFor such as "worksFor(Employee, Employer)""
type FormalRelation struct {
	// name of the relation
	name string
	// operands of the relation, in order
	operands []FormalClass
}

// NewFormalRelation returns a new relation
func NewFormalRelation(relationName string, classes []FormalClass) (FormalRelation, error) {
	var result FormalRelation
	if relationName == "" {
		return result, errors.New("invalid relation name")
	} else if len(classes) == 0 {
		return result, errors.New("invalid classes")
	}

	result.name = relationName
	result.operands = make([]FormalClass, len(classes))
	copy(result.operands, classes)
	return result, nil
}

// SameFormalRelation returns true if relations have same name (case insensitive)
// and same operands as a slice (same index, same value)
func (f FormalRelation) SameFormalRelation(other FormalRelation) bool {
	if !strings.EqualFold(other.name, f.name) {
		return false
	} else {
		fn := func(a, b FormalClass) bool { return a.SameFormalClass(b) }
		return slices.EqualFunc(other.operands, f.operands, fn)
	}
}

// SameNode is here to use formal relations in graphs
func (f FormalRelation) SameNode(n graphs.Node) bool {
	if fnode, ok := n.(FormalRelation); !ok {
		return false
	} else {
		return f.SameFormalRelation(fnode)
	}
}
