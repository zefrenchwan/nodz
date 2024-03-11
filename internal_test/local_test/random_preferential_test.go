package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestRandomBarabasiAlbertGraph(t *testing.T) {
	randomizer := local.RandomGenerator[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]{}
	result, errResult := randomizer.UndirectedBarabasiAlbertGraph(10, 20, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
	if errResult != nil {
		t.Fail()
	}

	it, errIt := result.AllNodes()
	if errIt != nil {
		t.Fail()
	}

	counter := 0
	bigNodesCounter := 0

	for has, errHas := it.Next(); has; has, errHas = it.Next() {
		if errHas != nil {
			t.Fail()
		}

		counter++

		if v, errV := it.Value(); errV != nil {
			t.Fail()
		} else if n, errN := result.Neighbors(v); errN != nil {
			t.Fail()
		} else if n.UndirectedDegree() >= 9 {
			bigNodesCounter++
		}
	}

	if counter != 20 {
		t.Error("missing nodes")
	}

	if bigNodesCounter < 10 {
		t.Error("no complete network for a start")
	}

}
