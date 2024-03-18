package graphs

import "errors"

// NetworkStatistics is about the basic statistics of a network: distribution degree, nodes and links counters
type NetworkStatistics struct {
	// DegreeDistribution counts the number of nodes for a given degree, but normalized.
	// It means that we count the nodes matching the degree for each degree, then sum all nodes.
	// Then, we divide each degree by the total number of nodes
	DegreeDistribution map[int64]float64
	// NodesSize is the number of nodes in the graph
	NodesSize int64
	// UndirectedSize is the number of undirected links
	UndirectedSize int64
	// DirectedSize is the number of outgoing links (and the same as incoming links).
	DirectedSize int64
}

// AverageDirectedDegree returns -1 for 0 denominator, the average degree of the directed graph otherwise
func (n NetworkStatistics) AverageDirectedDegree() float64 {
	if n.NodesSize == 0 {
		return -1.0
	}

	return float64(n.DirectedSize) / float64(n.NodesSize)
}

// AverageUndirectedDegree returns -1 for 0 denominator, the average degree of the undirected graph otherwise
func (n NetworkStatistics) AverageUndirectedDegree() float64 {
	if n.NodesSize == 0 {
		return -1.0
	}

	return 2.0 * float64(n.UndirectedSize) / float64(n.NodesSize)
}

// DirectedDensity calculates the fraction of links compared to its possible max.
// Default value for graphs with nodesSize < 2 is 0.0
func (n NetworkStatistics) DirectedDensity() float64 {
	if n.NodesSize < 2 {
		return 0.0
	}

	nodesCount := float64(n.NodesSize)
	linksCount := float64(n.DirectedSize)

	return linksCount / (nodesCount * (nodesCount - 1.0))
}

// UndirectedDensity calculates the fraction of links compared to its possible max.
// Default value for graphs with nodesSize < 2 is 0.0
func (n NetworkStatistics) UndirectedDensity() float64 {
	if n.NodesSize < 2 {
		return 0.0
	}

	nodesCount := float64(n.NodesSize)
	linksCount := float64(n.UndirectedSize)

	return 2.0 * linksCount / (nodesCount * (nodesCount - 1.0))
}

// CalculateNetworkStatistics returns the basic stats of a network.
// Sizes are calculated the standard way.
// Depending on the type of graph you work on, you may want undirected degree of a node in the distribution, or incoming, or...
// To deal with those situations once, counter maps a neighborhood to an int, and degree distribution is based on its result.
// Algorithm is to go through all the nodes, so, for large graphs, it may take time.
func CalculateNetworkStatistics[N Node, L Link[N]](g CentralStructureGraph[N, L], counter func(Neighborhood[N, L]) int64) (NetworkStatistics, error) {
	var result NetworkStatistics
	result.DegreeDistribution = make(map[int64]float64)

	it, errIt := g.AllNodes()
	if errIt != nil {
		return result, errIt
	}
	// for a counter value, the number of matching nodes
	degreesCounter := make(map[int64]int64)

	var globalErr error
	for has, err := it.Next(); has; has, err = it.Next() {
		if err != nil {
			globalErr = errors.Join(globalErr, err)
		}

		var neighborhood Neighborhood[N, L]
		if node, errNode := it.Value(); errNode != nil {
			globalErr = errors.Join(globalErr, errNode)
		} else if neighbors, errNeighbors := g.Neighbors(node); errNeighbors != nil {
			globalErr = errors.Join(globalErr, err)
		} else {
			neighborhood = neighbors
		}

		countValue := counter(neighborhood)
		degreesCounter[countValue] = degreesCounter[countValue] + 1
		result.NodesSize += 1
		result.DirectedSize += neighborhood.OutgoingDegree()
		result.UndirectedSize += neighborhood.UndirectedDegree()
	}

	// undirected links were counted twice
	result.UndirectedSize /= 2

	// perform normalization over degree distribution
	for k, v := range degreesCounter {
		result.DegreeDistribution[k] = float64(v) / float64(result.NodesSize)
	}

	return result, globalErr
}
