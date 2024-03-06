package local

import (
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
)

// AdjacencyMatrix contains the links for source and destination.
// Definition is often true or false for non valued graphs, the values for valued graphs.
// This implementation stores the link as a whole : given a source and destination, stored value is the link.
// It means that, for some algorithms involving matrix operations, link has to be processed to return an int or a float, or...
type AdjacencyMatrix[N graphs.Node, L graphs.Link[N]] struct {
	// bijective slice to link nodes and their index
	nodes []N
	// content is a map, keys are nodes index, values are nodes metadata and their outgoing links
	content map[int]adjacencyItem[N, L]
}

// NewAdjacencyMatrix returns a new empty adjacency matrix as a central structure graph
func NewAdjacencyMatrix[N graphs.Node, L graphs.Link[N]]() AdjacencyMatrix[N, L] {
	return AdjacencyMatrix[N, L]{
		nodes:   make([]N, 0),
		content: make(map[int]adjacencyItem[N, L]),
	}
}

// AddNode adds a node if it did not exist, does nothing otherwise.
func (am *AdjacencyMatrix[N, L]) AddNode(node N) error {
	index := am.nodeIndex(node)
	if index < 0 {
		am.nodes = append(am.nodes, node)
	}

	return nil
}

// RemoveNode removes a node and all its links
func (am *AdjacencyMatrix[N, L]) RemoveNode(node N) error {
	targetIndex := am.nodeIndex(node)
	if targetIndex < 0 {
		return nil
	}

	delete(am.content, targetIndex)
	for key, lines := range am.content {
		lines.removeNode(targetIndex)
		am.content[key] = lines
	}

	return nil
}

// AddLink adds a link (may be directed or undirected) if not already here
func (am *AdjacencyMatrix[N, L]) AddLink(link L) error {
	source := link.Source()
	dest := link.Destination()
	sourceIndex := am.nodeIndex(source)
	destIndex := am.nodeIndex(dest)

	if sourceIndex < 0 {
		sourceIndex = len(am.nodes)
		am.nodes = append(am.nodes, source)
	}

	if destIndex < 0 {
		destIndex = len(am.nodes)
		am.nodes = append(am.nodes, dest)
	}

	am.setLink(sourceIndex, destIndex, link)
	return nil
}

func (am *AdjacencyMatrix[N, L]) RemoveLink(link L) error {
	source := link.Source()
	dest := link.Destination()
	sourceIndex := am.nodeIndex(source)
	destIndex := am.nodeIndex(dest)

	if sourceIndex < 0 || destIndex < 0 {
		return nil
	}

	line := am.content[sourceIndex]
	removed := line.removeLink(destIndex, link)
	am.content[sourceIndex] = line

	switch {
	case removed && link.IsDirected():
		line := am.content[destIndex]
		line.incomingCounter = line.incomingCounter - 1
		am.content[destIndex] = line
	case removed && !link.IsDirected():
		line := am.content[destIndex]
		line.removeLink(sourceIndex, link)
		am.content[destIndex] = line
	}

	return nil
}

// AllNodes returns an iterator over nodes
func (am *AdjacencyMatrix[N, L]) AllNodes() (graphs.NodesIterator[N], error) {
	it := NewSlicesIterator(am.nodes)
	return &it, nil
}

// Neighbors returns the neighborhood of the node (metadata and iterators factory)
func (am *AdjacencyMatrix[N, L]) Neighbors(node N) (graphs.Neighborhood[N, L], error) {
	index := am.nodeIndex(node)
	if index < 0 {
		return nil, nil
	}

	adjacencyValue := am.content[index]

	return adjacencyValue.toNeighborhood(), nil
}

// nodeIndex gets the index of a node, -1 for no node found
func (am *AdjacencyMatrix[N, L]) nodeIndex(node N) int {
	if am == nil || am.nodes == nil {
		return -1
	}

	for index, value := range am.nodes {
		if node.SameNode(value) {
			return index
		}
	}

	return -1
}

