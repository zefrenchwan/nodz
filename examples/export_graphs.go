package examples

import (
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/storage/gexf"
)

// ExportGraphToGEXFFile writes a gexf file from a preferential attachment generated graph
func ExportGraphToGEXFFile() {
	var generator local.RandomGenerator[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]
	graph, errGraph := generator.UndirectedBarabasiAlbertGraph(10, 200, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errGraph != nil {
		panic(errGraph)
	}

	errWrite := gexf.ExportDataGraph(
		"D:\\test.gexf",
		graph,
		gexf.GexfBlankNodeExporter,
		gexf.GexfLinkBasicSerializer,
	)

	if errWrite != nil {
		panic(errWrite)
	}
}
