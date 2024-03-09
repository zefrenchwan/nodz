package internal

import "github.com/zefrenchwan/nodz.git/graphs"

// UndirectedSimpleLink is the simplest undirected link implementation
type UndirectedSimpleLink[N graphs.Node] struct {
	// linkSource1 is the source of the node
	linkSource1 N
	// linkSource2 is the destination of the node
	linkSource2 N
}

// NewUndirectedSimpleLink returns a link which extremities are source and destination.
func NewUndirectedSimpleLink[N graphs.Node](source, destination N) UndirectedSimpleLink[N] {
	var result UndirectedSimpleLink[N]
	result.linkSource1 = source
	result.linkSource2 = destination
	return result
}

// SameLink returns true if undirected links are the same :
// Same class, same sources and destinations or destination of one is the source of the other
func (l UndirectedSimpleLink[N]) SameLink(other graphs.Link[N]) bool {
	if other == nil {
		return false
	}

	if otherLink, ok := other.(UndirectedSimpleLink[N]); !ok {
		return false
	} else if l.linkSource1.SameNode(otherLink.linkSource1) && l.linkSource2.SameNode(otherLink.linkSource2) {
		return true
	} else {
		return l.linkSource1.SameNode(otherLink.linkSource2) && l.linkSource2.SameNode(otherLink.linkSource1)
	}
}

// Source returns the source of the link
func (l UndirectedSimpleLink[N]) Source() N {
	return l.linkSource1
}

// Destination returns the destination of the link
func (l UndirectedSimpleLink[N]) Destination() N {
	return l.linkSource2
}

// IsDirected returns false by definition
func (l UndirectedSimpleLink[N]) IsDirected() bool {
	return false
}
