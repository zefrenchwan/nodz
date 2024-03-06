package local_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestAdjacencyMatrixNodes(t *testing.T) {
	graph := local.NewAdjacencyMatrix[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	source := internal.NewPropertiesNode()
	dest1 := internal.NewPropertiesNode()
	dest2 := internal.NewPropertiesNode()
	dest3 := internal.NewPropertiesNode()
	intermed := internal.NewPropertiesNode()

	expectedNodes := []internal.PropertiesNode{
		source, dest1, dest2, dest3, intermed,
	}

	linkSourceDest1 := internal.NewDirectedValuedLink(&source, &dest1, 10)
	linkSourceDest2 := internal.NewDirectedValuedLink(&source, &dest2, 20)
	linkIntermedDest3 := internal.NewUndirectedValuedLink(&intermed, &dest3, 30)
	linkSourceIntermed := internal.NewUndirectedValuedLink(&source, &intermed, 40)

	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkIntermedDest3)
	graph.AddLink(linkSourceIntermed)

	allNodes := make([]internal.PropertiesNode, 0)

	if it, err := graph.AllNodes(); err != nil || it == nil {
		t.Error("failed to build iterator over nodes")
	} else {
		for found, err := it.Next(); found && err == nil; found, err = it.Next() {
			if v, errV := it.Value(); v == nil || errV != nil {
				t.Error("error while iterating over nodes, nil value or error")
			} else {
				allNodes = append(allNodes, *v)
			}
		}
	}

	if len(allNodes) != len(expectedNodes) {
		t.Error("missing nodes")
	}

	for _, expected := range expectedNodes {
		if !slices.ContainsFunc(allNodes, func(node internal.PropertiesNode) bool {
			return expected.SameNode(&node)
		}) {
			t.Error("some missing nodes")
		}
	}
}

func TestAdjacencyMatrixLinks(t *testing.T) {
	graph := local.NewAdjacencyMatrix[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	source := internal.NewPropertiesNode()
	dest1 := internal.NewPropertiesNode()
	dest2 := internal.NewPropertiesNode()
	dest3 := internal.NewPropertiesNode()
	intermed := internal.NewPropertiesNode()
	nonExistingNode := internal.NewPropertiesNode()

	linkSourceDest1 := internal.NewDirectedValuedLink(&source, &dest1, 10)
	linkSourceDest2 := internal.NewDirectedValuedLink(&source, &dest2, 20)
	linkIntermedDest3 := internal.NewUndirectedValuedLink(&intermed, &dest3, 30)
	linkSourceIntermed := internal.NewUndirectedValuedLink(&source, &intermed, 40)

	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkIntermedDest3)
	graph.AddLink(linkSourceIntermed)

	var neighbors graphs.Neighborhood[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]
	// test node not in the graph
	neighbors, errN := graph.Neighbors(&nonExistingNode)
	if errN != nil {
		t.Error("error when getting neighbors")
	} else if neighbors != nil {
		t.Error("neighbors finds values for non existing node")
	}

	// test directed link no incoming
	neighbors, errN = graph.Neighbors(&source)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 0 {
		t.Errorf("error for source incoming degree. Expected 0, got %d", neighbors.IncomingDegree())
	} else if neighbors.OutgoingDegree() != 2 {
		t.Errorf("error for source outgoing degree. Expected 2, got %d", neighbors.OutgoingDegree())
	} else if neighbors.UndirectedDegree() != 1 {
		t.Errorf("error for source undirected degree. Expected 1, got %d", neighbors.UndirectedDegree())
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else {
		// order may differ, no iteration
		localCompare := func(a, b internal.ValuedLink[*internal.PropertiesNode, int]) bool {
			return a.SameLink(b)
		}

		expected := []internal.ValuedLink[*internal.PropertiesNode, int]{linkSourceDest1, linkSourceDest2, linkSourceIntermed}
		if res, errComp := internal_test.CompareIteratorWithSlice(it, expected, localCompare, false); !res || errComp != nil {
			t.Error("expected links do not match for source")
		}
	}

	// test directed link no outgoing
	neighbors, errN = graph.Neighbors(&dest2)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 1 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 0 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); has || errHas != nil {
		t.Fail()
	}

	// test directed link no incoming
	neighbors, errN = graph.Neighbors(&intermed)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 0 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 2 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); !has || errHas != nil {
		t.Fail()
	} else if v, errV := it.Value(); errV != nil || !linkIntermedDest3.SameLink(v) {
		t.Fail()
	} else if has, errHas = it.Next(); !has || errHas != nil {
		t.Fail()
	} else if v, errV = it.Value(); errV != nil || !linkSourceIntermed.SameLink(v) {
		t.Fail()
	} else if has, errHas = it.Next(); has || errHas != nil {
		t.Fail()
	}
}

