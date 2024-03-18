package local

import (
	"errors"
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
)

// MapGraph contains the links for source and destination.
type MapGraph[N graphs.Node, L graphs.Link[N]] struct {
	// bijective slice to link nodes and their index
	nodes increasingMapping[N]
	// content is a map, keys are nodes index, values are nodes metadata and their outgoing links.
	// It does NOT contain all the nodes, just the ones with at least one link.
	content map[int]mapLine[N, L]
}

// NewMapGraph returns a new empty map matrix as a central structure graph
func NewMapGraph[N graphs.Node, L graphs.Link[N]]() MapGraph[N, L] {
	nodesTest := func(a, b N) bool {
		return a.SameNode(b)
	}

	return MapGraph[N, L]{
		nodes:   newIncreasingMapping(nodesTest),
		content: make(map[int]mapLine[N, L]),
	}
}

// AddNode adds a node if it did not exist, does nothing otherwise.
func (am *MapGraph[N, L]) AddNode(node N) error {
	am.nodes.addValue(node)
	return nil
}

// RemoveNode removes a node and all its links
func (am *MapGraph[N, L]) RemoveNode(node N) error {
	targetIndex, found := am.nodes.getValue(node)
	if !found {
		return nil
	}

	delete(am.content, targetIndex)
	for key, lines := range am.content {
		lines.removeNode(targetIndex)
		am.content[key] = lines
	}

	am.nodes.removeValue(node)

	return nil
}

// AddLink adds a link (may be directed or undirected) if not already here
func (am *MapGraph[N, L]) AddLink(link L) error {
	source := link.Source()
	dest := link.Destination()
	sourceIndex := am.nodes.addValue(source)
	destIndex := am.nodes.addValue(dest)
	am.setLink(sourceIndex, destIndex, link)

	return nil
}

func (am *MapGraph[N, L]) HasLink(link L) bool {
	source := link.Source()
	dest := link.Destination()
	sourceIndex := am.nodes.hasValue(source)
	destIndex := am.nodes.hasValue(dest)

	// basic case: nodes are not in the graph
	if sourceIndex == -1 || destIndex == -1 {
		return false
	}

	// test if the link from source to dest is in the graph
	isLinkSource := false
	if v, ok := am.content[sourceIndex]; ok {
		if l, lok := v.values[destIndex]; lok && len(l) != 0 {
			for _, currentLink := range l {
				if link.SameLink(currentLink) {
					isLinkSource = true
					break
				}
			}
		}
	}

	if isLinkSource {
		return true
	} else if link.IsDirected() {
		return false
	}

	// From here, link is not directed, so it may appear as the opposite
	if v, ok := am.content[destIndex]; ok {
		if l, lok := v.values[sourceIndex]; lok && len(l) != 0 {
			for _, currentLink := range l {
				if link.SameLink(currentLink) {
					return true
				}
			}
		}
	}

	// neither as (source, dest) nor as (dest, source)
	return false
}

