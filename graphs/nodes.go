package graphs

// Node is the most general definition of a node.
// A node does not have, in general, an id, it may be simpler.
// A node is NOT comparable out of the box (so no map...)
type Node interface {
	// SameNode tests if another node is "the same as" this one.
	// It generally means same implementation and same value (value based) or id (id based)
	SameNode(other Node) bool
}

// NodesIterator defines a general iterator.
// Data may come from a graph database, another storage system, in memory iterator
type NodesIterator[N Node] GeneralIterator[N]
