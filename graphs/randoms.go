package graphs

// RandomNodeGenerator generates a new random node
type RandomNodeGenerator[N Node] func() N

// RandomLinkGenerator generates a new random link from a given source to a given destination.
// Note that link may not be directed, it is a random link, nothing more expected.
type RandomLinkGenerator[N Node, L Link[N]] func(source, destination N) L

// RandomGraphGenerator generates a random graph.
// There are many kinds of random graphs.
// So interface defintion may change in the future to include new kinds.
type RandomGraphGenerator[N Node, L Link[N]] interface {
	// DirectedGNP returns a graph with size nodes, each couple of nodes is linked with a given probabilty.
	// Each node and each link are generated via generators (one for nodes, one for links).
	// Special case for size == 0, result is an empty graph.
	// Result is either the graph, or an error.
	// For instance, a negative size returns an error, or a probability outside [0,1] makes an error too.
	//
	// ATTENTION: expected number of links is probability * ( size - 1)
	DirectedGNP(size int, probability float64, nodeGenerator RandomNodeGenerator[N], linkGenerator RandomLinkGenerator[N, L]) (CentralStructureGraph[N, L], error)

	// UndirectedGNP returns a graph with size nodes.
	// Each node and each link are generated via generators (one for nodes, one for links).
	// Special case for size == 0, result is an empty graph.
	// Result is either the graph, or an error.
	// For instance, a negative size returns an error, or a probability outside [0,1] makes an error too.
	//
	// ATTENTION: expected number of links is probability * 0.5 * (size - 1).
	// Given a pair (source, dest), either source - dest or dest - source is tested, not both.
	UndirectedGNP(size int, probability float64, nodeGenerator RandomNodeGenerator[N], linkGenerator RandomLinkGenerator[N, L]) (CentralStructureGraph[N, L], error)

	// UndirectedPreferentialAttachement returns a graph based on the Barabasi Albert model.
	// All links are UNdirected.
	// From a complete graph with initialSize nodes, add nodes until maxSize is reached.
	// Each link goes from a new node to node i, with a probability deg(i) / sum all deg
	UndirectedBarabasiAlbertGraph(initialSize int, maxSize int, nodeGenerator RandomNodeGenerator[N], linkGenerator RandomLinkGenerator[N, L]) (CentralStructureGraph[N, L], error)
}
