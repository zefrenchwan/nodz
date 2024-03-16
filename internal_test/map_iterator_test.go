package internal_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestMapIterator(t *testing.T) {
	baseIterator := local.NewSlicesIterator([]string{"aaaa", "bb", "c"})
	it := graphs.MapIterator[string, int]{
		Iterator: &baseIterator,
		Mapper:   func(a string) int { return len(a) },
	}

	if ok, err := CompareIteratorWithSlice[int](&it, []int{4, 2, 1}, func(a, b int) bool { return a == b }, true); !ok || err != nil {
		t.Fail()
	}
}
