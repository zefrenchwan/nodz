package patterns

import (
	"slices"
	"strings"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// FormalClass is a name and set of attributes.
// Its typical use is to represent a "real life" entity.
type FormalClass struct {
	// name of the class, should be unique (name and case insensitive name)
	name string
	// attributes of the class, case insensitive based
	attributes []string
}

// NewFormalClass returns a new formal class by name
func NewFormalClass(className string) FormalClass {
	return FormalClass{
		name:       className,
		attributes: make([]string, 0),
	}
}

// SameNode to implement Node and then make inheritance trees
func (c *FormalClass) SameNode(n graphs.Node) bool {
	if c == nil {
		return false
	} else if node, ok := n.(*FormalClass); !ok {
		return false
	} else {
		return node.SameFormalClass(*c)
	}
}

// SameFormalClass returns true if both classes have the same name, case insensitive
func (c *FormalClass) SameFormalClass(other FormalClass) bool {
	if c == nil {
		return false
	}

	return strings.EqualFold(other.name, c.name)
}

// AddAttribute adds an attribute in the definition of the formal class
func (c *FormalClass) AddAttribute(attribute string) {
	c.RemoveAttribute(attribute)
	c.attributes = append(c.attributes, attribute)
}

// RemoveAttribute removes an attribute in the definition of the formal class
func (c *FormalClass) RemoveAttribute(attribute string) {
	var matchingAttr string
	for _, attr := range c.attributes {
		if strings.EqualFold(attr, attribute) {
			matchingAttr = attr
		}
	}

	if matchingAttr != "" {
		localMatch := func(value string) bool { return value == matchingAttr }
		c.attributes = slices.DeleteFunc(c.attributes, localMatch)
	}
}

// ListAttributes returns all the attributes of the formal class, sorted
func (c *FormalClass) ListAttributes() []string {
	result := make([]string, len(c.attributes))
	copy(result, c.attributes)
	slices.Sort(result)
	return result
}
