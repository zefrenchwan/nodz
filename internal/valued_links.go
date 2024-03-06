package internal

import "github.com/zefrenchwan/nodz.git/graphs"

// ValuedLink is a link that carries a value.
// It may be directed or not
type ValuedLink[N graphs.Node, V comparable] struct {
	// nodeSource is the source of the link, or an extremity for undirected links
	nodeSource N
	// nodeDestination is the destination of the link, or an extremity for undirected links
	nodeDestination N
	// value is the value of the link. May be an int, a float, a string, etc
	value V
	// directed is true for directed links, false for undirected
	directed bool
}

// NewDirectedValuedLink returns a valued link from a source to a destination
func NewDirectedValuedLink[N graphs.Node, V comparable](source, destination N, value V) ValuedLink[N, V] {
	return ValuedLink[N, V]{
		nodeSource:      source,
		nodeDestination: destination,
		value:           value,
		directed:        true,
	}
}

// NewUnirectedValuedLink returns a valued link with source and destination as extremities
func NewUndirectedValuedLink[N graphs.Node, V comparable](source, destination N, value V) ValuedLink[N, V] {
	return ValuedLink[N, V]{
		nodeSource:      source,
		nodeDestination: destination,
		value:           value,
		directed:        false,
	}
}

// SameLink returns true for a same direction, same source and destination, same value, false otherwise.
// For directed nodes, it means sources should be equal, and destinations should be equal.
// For undirected nodes, it means source of one is either source or destination of the other, and vice versa.
func (vl ValuedLink[N, V]) SameLink(other graphs.Link[N]) bool {
	if other == nil {
		return false
	}

	otherLink, ok := other.(ValuedLink[N, V])
	if !ok {
		return false
	}

	if otherLink.directed != vl.directed || vl.value != otherLink.value {
		return false
	}

	switch {
	case vl.directed:
		return vl.nodeSource.SameNode(otherLink.nodeSource) && vl.nodeDestination.SameNode(otherLink.nodeDestination)
	case vl.nodeSource.SameNode(otherLink.nodeDestination):
		return vl.nodeDestination.SameNode(otherLink.nodeSource)
	case vl.nodeSource.SameNode(otherLink.nodeSource):
		return vl.nodeDestination.SameNode(otherLink.nodeDestination)
	default:
		return false
	}
}

// Source returns the source of the link as defined when link was built.
// For undirected links, it matters because source may be your destination
func (vl ValuedLink[N, V]) Source() N {
	return vl.nodeSource
}

// Destination resturns the destination of the link as defined when built.
// For undirected links, it matters because destination may be your source
func (vl ValuedLink[N, V]) Destination() N {
	return vl.nodeDestination
}

// IsDirected returns true for directed links, false otherwise
func (vl ValuedLink[N, V]) IsDirected() bool {
	return vl.directed
}
