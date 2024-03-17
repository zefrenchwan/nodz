package graphs_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestGraphConnectedComponents(t *testing.T) {
	graph := local.NewMapGraph[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]()

	// isolated node
	graph.AddNode(internal.NewRandomIdNode())

	c11 := internal.NewRandomIdNode()
	c12 := internal.NewRandomIdNode()
	c13 := internal.NewRandomIdNode()
	c14 := internal.NewRandomIdNode()

	c21 := internal.NewRandomIdNode()
	c22 := internal.NewRandomIdNode()
	c23 := internal.NewRandomIdNode()

	graph.AddLink(internal.NewUndirectedSimpleLink(c11, c12))
	graph.AddLink(internal.NewUndirectedSimpleLink(c12, c14))
	graph.AddLink(internal.NewUndirectedSimpleLink(c14, c13))

	graph.AddLink(internal.NewUndirectedSimpleLink(c21, c22))
	graph.AddLink(internal.NewUndirectedSimpleLink(c22, c23))
	graph.AddLink(internal.NewUndirectedSimpleLink(c23, c21))

	setBuilder := func(f graphs.SetEqualsFunction[internal.IdNode]) (graphs.AbstractSet[internal.IdNode], error) {
		result := local.NewSlicesSet(f)
		return &result, nil
	}

	itBuilder := func() (graphs.DynamicIterator[internal.IdNode], error) {
		result := local.NewDynamicSlicesIterator[internal.IdNode]()
		return &result, nil
	}

	if count, err := graphs.CountConnectedComponents(&graph, setBuilder, itBuilder); err != nil {
		t.Fail()
	} else if count != 3 {
		t.Error("expected 3 connected components")
	}
}
