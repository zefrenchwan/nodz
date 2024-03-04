package internal_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal"
)

func TestTypePropertiesLinkCreation(t *testing.T) {
	a := internal.NewLabelsPropertiesNode()
	b := internal.NewLabelsPropertiesNode()
	link := internal.NewTypePropertiesLink("link", &a, &b)
	reverseLink := internal.NewTypePropertiesLink("link", &b, &a)

	if !a.SameNode(link.Source()) {
		t.Fail()
	}

	if !b.SameNode(link.Destination()) {
		t.Fail()
	}

	if link.SameLink(reverseLink) {
		t.Fail()
	}

	if link.LinkType() != "link" {
		t.Fail()
	}
}

func TestTypePropertiesLinkProperties(t *testing.T) {
	a := internal.NewLabelsPropertiesNode()
	b := internal.NewLabelsPropertiesNode()
	link := internal.NewTypePropertiesLink("link", &a, &b)

	link.SetProperty("k", "v")

	if v, ok := link.GetProperty("k"); v != "v" || !ok {
		t.Fail()
	}

	link.RemoveProperty("z")

	if v, ok := link.GetProperty("k"); v != "v" || !ok {
		t.Fail()
	}

	link.RemoveProperty("k")

	if v, ok := link.GetProperty("k"); v != "" || ok {
		t.Fail()
	}

}
