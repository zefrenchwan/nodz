package internal

import "slices"

// Node is the most general definition of a node.
// A node has an id, unique, that is, same id implies same node
type Node interface {
	// SameNode tests if another node is "the same as" this one.
	// It generally means same implementation and same value (value based) or id (id based)
	SameNode(other Node) bool
}

// NodesIterator defines a general iterator.
// Data may come from a graph database, another storage system, in memory iterator
type NodesIterator[N Node] GeneralIterator[N]

// LabelPropertiesNode is a node with labels (as a list) and properties (as a map).
// Neo4j graph database uses this model
type LabelPropertiesNode struct {
	// nodeId is the id of the node
	nodeId string
	// nodeLabels is a set of labels
	nodeLabels map[string]bool
	// nodeProperties as a map. Note that Neo4j uses more types (list, date, etc)
	nodeProperties map[string]string
}

// Id returns the id of the node, default value for nil node
func (lpn *LabelPropertiesNode) Id() string {
	var result string
	if lpn != nil {
		result = lpn.nodeId
	}
	return result
}

// SameNodes returns true if the nodes have the same id (or are both nil), false otherwise
func (lpn *LabelPropertiesNode) SameNode(other Node) bool {
	switch {
	case other == nil:
		return lpn == nil
	case lpn == nil:
		return other == nil
	default:
		otherLPN, ok := other.(*LabelPropertiesNode)
		return ok && lpn.Id() == otherLPN.Id()
	}
}

// NewLabelPropertiesNode returns a new initialized node
func NewLabelPropertiesNode() LabelPropertiesNode {
	return LabelPropertiesNode{
		nodeId:         NewUniqueId(),
		nodeLabels:     make(map[string]bool),
		nodeProperties: make(map[string]string),
	}
}

// AddLabel appends a label to the set of labels (no duplicate)
func (lpn *LabelPropertiesNode) AddLabel(label string) {
	lpn.nodeLabels[label] = true
}

// Labels returns the labels as a sorted slice and nil for nil
// Sorting allows to compare labels using the order.
func (lpn *LabelPropertiesNode) Labels() []string {
	if lpn == nil {
		return nil
	}

	result := make([]string, len(lpn.nodeLabels))
	index := 0

	for k := range lpn.nodeLabels {
		result[index] = k
		index++
	}

	slices.Sort(result)

	return result
}

// RemoveLabel removes, if any, the label key
func (lpn *LabelPropertiesNode) RemoveLabel(key string) {
	if lpn != nil {
		delete(lpn.nodeLabels, key)
	}
}

// Properties returns a copy of the properties of the node, nil for nil
func (lpn *LabelPropertiesNode) Properties() map[string]string {
	if lpn == nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range lpn.nodeProperties {
		result[k] = v
	}

	return result
}

// SetProperty adds (or changes) this property of a node
func (lpn *LabelPropertiesNode) SetProperty(key, value string) {
	if lpn != nil {
		lpn.nodeProperties[key] = value
	}
}

// GetProperty returns the value if any, true if node has the value and false otherwise
func (lpn *LabelPropertiesNode) GetProperty(key string) (string, bool) {
	if lpn == nil {
		return "", false
	}

	value, has := lpn.nodeProperties[key]
	return value, has
}

// RemoveProperty removes a property per key
func (lpn *LabelPropertiesNode) RemoveProperty(key string) {
	if lpn != nil {
		delete(lpn.nodeProperties, key)
	}
}
