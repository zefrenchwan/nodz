package internal_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestMapNoFilterIterator(t *testing.T) {
	baseIterator := local.NewSlicesIterator([]string{"aaaa", "bb", "c"})
	it := graphs.MapFilterIterator[string, int]{
		Iterator: &baseIterator,
		Mapper:   func(a string) (int, error) { return len(a), nil },
	}

	if ok, err := CompareIteratorWithSlice[int](&it, []int{4, 2, 1}, func(a, b int) bool { return a == b }, true); !ok || err != nil {
		t.Fail()
	}
}

func TestMapFilterIterator(t *testing.T) {
	baseIterator := local.NewSlicesIterator([]string{"aaaa", "bb", "c"})
	it := graphs.MapFilterIterator[string, int]{
		Iterator: &baseIterator,
		Mapper:   func(a string) (int, error) { return len(a), nil },
		Filter:   func(a int) bool { return a%2 != 0 },
	}

	if ok, err := CompareIteratorWithSlice[int](&it, []int{1}, func(a, b int) bool { return a == b }, true); !ok || err != nil {
		t.Fail()
	}

	// test the opposite condition
	baseIterator = local.NewSlicesIterator([]string{"aaaa", "bb", "c"})
	it = graphs.MapFilterIterator[string, int]{
		Iterator: &baseIterator,
		Mapper:   func(a string) (int, error) { return len(a), nil },
		Filter:   func(a int) bool { return a%2 == 0 },
	}

	if ok, err := CompareIteratorWithSlice[int](&it, []int{4, 2}, func(a, b int) bool { return a == b }, true); !ok || err != nil {
		t.Fail()
	}
}
