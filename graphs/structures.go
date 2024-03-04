package graphs

// ValueBasedGraph is the general definition of a graph that is sort of "ready for use":
// definition does not provide nodes, links, neighborhoods.
type ValueBasedGraph[NV comparable, LV comparable] interface {
	// SetLink set a value from a source to a destination, raises an error when failing.
	// SetLink changes the value if any, or create the link otherwise.
	// If nodes were not in the graph, they are created
	SetLink(source, destination NV, value LV) error
	// RemoveLink will remove the link from the source to the destination.
	// If the implementation is undirected,
	// there should be no link between destination and source too
	RemoveLink(source, destination NV) error
	// AddNode adds a node in the graph.
	// If it was not there already, it is added as an isolated node.
	// If it was there before, no change.
	AddNode(NV) error
	// RemoveNode removes the node and all links around it.
	RemoveNode(NV) error
	// Neighbors of a value is the map of neighbors with their value, or an error
	Neighbors(NV) (map[NV]LV, error)
	// LinkValue returns true and the value if any, false and default value if no link, an error if any
	LinkValue(source, destination NV) (LV, bool, error)
}
