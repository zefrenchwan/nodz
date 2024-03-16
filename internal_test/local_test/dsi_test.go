package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestDynamicSliceIterator(t *testing.T) {
	it := local.NewDynamicSlicesIterator[int]()
	it.AddNextValue(2)
	it.AddNextValue(1)
	it.AddLastValue(3)

	if res, err := internal_test.CompareIteratorWithSlice(&it, []int{1, 2, 3}, func(a, b int) bool { return a == b }, true); err != nil || !res {
		t.Fail()
	}
}

func TestDynamicSliceIteratorEmpty(t *testing.T) {
	it := local.NewDynamicSlicesIterator[int]()
	if has, err := it.Next(); has || err != nil {
		t.Fail()
	}
}

func TestDynamicSliceIteratorHalt(t *testing.T) {
	it := local.NewDynamicSlicesIterator[int]()

	it.AddNextValue(2)
	it.AddNextValue(1)
	it.AddLastValue(3)

	it.Halt()
	if has, err := it.Next(); has || err != nil {
		t.Fail()
	}
}
