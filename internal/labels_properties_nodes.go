package internal

import (
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// LabelsPropertiesNode is a node with labels (as a list) and properties (as a map).
// Neo4j graph database uses this model
type LabelsPropertiesNode struct {
	// nodeId is the id of the node
	nodeId string
	// nodeLabels is a set of labels
	nodeLabels map[string]bool
	// nodeProperties as a map. Note that Neo4j uses more types (list, date, etc)
	nodeProperties map[string]string
}

// Neo4jNode is just an alias to ease use of Neo4j
type Neo4jNode LabelsPropertiesNode

// Id returns the id of the node, default value for nil node
func (lpn *LabelsPropertiesNode) Id() string {
	var result string
	if lpn != nil {
		result = lpn.nodeId
	}
	return result
}

// SameNodes returns true if the nodes have the same id (or are both nil), false otherwise
func (lpn *LabelsPropertiesNode) SameNode(other graphs.Node) bool {
	switch {
	case other == nil:
		return lpn == nil
	case lpn == nil:
		return other == nil
	default:
		otherLPN, ok := other.(*LabelsPropertiesNode)
		return ok && lpn.Id() == otherLPN.Id()
	}
}

// NewLabelsPropertiesNode returns a new initialized node
func NewLabelsPropertiesNode() LabelsPropertiesNode {
	return LabelsPropertiesNode{
		nodeId:         graphs.NewUniqueId(),
		nodeLabels:     make(map[string]bool),
		nodeProperties: make(map[string]string),
	}
}

// AddLabel appends a label to the set of labels (no duplicate)
func (lpn *LabelsPropertiesNode) AddLabel(label string) {
	lpn.nodeLabels[label] = true
}

// Labels returns the labels as a sorted slice and nil for nil
// Sorting allows to compare labels using the order.
func (lpn *LabelsPropertiesNode) Labels() []string {
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
func (lpn *LabelsPropertiesNode) RemoveLabel(key string) {
	if lpn != nil {
		delete(lpn.nodeLabels, key)
	}
}

// Properties returns a copy of the properties of the node, nil for nil
func (lpn *LabelsPropertiesNode) Properties() map[string]string {
	if lpn == nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range lpn.nodeProperties {
		result[k] = v
	}

	return result
}

func (lpn *LabelsPropertiesNode) PropertyKeys() []string {
	if lpn == nil {
		return nil
	}

	result := make([]string, 0)
	for k := range lpn.nodeProperties {
		result = append(result, k)
	}

	return result
}

// SetProperty adds (or changes) this property of a node
func (lpn *LabelsPropertiesNode) SetProperty(key, value string) {
	if lpn != nil {
		lpn.nodeProperties[key] = value
	}
}

// GetProperty returns the value if any, true if node has the value and false otherwise
func (lpn *LabelsPropertiesNode) GetProperty(key string) (string, bool) {
	if lpn == nil {
		return "", false
	}

	value, has := lpn.nodeProperties[key]
	return value, has
}

// RemoveProperty removes a property per key
func (lpn *LabelsPropertiesNode) RemoveProperty(key string) {
	if lpn != nil {
		delete(lpn.nodeProperties, key)
	}
}
