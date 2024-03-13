package main

import (
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/storage/gexf"
)

func main() {
	graph, errGraph := local.GenerateCompleteUndirectedGraph(10, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errGraph != nil {
		panic(errGraph)
	}

	errWrite := gexf.ExportDataGraph(
		"D:\\test.gexf",
		&graph,
		func(internal.IdNode) (string, map[string]string) { return "Hello", nil },
		gexf.GexfLinkBasicSerializer,
	)

	if errWrite != nil {
		panic(errWrite)
	}
}
