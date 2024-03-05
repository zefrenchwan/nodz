package local

import (
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// AdjacencyMatrix contains the links for source and destination.
// Definition is often true or false for non valued graphs, the values for valued graphs.
// This implementation stores the link as a whole : given a source and destination, stored value is the link.
// It means that, for some algorithms involving matrix operations, link has to be processed to return an int or a float, or...
type AdjacencyMatrix[N graphs.Node, L graphs.Link[N]] struct {
	// bijective slice to link nodes and their index
	nodes []N
	// basic implementation with nodes index and related links.
	// Invariant are:
	// for indexes i,j, content[i][j] contains all the links that have i and j as source and destination
	// directed links are inserted as is
	// UNdirected links are stored once, with smaller index then bigger index
	content doubleMap[int, []L]
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

// AddNode adds a node if it did not exist, does nothing otherwise.
func (am *AdjacencyMatrix[N, L]) AddNode(node N) error {
	index := am.nodeIndex(node)
	if index < 0 {
		am.nodes = append(am.nodes, node)
	}

	return nil
}

func (am *AdjacencyMatrix[N, L]) RemoveNode(node N) error {
	targetIndex := am.nodeIndex(node)
	if targetIndex < 0 {
		return nil
	}

	for index := range am.content {
		delete(am.content[index], targetIndex)
	}

	delete(am.content, targetIndex)
	return nil
}

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
	} else if link.IsDirected() && sourceIndex > destIndex {
		sourceIndex, destIndex = destIndex, sourceIndex
	}

	values, found := am.content.getValue(sourceIndex, destIndex)
	if !found || values == nil {
		return nil
	}

	values = slices.DeleteFunc(values, func(other L) bool {
		return link.SameLink(other)
	})

	am.content.putValue(sourceIndex, destIndex, values)
	return nil
}

func (am *AdjacencyMatrix[N, L]) AllNodes() (graphs.NodesIterator[N], error) {
	it := NewSlicesIterator(am.nodes)
	return &it, nil
}

func (am *AdjacencyMatrix[N, L]) setLink(i, j int, link L) {
	// deal with undirected links to avoid double insert
	sourceIndex, destIndex := i, j
	if link.IsDirected() && j > i {
		sourceIndex, destIndex = j, i
	}

	// add value if there is no same link already there
	previousLinks, found := am.content.getValue(sourceIndex, destIndex)
	if found && previousLinks != nil {
		for _, previousLink := range previousLinks {
			if link.SameLink(previousLink) {
				return
			}
		}
	}

	// either no previous value, or no same link inserted
	previousLinks = append(previousLinks, link)
	am.content.putValue(sourceIndex, destIndex, previousLinks)
}
