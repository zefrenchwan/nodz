package internal

import "github.com/zefrenchwan/nodz.git/graphs"

type ValuedLink[N graphs.Node, V comparable] struct {
	nodeSource      N
	nodeDestination N
	value           V
	directed        bool
}

func NewDirectedValuedLink[N graphs.Node, V comparable](source, destination N, value V) ValuedLink[N, V] {
	return ValuedLink[N, V]{
		nodeSource:      source,
		nodeDestination: destination,
		value:           value,
		directed:        true,
	}
}

func NewUndirectedValuedLink[N graphs.Node, V comparable](source, destination N, value V) ValuedLink[N, V] {
	return ValuedLink[N, V]{
		nodeSource:      source,
		nodeDestination: destination,
		value:           value,
		directed:        false,
	}
}

func (vl ValuedLink[N, V]) SameLink(other graphs.Link[N]) bool {
	if other == nil {
		return false
	}

	otherLink, ok := other.(ValuedLink[N, V])
	if !ok {
		return false
	}

	return otherLink.directed == vl.directed &&
		vl.nodeSource.SameNode(otherLink.nodeSource) &&
		vl.nodeDestination.SameNode(otherLink.nodeDestination) &&
		vl.value == otherLink.value

}

func (vl ValuedLink[N, V]) Source() N {
	return vl.nodeSource
}

func (vl ValuedLink[N, V]) Destination() N {
	return vl.nodeDestination
}

func (vl ValuedLink[N, V]) IsDirected() bool {
	return vl.directed
}
