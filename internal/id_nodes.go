package internal

import "github.com/zefrenchwan/nodz.git/graphs"

// IdNode is just a node with an id
type IdNode struct {
	// nodeId is the unique id of the node
	nodeId string
}

// NewIdNode returns a new node with the given id
func NewIdNode(id string) IdNode {
	return IdNode{
		nodeId: id,
	}
}

// NewRandomIdNode returns a new node with a random id
func NewRandomIdNode() IdNode {
	return NewIdNode(graphs.NewUniqueId())
}

// Id returns the id of the node
func (in IdNode) Id() string {
	return in.nodeId
}

// SameNode returns true if nodes are the same, based on id.
// It is not necessary that other is an instance of IdNode,
// it just needs an id somehow because it is already a node.
func (in IdNode) SameNode(other graphs.Node) bool {
	if nid, ok := other.(graphs.WithId); !ok {
		return false
	} else {
		return nid.Id() == in.Id()
	}
}
