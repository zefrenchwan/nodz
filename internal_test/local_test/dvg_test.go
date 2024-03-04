package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestDirectedValuesGraph(t *testing.T) {
	graph := local.NewDirectedValuesGraph[string, int]()

	if value, found, err := graph.LinkValue("a", "b"); err != nil || found || value != 0 {
		t.Fail()
	}

	graph.AddNode("orphan")
	if value, found, err := graph.LinkValue("orphan", "a"); err != nil || found || value != 0 {
		t.Fail()
	}

	// graph is directed
	graph.SetLink("a", "b", 50)
	if value, found, err := graph.LinkValue("a", "b"); err != nil || !found || value != 50 {
		t.Fail()
	} else if value, found, err := graph.LinkValue("b", "a"); err != nil || found || value != 0 {
		t.Fail()
	}

	// adding an existing node should change nothing
	graph.AddNode("a")
	if value, found, err := graph.LinkValue("a", "b"); err != nil || !found || value != 50 {
		t.Fail()
	} else if value, found, err := graph.LinkValue("b", "a"); err != nil || found || value != 0 {
		t.Fail()
	}

	// removing link should be visible, and not affect other links
	graph.SetLink("a", "c", 60)
	graph.RemoveLink("a", "b")
	if value, found, err := graph.LinkValue("a", "c"); err != nil || !found || value != 60 {
		t.Fail()
	} else if value, found, err := graph.LinkValue("a", "b"); err != nil || found || value != 0 {
		t.Fail()
	}

	// test neighbors
	graph.SetLink("a", "d", 80)
	if n, err := graph.Neighbors("a"); err != nil {
		t.Fail()
	} else if n["d"] != 80 {
		t.Fail()
	} else if n["c"] != 60 {
		t.Fail()
	}

	// remove linked node
	graph.RemoveNode("a")
	if value, found, err := graph.LinkValue("a", "c"); err != nil || found || value != 0 {
		t.Fail()
	} else if value, found, err := graph.LinkValue("a", "d"); err != nil || found || value != 0 {
		t.Fail()
	}
}
