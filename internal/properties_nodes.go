package internal

import "github.com/zefrenchwan/nodz.git/graphs"

type PropertiesNode struct {
	nodeId         string
	nodeProperties map[string]string
}

func NewPropertiesNode() PropertiesNode {
	return PropertiesNode{
		nodeId:         graphs.NewUniqueId(),
		nodeProperties: make(map[string]string),
	}
}

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

func (n *PropertiesNode) Id() string {
	var result string
	if n != nil {
		result = n.nodeId
	}

	return result
}

func (n *PropertiesNode) SetProperty(key, value string) {
	if n != nil {
		n.nodeProperties[key] = value
	}
}

func (n *PropertiesNode) GetProperty(key string) (string, bool) {
	var result string
	if n != nil && n.nodeProperties != nil {
		result, ok := n.nodeProperties[key]
		return result, ok
	}

	return result, false
}

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
