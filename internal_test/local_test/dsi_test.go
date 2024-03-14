package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestDynamicSlicesIterators(t *testing.T) {
	firstElements := local.NewSlicesIterator([]int{1, 2, 3})
	secondElements := local.NewSlicesIterator([]int{4, 5, 6})
	replacedIterator := local.NewSlicesIterator([]int{7})
	emptyElements := local.NewSlicesIterator([]int{})
	thirdElements := local.NewSlicesIterator([]int{8, 9})

	it := local.NewDynamicSlicesIterator(&replacedIterator)
	it.ForceCurrent(&secondElements)
	it.PostponeCurrent(&firstElements)
	it.AddNext(&emptyElements)
	it.AddLast(&thirdElements)

	if v, err := internal_test.CompareIteratorWithSlice(&it, []int{1, 2, 3, 4, 5, 6, 8, 9}, func(a, b int) bool { return a == b }, true); err != nil {
		t.Fail()
	} else if !v {
		t.Error("it failure for test")
	}
}
