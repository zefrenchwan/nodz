package graphs

// Neighborhood of a node is the links from that node.
// Still, no need to load all the links to have basic information about the neighborhood.
// So, to include all options, three functions appear for degrees : incoming, outgoing, and undirected.
type Neighborhood[N Node, L Link[N]] interface {
	// Links returns the neighborhood as a way to iterate over links.
	// But no need to embed a full iterator if we want to have node metadata, so you may lazy load the iterator.
	// NOTE : there is absolutely no warranty about the links order !
	// For instance, some implementations use maps, so no order is provided when iterating over the values.
	Links() (LinksIterator[N, L], error)
	// IncomingDegree returns the number of nodes that have current node as their destination.
	// For undirected link, just use undirected degree
	IncomingDegree() int64
	// OutgoingDegree returns the number of nodes that have current node as their source.
	// For undirected link, just use undirected degree
	OutgoingDegree() int64
	// UndirectedDegree returns the degree of the node for undirected links
	UndirectedDegree() int64
	// CenterNode returns the node we get neighborhood for
	CenterNode() N
}

// IsIsolatedNeighborhood returns true for a node with no link, false otherwise.
// It returns true also for a default neighborhood
func IsIsolatedNeighborhood[N Node, L Link[N]](n Neighborhood[N, L]) bool {
	return n.IncomingDegree() == 0 && n.OutgoingDegree() == 0 && n.UndirectedDegree() == 0
}

// NeighborhoodIterator is just an iterator over all the neighboors of a node
type NeighborhoodIterator[N Node, L Link[N]] GeneralIterator[Neighborhood[N, L]]
