package internal_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal"
)

func TestLocalIteratorNoValue(t *testing.T) {
	values := make([]int, 0)
	it := internal.NewLocalIterator(values)

	if _, err := it.Value(); err == nil {
		t.Fail()
	}

	if v, err := it.Next(); err != nil || v {
		t.Fail()
	}
}

func TestLocalIteratorValues(t *testing.T) {
	values := []int{10, 20}
	it := internal.NewLocalIterator(values)

	if _, err := it.Value(); err == nil {
		t.Fail()
	}

	if v, err := it.Next(); err != nil || !v {
		t.Fail()
	}

	if v, err := it.Value(); v != 10 || err != nil {
		t.Fail()
	}

	if v, err := it.Next(); err != nil || !v {
		t.Fail()
	}

	if v, err := it.Value(); v != 20 || err != nil {
		t.Fail()
	}

	if v, err := it.Next(); err != nil || v {
		t.Fail()
	}
}
