package internal

import "github.com/zefrenchwan/nodz.git/graphs"

// NeighborsIterator is a basic implementation for a neighborhood.
// It contains main stats and decorates an iterator factory
type NeighborsIterator[N graphs.Node, L graphs.Link[N]] struct {
	// CurrentNode is the center of the neighborhood
	CurrentNode N
	// IncomingCounter is the number of incoming links from the current node's point of view
	IncomingCounter int64
	// OutgoingCounter is the number of outgoing links from the current node's point of view
	OutgoingCounter int64
	// UndirectedCounter is the number of undirected links from the current node's point of view
	UndirectedCounter int64
	// IteratorsFactory returns an iterator over the links from that node.
	// We may use stats and not link iteration, so we do not embed a slice so far.
	// Build it as late as you can to avoid dealing with memory issues.
	IteratorsFactory func() graphs.LinksIterator[N, L]
}

// CenterNode is the node that we look neighbors of
func (it NeighborsIterator[N, L]) CenterNode() N {
	return it.CurrentNode
}

// IncomingDegree is the number of incoming links to current node
func (it NeighborsIterator[N, L]) IncomingDegree() int64 {
	return it.IncomingCounter
}

// OutgoingDegree is the number of outgoing links to current node
func (it NeighborsIterator[N, L]) OutgoingDegree() int64 {
	return it.OutgoingCounter
}

// UndirectedDegree is the number of undirected links to current node
func (it NeighborsIterator[N, L]) UndirectedDegree() int64 {
	return it.UndirectedCounter
}

// Links returns a lazy loaded iterator over the links of the node
func (it NeighborsIterator[N, L]) Links() (graphs.LinksIterator[N, L], error) {
	return it.IteratorsFactory(), nil
}
