package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

// TestGnpProbaZero is gnp with proba = 0.
// It ensures that we know the graph: all nodes are disconnected
func TestGnpProbaZero(t *testing.T) {
	randomizer := local.RandomGenerator[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]{}
	result, errGraph := randomizer.GNP(5, 0.0, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errGraph != nil {
		t.Fail()
	}

	nodesCounter := 0
	it, errIt := result.AllNodes()
	if errIt != nil {
		t.Fail()
	}

	for has, err := it.Next(); has; has, err = it.Next() {
		if err != nil {
			t.Fail()
		}

		node, errNode := it.Value()
		if errNode != nil {
			t.Fail()
		}

		if neighbors, errNeighbors := result.Neighbors(node); errNeighbors != nil {
			t.Fail()
		} else if neighbors.UndirectedDegree() != 0 {
			t.Fail()
		}

		nodesCounter++
	}

	if nodesCounter != 5 {
		t.Fail()
	}
}

// TestGnpProbaOne is gnp with proba = 1.0.
// It means that graph should be complete
func TestGnpProbaOne(t *testing.T) {
	randomizer := local.RandomGenerator[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]{}
	result, errGraph := randomizer.GNP(5, 1.0, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errGraph != nil {
		t.Fail()
	}

	nodesCounter := 0
	it, errIt := result.AllNodes()
	if errIt != nil {
		t.Fail()
	}

	for has, err := it.Next(); has; has, err = it.Next() {
		if err != nil {
			t.Fail()
		}

		node, errNode := it.Value()
		if errNode != nil {
			t.Fail()
		}

		if neighbors, errNeighbors := result.Neighbors(node); errNeighbors != nil {
			t.Fail()
		} else if neighbors.UndirectedDegree() != 4 {
			t.Fail()
		}

		nodesCounter++
	}

	if nodesCounter != 5 {
		t.Fail()
	}
}