// RemoveLink removes a link if any, does nothing otherwise
func (am *MapGraph[N, L]) RemoveLink(link L) error {
	source := link.Source()
	dest := link.Destination()
	sourceIndex, foundSource := am.nodes.getValue(source)
	destIndex, foundDest := am.nodes.getValue(dest)

	if !foundSource || !foundDest {
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
func (am *MapGraph[N, L]) AllNodes() (graphs.NodesIterator[N], error) {
	return am.nodes.toIterator(), nil
}

// Neighbors returns the neighborhood of the node (metadata and iterators factory)
func (am *MapGraph[N, L]) Neighbors(node N) (graphs.Neighborhood[N, L], error) {
	index, found := am.nodes.getValue(node)
	if !found {
		return nil, nil
	}

	mapValue := am.content[index]

	return mapValue.toNeighborhood(node), nil
}

// ToMatrix maps a graph as a matrix of nodes with values.
// N is nodes type
// L is link type
// S is type of the elements of the matrix
// linksMapper maps all links from i to j to an instance of S
// Result is mapping between indexes and nodes, and the matrix per se.
// And, at position (i,j), the instance of S is linksMapper of all the links from i to j.
// Remember that N is not comparable, so []N result is a way to say without a map: index of this node is this int.
// As an example of use, to get the adjacency matrix (the number of links from i to j),
// linksMapper would just count the number of directed links.
func ToMatrix[N graphs.Node, L graphs.Link[N], S any](am *MapGraph[N, L], linksMapper func([]L) S) ([]N, graphs.Matrix[S]) {
	if am == nil || linksMapper == nil {
		return nil, nil
	}

	// make result
	defaultValue := linksMapper(nil)
	expectedSize := am.nodes.size()
	// size has been tested before, cannot raise error
	result, _ := NewMapMatrix[S](expectedSize, defaultValue)

	// index of elements still in the map may not be contiguous ints (due to removes)
	matrixIndexMapping := am.nodes.toIncreasingIndexes()
	// for each link, find index of source and destination in matrix and fill matrix
	for sourceIndex, line := range am.content {
		matrixSourceIndex := matrixIndexMapping[sourceIndex]
		for destIndex, links := range line.values {
			matrixDestIndex := matrixIndexMapping[destIndex]
			result.SetValue(matrixSourceIndex, matrixDestIndex, linksMapper(links))
		}
	}

	arrayOfNodes := am.nodes.toIncreasingValues()
	return arrayOfNodes, &result
}

// GenerateCompleteUndirectedGraph returns a complete undirected graph with nodesSize nodes
func GenerateCompleteUndirectedGraph[N graphs.Node, L graphs.Link[N]](
	nodesSize int, // number of nodes in the result
	nodeGenerator graphs.RandomNodeGenerator[N], // generates a new random node at each call
	linkGenerator graphs.RandomLinkGenerator[N, L], // generates a new random undirected link at each call
) (
	MapGraph[N, L], // result
	error, // error, for instance if nodesSize makes no sense or whether linkGenerator makes directed links
) {
	var result MapGraph[N, L]

	if nodesSize < 0 {
		return result, errors.New("invalid size")
	} else if nodesSize == 0 {
		return result, nil
	}

	result = NewMapGraph[N, L]()
	// generate nodes
	nodes := make([]N, nodesSize)
	for index := 0; index < nodesSize; index++ {
		newNode := nodeGenerator()
		result.AddNode(newNode)
		nodes[index] = newNode
	}

	for i, source := range nodes {
		for j, dest := range nodes {
			if i == j {
				continue
			}

			link := linkGenerator(source, dest)
			if link.IsDirected() {
				return result, errors.New("undirected links only")
			}

			result.AddLink(link)
		}
	}

	return result, nil
}

// setLink adds a link at a given index
func (am *MapGraph[N, L]) setLink(sourceIndex, destIndex int, link L) {
	mapValue := am.content[sourceIndex]
	mapValue.addLink(destIndex, link)
	am.content[sourceIndex] = mapValue

	if !link.IsDirected() {
		mapValue = am.content[destIndex]
		mapValue.addLink(sourceIndex, link)
		am.content[destIndex] = mapValue
	} else {
		mapValue = am.content[destIndex]
		mapValue.incomingCounter = mapValue.incomingCounter + 1
		am.content[destIndex] = mapValue
	}
}

// mapLine is the outgoings links (or undirected links seen as outgoing links) and node stats.
// It makes no sense alone, it is always to consider with the node it represents (the source) in mind.
type mapLine[N graphs.Node, L graphs.Link[N]] struct {
	// incomingCounter keeps counter for directed incoming links
	incomingCounter int64
	// outgoingCounter keeps counter for directed outgoing links
	outgoingCounter int64
	// undirectedCounter keeps counter for undirected nodes seen as outgoing links
	undirectedCounter int64
	// Values are the links outgoing from source node.
	// But, in general, there may be many nodes with the same source and destination.
	// So we store index of destination node and all links with same source and destination.
	// In particular, for undirected links, they are stored twice.
	values map[int][]L
}

// removeNode removes the nodes and changes counters.
// So, set back the value in the graph map
func (a *mapLine[N, L]) removeNode(nodeIndex int) {
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

// addLink adds the link if not already there, returns true if one link was inserted, false otherwise.
// If there was an actual insertion, set back the value in the graph map
func (a *mapLine[N, L]) addLink(destinationIndex int, link L) bool {
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

// removeLink removes the link if any, returns true if one link was removed, false otherwise.
// If there was an actual deletion, set back the value in the graph map
func (a *mapLine[N, L]) removeLink(destinationIndex int, link L) bool {
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
func (a *mapLine[N, L]) toNeighborhood(center N) graphs.Neighborhood[N, L] {
	result := internal.NeighborsIterator[N, L]{}
	if a == nil {
		return result
	}

	result.IncomingCounter = a.incomingCounter
	result.OutgoingCounter = a.outgoingCounter
	result.UndirectedCounter = a.undirectedCounter
	result.CurrentNode = center
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
