package graphs

// Neighborhood of a node is the links from that node.
// Still, no need to load all the links to have basinc information about the neighborhood.
// So, to include all options, three functions appear for degrees : incoming, outgoing, and undirected
type Neighborhood[N Node, L Link[N]] interface {
	// Neighbors returns an iterator over the links starting from the source
	Neighbors() (LinksIterator[N, L], error)
	// IncomingDegree returns the number of nodes that have current node as their destination.
	// For undirected link, just use undirected degree
	IncomingDegree() int64
	// OutgoingDegree returns the number of nodes that have current node as their source.
	// For undirected link, just use undirected degree
	OutgoingDegree() int64
	// Degree of the node for undirected links
	UndirectedDegree() int64
}
