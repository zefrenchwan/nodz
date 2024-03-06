package internal

import "github.com/zefrenchwan/nodz.git/graphs"

// PropertiesNode is a node with an id and properties in a map
type PropertiesNode struct {
	// nodeId is the unique id of the node
	nodeId string
	// nodeProperties are a map of key value stored as strings
	nodeProperties map[string]string
}

// NewPropertiesNode returns an empty properies node, with an id
func NewPropertiesNode() PropertiesNode {
	return PropertiesNode{
		nodeId:         graphs.NewUniqueId(),
		nodeProperties: make(map[string]string),
	}
}

// SameNode returns true based on the id of the nodes
func (n *PropertiesNode) SameNode(other graphs.Node) bool {
	if n == nil || other == nil {
		return n == nil && other == nil
	}

	if otherNode, ok := other.(*PropertiesNode); ok {
		return otherNode.nodeId == n.nodeId
	} else {
		return false
	}
}

// Id returns the unique id of the node, empty for nil
func (n *PropertiesNode) Id() string {
	var result string
	if n != nil {
		result = n.nodeId
	}

	return result
}

// SetProperty forces a value linked to a key, no matter the previous one if any
func (n *PropertiesNode) SetProperty(key, value string) {
	if n != nil {
		n.nodeProperties[key] = value
	}
}

// GetProperty returns the value, if any, linked to the key.
// If there is a value, second result is true, false otherwise
func (n *PropertiesNode) GetProperty(key string) (string, bool) {
	var result string
	if n != nil && n.nodeProperties != nil {
		result, ok := n.nodeProperties[key]
		return result, ok
	}

	return result, false
}

// PropertyKeys returns the slice of all keys, in no particular order
func (n *PropertiesNode) PropertyKeys() []string {
	if n == nil {
		return nil
	}

	result := make([]string, 0)
	if n.nodeProperties != nil {
		for key := range n.nodeProperties {
			result = append(result, key)
		}
	}

	return result
}