// setLink adds a link at a given index
func (am *AdjacencyMatrix[N, L]) setLink(sourceIndex, destIndex int, link L) {
	adjacencyValue := am.content[sourceIndex]
	adjacencyValue.addLink(destIndex, link)
	am.content[sourceIndex] = adjacencyValue

	if !link.IsDirected() {
		adjacencyValue = am.content[destIndex]
		adjacencyValue.addLink(sourceIndex, link)
		am.content[destIndex] = adjacencyValue
	} else {
		adjacencyValue = am.content[destIndex]
		adjacencyValue.incomingCounter = adjacencyValue.incomingCounter + 1
		am.content[destIndex] = adjacencyValue
	}
}

type adjacencyItem[N graphs.Node, L graphs.Link[N]] struct {
	incomingCounter   int64
	outgoingCounter   int64
	undirectedCounter int64
	values            map[int][]L
}

func (a *adjacencyItem[N, L]) removeNode(nodeIndex int) {
	if a == nil || a.values == nil {
		return
	}

	allLinks, found := a.values[nodeIndex]
	if !found || allLinks == nil {
		return
	}

	var countDirected, countUndirected int64
	for _, link := range allLinks {
		// source is implicit, destination is nodeIndex
		if link.IsDirected() {
			countDirected++
		} else {
			countUndirected++
		}
	}

	delete(a.values, nodeIndex)
	// source is not nodeIndex, and we found all links such as destination is nodeIndex.
	// So, decrease the incoming counter from the source point of view
	a.incomingCounter = a.incomingCounter - countDirected
	a.undirectedCounter = a.undirectedCounter - countUndirected
}

func (a *adjacencyItem[N, L]) addLink(destinationIndex int, link L) bool {
	if a == nil {
		return false
	} else if a.values == nil {
		a.values = make(map[int][]L)
	}

	added := false

	links, found := a.values[destinationIndex]
	if !found || links == nil {
		a.values[destinationIndex] = []L{link}
		added = true
	} else {
		sameLinkFound := false
		for _, l := range a.values[destinationIndex] {
			if link.SameLink(l) {
				sameLinkFound = true
				break
			}
		}

		if !sameLinkFound {
			a.values[destinationIndex] = append(a.values[destinationIndex], link)
			added = true
		}
	}

	directed := link.IsDirected()
	switch {
	case added && directed:
		a.outgoingCounter = a.outgoingCounter + 1
	case added && !directed:
		a.undirectedCounter = a.undirectedCounter + 1
	}

	return added
}

func (a *adjacencyItem[N, L]) removeLink(destinationIndex int, link L) bool {
	if a == nil {
		return false
	}

	links, found := a.values[destinationIndex]
	if !found || links == nil {
		return false
	}

	sizeBefore := len(links)
	links = slices.DeleteFunc(links, func(l L) bool {
		return link.SameLink(l)
	})
	erased := len(links) != sizeBefore
	directed := link.IsDirected()
	switch {
	case erased && directed:
		a.outgoingCounter = a.outgoingCounter - 1
	case erased && !directed:
		a.undirectedCounter = a.undirectedCounter - 1
	}

	a.values[destinationIndex] = links

	return erased
}

// toNeighborhood constructs a new neighborhood for a given source
func (a *adjacencyItem[N, L]) toNeighborhood() graphs.Neighborhood[N, L] {
	result := internal.NeighborsIterator[N, L]{}
	if a == nil {
		return result
	}

	result.IncomingCounter = a.incomingCounter
	result.OutgoingCounter = a.outgoingCounter
	result.UndirectedCounter = a.undirectedCounter
	// make the union of all links. Because destinations are disjoin, links are always different
	result.IteratorsFactory = func() graphs.LinksIterator[N, L] {
		allLinks := make([]L, 0)
		for _, links := range a.values {
			if len(links) > 0 {
				allLinks = append(allLinks, links...)
			}
		}

		it := NewSlicesIterator(allLinks)
		return &it
	}

	return result
}
