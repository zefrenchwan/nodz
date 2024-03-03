package internal_test

import (
	"slices"
	"testing"

	"github.com/zefrenchwan/nodz.git/internal"
)

func TestEqualityForLPN(t *testing.T) {
	node := internal.NewLabelPropertiesNode()
	otherNode := internal.NewLabelPropertiesNode()

	if node.SameNode(nil) {
		t.Fail()
	}

	if node.SameNode(&otherNode) {
		t.Fail()
	}

	if !node.SameNode(&node) {
		t.Fail()
	}
}

func TestLabelsForLPN(t *testing.T) {
	node := internal.NewLabelPropertiesNode()

	node.AddLabel("b")
	node.AddLabel("a")
	node.AddLabel("c")

	if slices.Compare(node.Labels(), []string{"a", "b", "c"}) != 0 {
		t.Fail()
	}

	node.RemoveLabel("a")
	if slices.Compare(node.Labels(), []string{"b", "c"}) != 0 {
		t.Fail()
	}
}

func TestPropertiesForLPN(t *testing.T) {
	node := internal.NewLabelPropertiesNode()

	if v, found := node.GetProperty("a"); v != "" || found {
		t.Fail()
	}

	node.SetProperty("key", "value")

	if v, found := node.GetProperty("key"); v != "value" || !found {
		t.Fail()
	}

	node.RemoveProperty("key")

	if v, found := node.GetProperty("key"); v != "" || found {
		t.Fail()
	}
}
