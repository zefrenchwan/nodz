package graphs

// ValueBasedGraph is the general definition of a graph that is sort of "ready for use":
// definition does not provide nodes, links, neighborhoods.
// For instance, you may use to link cities (by name, NV = string) with distances (LV=float32).
// Although it is an "easy to use" data structure, it is also a weak one:  Neighbors returns a map.
// NOT an iterator that would allow partial loading, a map with all the content in memory.
// Use it for tests or to deal with small graphs.
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

// CentralStructureGraph is a graph that allows global operations, such as nodes or links iterations.
// Its definition should allow many implementations, from a "in memory" implementation to a distributed one.
// It also should deal with many types of links ((un)directed, valued, etc) and nodes (with data in it, or just id based nodes)
// This is why this definition provides N and L, and no direct api.
// N is then any implementation of a node, and L is any implementation of a link between two nodes (as instances of N).
// For instance, consider a graph with city objects and directed valued links.
// N and L are NOT city and value, but N is a struct that implements Node, and L is a struct that deals with direction and value.
// This structure is less intuitive (depending on your intuition...) than a value based graph, but it offers way more options.
// Because it allows a distributed storage version, all functions may return an error.
type CentralStructureGraph[N Node, L Link[N]] interface {
	// AddLink adds a node in the graph, upserts its value if any, does nothing for same content
	AddLink(L) error
	// RemoveLink removes a link but keeps the nodes
	RemoveLink(L) error
	// AddNode adds a non existing node, does nothing for an existing one
	AddNode(N) error
	// RemoveNode removes a node, does nothing if the node did not appear in the graph
	RemoveNode(N) error
	// AllNodes returns an iterator over all the nodes. Each node appearts exactly once
	AllNodes() (NodesIterator[N], error)
	// Neighbors returns the neighborhood of a node
	Neighbors(N) (Neighborhood[N, L], error)
}