func TestAdjacencyMatrixRemoveLink(t *testing.T) {
	graph := local.NewAdjacencyMatrix[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	source := internal.NewPropertiesNode()
	dest1 := internal.NewPropertiesNode()
	dest2 := internal.NewPropertiesNode()
	dest3 := internal.NewPropertiesNode()
	intermed := internal.NewPropertiesNode()

	linkSourceDest1 := internal.NewDirectedValuedLink(&source, &dest1, 10)
	linkSourceDest2 := internal.NewDirectedValuedLink(&source, &dest2, 20)
	linkIntermedDest3 := internal.NewUndirectedValuedLink(&intermed, &dest3, 30)
	linkSourceIntermed := internal.NewUndirectedValuedLink(&source, &intermed, 40)
	linkDest1Dest3 := internal.NewUndirectedValuedLink(&dest1, &dest3, 50)
	linkDest3Dest2 := internal.NewDirectedValuedLink(&dest3, &dest2, 60)

	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkDest1Dest3)
	graph.AddLink(linkDest3Dest2)
	graph.AddLink(linkIntermedDest3)
	graph.AddLink(linkSourceIntermed)

	graph.RemoveLink(linkDest3Dest2)
	graph.RemoveLink(linkDest1Dest3)

	var neighbors graphs.Neighborhood[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]

	// no change on source
	neighbors, errN := graph.Neighbors(&source)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 0 {
		t.Error("source stats failure")
	} else if neighbors.OutgoingDegree() != 2 {
		t.Error("source stats failure")
	} else if neighbors.UndirectedDegree() != 1 {
		t.Error("source stats failure")
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else {
		// order may differ, no iteration
		localCompare := func(a, b internal.ValuedLink[*internal.PropertiesNode, int]) bool {
			return a.SameLink(b)
		}

		expected := []internal.ValuedLink[*internal.PropertiesNode, int]{linkSourceDest1, linkSourceDest2, linkSourceIntermed}
		if res, errComp := internal_test.CompareIteratorWithSlice(it, expected, localCompare, false); !res || errComp != nil {
			t.Error("expected links do not match for source")
		}
	}

	// dest2 should have one incoming node
	neighbors, errN = graph.Neighbors(&dest2)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 1 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 0 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); has || errHas != nil {
		t.Fail()
	}

	// dest1 should have one incoming
	neighbors, errN = graph.Neighbors(&dest1)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 1 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 0 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); has || errHas != nil {
		t.Fail()
	}

	// dest3 should have one undirected
	neighbors, errN = graph.Neighbors(&dest3)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 0 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 1 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); !has || errHas != nil {
		t.Fail()
	} else if v, errV := it.Value(); errV != nil || !linkIntermedDest3.SameLink(v) {
		t.Fail()
	} else if has, errHas = it.Next(); has || errHas != nil {
		t.Fail()
	}
}

func TestAdjacencyMatrixRemoveNode(t *testing.T) {
	graph := local.NewAdjacencyMatrix[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	source := internal.NewPropertiesNode()
	dest1 := internal.NewPropertiesNode()
	dest2 := internal.NewPropertiesNode()
	dest3 := internal.NewPropertiesNode()
	intermed := internal.NewPropertiesNode()

	linkSourceDest1 := internal.NewDirectedValuedLink(&source, &dest1, 10)
	linkSourceDest2 := internal.NewDirectedValuedLink(&source, &dest2, 20)
	linkIntermedDest3 := internal.NewUndirectedValuedLink(&intermed, &dest3, 30)
	linkSourceIntermed := internal.NewUndirectedValuedLink(&source, &intermed, 40)
	linkDest1Dest3 := internal.NewUndirectedValuedLink(&dest1, &dest3, 50)
	linkDest3Dest2 := internal.NewDirectedValuedLink(&dest3, &dest2, 60)

	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkDest1Dest3)
	graph.AddLink(linkDest3Dest2)
	graph.AddLink(linkIntermedDest3)
	graph.AddLink(linkSourceIntermed)

	// remove intermed means that graph should be
	// source --> dest1 -- dest3 and source --> dest2
	graph.RemoveNode(&intermed)

	var neighbors graphs.Neighborhood[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]

	// source expects two directed outgoing nodes
	neighbors, errN := graph.Neighbors(&source)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 0 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 2 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 0 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else {
		// order may differ, no iteration
		localCompare := func(a, b internal.ValuedLink[*internal.PropertiesNode, int]) bool {
			return a.SameLink(b)
		}

		expected := []internal.ValuedLink[*internal.PropertiesNode, int]{linkSourceDest1, linkSourceDest2}
		if res, errComp := internal_test.CompareIteratorWithSlice(it, expected, localCompare, false); !res || errComp != nil {
			t.Error("expected links do not match for source")
		}
	}

	graph.RemoveNode(&dest2)
	// graph should now be source --> dest1 -- dest3
	neighbors, errN = graph.Neighbors(&dest1)
	if errN != nil {
		t.Fail()
	} else if neighbors.IncomingDegree() != 1 {
		t.Fail()
	} else if neighbors.OutgoingDegree() != 0 {
		t.Fail()
	} else if neighbors.UndirectedDegree() != 1 {
		t.Fail()
	} else if it, errIt := neighbors.Links(); errIt != nil {
		t.Fail()
	} else if has, errHas := it.Next(); !has || errHas != nil {
		t.Fail()
	} else if v, errV := it.Value(); errV != nil || !linkDest1Dest3.SameLink(v) {
		t.Fail()
	} else if has, errHas = it.Next(); has || errHas != nil {
		t.Fail()
	}
}
