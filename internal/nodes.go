package internal

// Node is the most general definition of a node.
// A node has an id, unique, that is, same id implies same node
type Node interface {
	// Id returns an unique id of a node
	Id() string
}

// SameNodes returns true if the nodes have the same id (or are both nil), false otherwise
func SameNodes(a, b Node) bool {
	switch {
	case a == nil:
		return b == nil
	case b == nil:
		return a == nil
	default:
		return a.Id() == b.Id()
	}
}

// NodesIterator defines a general iterator.
// Data may come from a graph database, another storage system, in memory iterator
type NodesIterator[N Node] GeneralIterator[N]
