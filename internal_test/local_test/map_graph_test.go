package local_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestNeighbors(t *testing.T) {
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	source := internal.NewPropertiesNode()
	dest1 := internal.NewPropertiesNode()
	dest2 := internal.NewPropertiesNode()
	dest3 := internal.NewPropertiesNode()
	notInGraph := internal.NewPropertiesNode()

	linkSourceDest1 := internal.NewDirectedValuedLink(&source, &dest1, 10)
	linkSourceDest2 := internal.NewDirectedValuedLink(&source, &dest2, 20)
	linkDest1Dest3 := internal.NewDirectedValuedLink(&dest1, &dest3, 30)

	graph.AddLink(linkSourceDest1)
	graph.AddLink(linkSourceDest2)
	graph.AddLink(linkDest1Dest3)

	var errG error
	var neighbors graphs.Neighborhood[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]
	//neighborhood of something not in the graph should return nothing
	neighbors, errG = graph.Neighbors(&notInGraph)
	if errG != nil {
		t.Error("getting neighbors of a non existing node should not return an error")
	} else if neighbors != nil {
		t.Error("expecting nil for no neighboors")
	}

	// neighborhood of source shoud be source -> dest1 and source->dest2
	neighbors, errG = graph.Neighbors(&source)
	if errG != nil {
		t.Error("getting neighbors should not return an error")
	} else if neighbors == nil {
		t.Error("expecting non nil for neighboors")
	} else if neighbors.OutgoingDegree() != 2 {
		t.Error("expecting two outgoing nodes")
	} else if neighbors.IncomingDegree() != 0 {
		t.Error("no incoming link expected")
	} else if !neighbors.CenterNode().SameNode(&source) {
		t.Error("cannot get center of neighborhood")
	}

	// dest1 has one incoming node and one outgoing node
	neighbors, errG = graph.Neighbors(&dest1)
	if errG != nil {
		t.Error("getting neighbors should not return an error")
	} else if neighbors == nil {
		t.Error("expecting non nil for neighboors")
	} else if neighbors.OutgoingDegree() != 1 {
		t.Error("expecting two outgoing nodes")
	} else if neighbors.IncomingDegree() != 1 {
		t.Error("no incoming link expected")
	} else if !neighbors.CenterNode().SameNode(&dest1) {
		t.Error("cannot get center of neighborhood")
	}
}

func TestAdjacencyMatrixNodes(t *testing.T) {
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

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
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

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
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

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
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

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

	it, errIt := graph.AllNodes()
	if errIt != nil {
		t.Fail()
	}

	localCompare := func(a, b *internal.PropertiesNode) bool {
		return a.SameNode(b)
	}

	expected := []*internal.PropertiesNode{&source, &dest1, &dest3}
	if res, errComp := internal_test.CompareIteratorWithSlice(it, expected, localCompare, false); !res || errComp != nil {
		t.Error("expected nodes do not match for graph")
	}
}

func TestMatrix(t *testing.T) {
	graph := local.NewMapGraph[*internal.PropertiesNode, internal.ValuedLink[*internal.PropertiesNode, int]]()

	dest := internal.NewPropertiesNode()
	source1 := internal.NewPropertiesNode()
	source2 := internal.NewPropertiesNode()
	alone := internal.NewPropertiesNode()

	linkSource1Dest := internal.NewDirectedValuedLink(&source1, &dest, 10)
	linkSource2Dest1 := internal.NewDirectedValuedLink(&source2, &dest, 20)
	linkSource2Dest2 := internal.NewDirectedValuedLink(&source2, &dest, 30)

	graph.AddLink(linkSource1Dest)
	graph.AddLink(linkSource2Dest1)
	graph.AddLink(linkSource2Dest2)
	graph.AddNode(&alone)

	counter := func(links []internal.ValuedLink[*internal.PropertiesNode, int]) int {
		return len(links)
	}

	mapping, matrix := local.ToMatrix(&graph, counter)

	if len(mapping) != 4 {
		t.Fail()
	}

	// find index of each node.
	// Node not comparable means no map, so .... Here we go....
	var indexSource1, indexSource2, indexDest int
	for index, node := range mapping {
		switch {
		case source1.SameNode(node):
			indexSource1 = index
		case source2.SameNode(node):
			indexSource2 = index
		case alone.SameNode(node):
			_ = index
		case dest.SameNode(node):
			indexDest = index
		default:
			t.Error("unexpected node")
		}
	}

	if matrix.Size() != 4 {
		t.Fail()
	}

	var value int
	if value, _, _ := matrix.GetValue(indexSource1, indexDest); value != 1 {
		t.Error("expected one link from source1 to dest")
	}

	if value, _, _ := matrix.GetValue(indexSource2, indexDest); value != 2 {
		t.Error("expected two links from source2 to dest")
	}

	// rest should be 0, so global sum should be 3
	sum := 0
	for i := 0; i < matrix.Size(); i++ {
		for j := 0; j < matrix.Size(); j++ {
			value, _, _ = matrix.GetValue(i, j)
			sum = sum + value
		}
	}

	if sum != 3 {
		t.Errorf("expected 3 as sum, got %d", sum)
	}
}

func TestCompleteGraphGeneration(t *testing.T) {
	result, errResult := local.GenerateCompleteUndirectedGraph[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]](10, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errResult != nil {
		t.Fail()
	}

	it, errIt := result.AllNodes()
	if errIt != nil {
		t.Fail()
	}

	counter := 0
	for has, errHas := it.Next(); has; has, errHas = it.Next() {
		if errHas != nil {
			t.Fail()
		}

		counter++

		if v, errV := it.Value(); errV != nil {
			t.Fail()
		} else if n, errN := result.Neighbors(v); errN != nil {
			t.Fail()
		} else if n.UndirectedDegree() != 9 {
			t.Fail()
		}
	}

	if counter != 10 {
		t.Fail()
	}
}
