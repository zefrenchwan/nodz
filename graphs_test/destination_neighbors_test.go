package graphs_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestDestinationNeighbors(t *testing.T) {
	graph := local.NewMapGraph[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]()

	source := internal.NewRandomIdNode()
	dest1 := internal.NewRandomIdNode()
	dest2 := internal.NewRandomIdNode()
	dest3 := internal.NewRandomIdNode()
	notInGraph := internal.NewRandomIdNode()
	isolated := internal.NewRandomIdNode()

	linkSourceDest1 := internal.NewUndirectedSimpleLink(source, dest1)
	linkSourceDest2 := internal.NewUndirectedSimpleLink(source, dest2)
	linkDest1Dest3 := internal.NewUndirectedSimpleLink(dest1, dest3)

	graph.AddNode(isolated)
	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkDest1Dest3)

	// not in graph, expects nil
	if n, err := graphs.DestinationNeighbors(notInGraph, &graph); err != nil || n != nil {
		t.Fail()
	}

	// isolated in graph, expects empty
	if n, err := graphs.DestinationNeighbors(isolated, &graph); err != nil || n == nil {
		t.Fail()
	} else if has, err := n.Next(); err != nil || has {
		t.Fail()
	}

	// in graph, expects all its neighbors.
	// Expected values are the neighborhood of source and dest3
	// First part, calculate expected
	expected := make([]graphs.Neighborhood[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]], 0)
	if n, err := graph.Neighbors(source); err != nil {
		t.Fail()
	} else {
		expected = append(expected, n)
	}

	if n, err := graph.Neighbors(dest3); err != nil {
		t.Fail()
	} else {
		expected = append(expected, n)
	}

	testNeighbors := func(a, b graphs.Neighborhood[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]) bool {
		return a.IncomingDegree() == b.IncomingDegree() &&
			a.OutgoingDegree() == b.OutgoingDegree() &&
			a.UndirectedDegree() == b.UndirectedDegree() &&
			a.CenterNode().SameNode(b.CenterNode())
	}

	// proceed to test
	if all, err := graphs.DestinationNeighbors(dest1, &graph); err != nil {
		t.Fail()
	} else if v, errV := internal_test.CompareIteratorWithSlice(all, expected, testNeighbors, false); errV != nil {
		t.Fail()
	} else if !v {
		t.Error("neighborhoods differ")
	}
}
