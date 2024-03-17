package graphs

import (
	"errors"
)

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

// StructuredGraph just says: this graph knows about its neighbors given a node.
// This is the common denominator with central structure graphs (basically, graph deals with nodes and links)
// and peers graphs (basically, nodes deal with graph structure without a "all mighty" common structure).
type StructuredGraph[N Node, L Link[N]] interface {
	// Neighbors returns the neighborhood of a node.
	// Formally, it returns nil if the node is NOT in the graph.
	// It returns an isolated non empty neighborhood for an isolated node.
	// It returns a full neighborhood for a linked node.
	Neighbors(N) (Neighborhood[N, L], error)
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
	// a central structure graph is a structured graph
	StructuredGraph[N, L]
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
}

// CountConnectedComponents applies to undirected graphs and returns the number of connected components of the graph.
// General algorithm is simple:
// Create a counter set to 0
// Mark all nodes in the graph as not seen
// Peek a node
// Find its connected component, and unmark each node of the same connected component
// If there is a left node, do the same with it
// Otherwise, return the counter of connected components
//
// To find the connected components of a node, we use a breadth first search.
// Algorithm is:
// create a fifo
// add first node (the one we start connected component with)
// while there is a node in the fifo
// pick it, and add all its neighbors in the fifo
func CountConnectedComponents[N Node, L Link[N]](
	graph CentralStructureGraph[N, L], // graph to find connected components within
	setBuilder AbstractSetBuilder[N], // to make a set implementation able to deal with the graph
	dynamicBuilder DynamicIteratorBuilder[N], // to make a dynamic builder able to deal with the graph
) (int64, // number of connected components in the graph
	error, // for any error
) {
	itNodes, errItNodes := graph.AllNodes()
	if errItNodes != nil {
		return -1, errItNodes
	}

	markedNodes, errSet := setBuilder(func(a, b N) bool { return a.SameNode(b) })
	if errSet != nil {
		return -1, errSet
	}

	var globalErr error
	for has, errHas := itNodes.Next(); has; has, errHas = itNodes.Next() {
		if errHas != nil {
			globalErr = errors.Join(globalErr, errHas)
			continue
		}

		if v, errV := itNodes.Value(); errV != nil {
			globalErr = errors.Join(globalErr, errV)
			continue
		} else if err := markedNodes.Add(v); err != nil {
			globalErr = errors.Join(globalErr, err)
			continue
		}
	}

	if globalErr != nil {
		return -1, globalErr
	}

	// counter of connected components
	var result int64
	for {
		if empty, err := markedNodes.IsEmpty(); err != nil {
			globalErr = errors.Join(globalErr, err)
			break
		} else if empty {
			break
		}

		nextOne, errPop := markedNodes.Peek()
		if errPop != nil {
			globalErr = errors.Join(globalErr, errPop)
			continue
		}

		fifo, errFifo := dynamicBuilder()
		if errFifo != nil {
			globalErr = errors.Join(globalErr, errPop)
			break
		}

		if err := fifo.AddNextValue(nextOne); err != nil {
			globalErr = errors.Join(globalErr, err)
			continue
		}

		// find its connected component using a breadth first search
		for {
			// node to find connected component for
			var currentNode N
			// go until fifo is empty
			if next, errNext := fifo.Next(); errNext != nil || !next {
				globalErr = errors.Join(globalErr, errNext)
				break
			} else if v, errV := fifo.Value(); errV != nil {
				globalErr = errors.Join(globalErr, errNext)
				continue
			} else if has, errHas := markedNodes.Has(v); errHas != nil || !has {
				globalErr = errors.Join(globalErr, errHas)
				continue
			} else {
				currentNode = v
			}

			// unmark the node
			if err := markedNodes.Remove(currentNode); err != nil {
				globalErr = errors.Join(globalErr, err)
				break
			}

			// next step is to keep walking through the graph by adding neighbors to the fifo
			var neighbors NeighborhoodIterator[N, L]
			if n, errN := DestinationNeighbors(currentNode, graph); errN != nil {
				globalErr = errors.Join(globalErr, errN)
				continue
			} else {
				neighbors = n
			}

			for has, errHas := neighbors.Next(); has; has, errHas = neighbors.Next() {
				if errHas != nil {
					globalErr = errors.Join(globalErr, errHas)
					continue
				} else if v, errV := neighbors.Value(); errV != nil {
					globalErr = errors.Join(globalErr, errHas)
					continue
				} else if errAdd := fifo.AddLastValue(v.CenterNode()); errAdd != nil {
					globalErr = errors.Join(globalErr, errAdd)
					continue
				}
			}
		}

		if globalErr != nil {
			return -1, globalErr
		}

		// else, we went through a full connected component, move to the next one
		result++
	}

	return result, globalErr
}

// DestinationNeighbors returns all the neighbors of the destinations of each link from origin.
// Get a node, get its links, get destinations, and return the neighborhoods of those destination.
// As an example, consider this graph:
// a -- b -- c
// Result of DestinationNeighbors(b, graph) is { neighborhood of a , neighborhood of b }
//
// Now, formally:
// * if node is not in the graph, it returns nil, nil.
// * if the node is isolated, result is an empty iterator, nil
// * otherwise, node is not isolated and all its neighbors form a set that is the expected result
func DestinationNeighbors[N Node, L Link[N]](origin N, graph StructuredGraph[N, L]) (NeighborhoodIterator[N, L], error) {
	if graph == nil {
		return nil, errors.New("nil graph")
	}

	neighbors, errNeighbors := graph.Neighbors(origin)
	if errNeighbors != nil {
		return nil, errNeighbors
	} else if neighbors == nil {
		// node not in the graph
		return nil, nil
	}

	links, errLinks := neighbors.Links()
	if errLinks != nil {
		return nil, errLinks
	}

	var result MapFilterIterator[L, Neighborhood[N, L]]
	result.Iterator = links
	result.Filter = func(n Neighborhood[N, L]) bool { return !IsIsolatedNeighborhood(n) }
	result.Mapper = func(link L) (Neighborhood[N, L], error) {
		_, destinationNode := FollowLink(origin, link)
		if n, err := graph.Neighbors(destinationNode); err != nil {
			return nil, err
		} else {
			return n, nil
		}
	}

	return &result, nil
}
