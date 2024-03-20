package examples

import (
	"errors"
	"fmt"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

// InheritanceDemo presents a simple breadth first walkthrough.
// Nodes are classes, links are "extends", and we test transitive relations.
func InheritanceTreeDemo() {
	// build the classes
	humans := internal.NewLabelsPropertiesNode()
	humans.AddLabel("humans")
	humans.SetProperty("name", "string")
	mortals := internal.NewLabelsPropertiesNode()
	mortals.AddLabel("mortals")
	entities := internal.NewLabelsPropertiesNode()
	entities.AddLabel("entities")
	// build the links (extends)
	humansMortalsLink := internal.NewTypePropertiesLink("extends", &humans, &mortals)
	mortalsEntitiesLink := internal.NewTypePropertiesLink("extends", &mortals, &entities)
	// build the tree
	inheritanceTree := local.NewMapGraph[*internal.LabelsPropertiesNode, *internal.TypePropertiesLink[*internal.LabelsPropertiesNode]]()
	inheritanceTree.AddLink(&humansMortalsLink)
	inheritanceTree.AddLink(&mortalsEntitiesLink)

	// put all superclasses of humans in a set
	fifo := local.NewDynamicSlicesIterator[*internal.LabelsPropertiesNode]()
	fifo.AddNextValue(&humans)
	superclasses := local.NewSlicesSet(func(a, b string) bool { return a == b })

	// not including error processing on purpose
	for has, _ := fifo.Next(); has; has, _ = fifo.Next() {
		current, _ := fifo.Value()
		neighbors, _ := graphs.DestinationNeighbors(current, &inheritanceTree)

		for h, _ := neighbors.Next(); h; h, _ = neighbors.Next() {
			// add in set
			v, _ := neighbors.Value()
			title := v.CenterNode().Labels()[0]
			superclasses.Add(title)
			// add in fifo to proceed
			fifo.AddLastValue(v.CenterNode())
		}
	}

	// and humans are mortal
	if has, _ := superclasses.Has("entities"); has {
		fmt.Println("Humans are entities (and mortal)")
	} else {
		panic(errors.New("inheritance failure"))
	}
}
