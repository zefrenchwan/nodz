package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestSlicesIterator(t *testing.T) {
	values := []int{10, 20}
	it := local.NewSlicesIterator(values)

	var found bool
	var value int
	var errIt error

	if found, errIt = it.Next(); !found || errIt != nil {
		t.Fail()
	} else if value, errIt = it.Value(); value != values[0] || errIt != nil {
		t.Fail()
	}

	if found, errIt = it.Next(); !found || errIt != nil {
		t.Fail()
	} else if value, errIt = it.Value(); value != values[1] || errIt != nil {
		t.Fail()
	}

	if found, errIt = it.Next(); found || errIt != nil {
		t.Fail()
	}
}

func TestSlicesIteratorNoValue(t *testing.T) {
	values := make([]int, 0)
	it := local.NewSlicesIterator(values)

	if _, err := it.Value(); err == nil {
		t.Fail()
	}

	if v, err := it.Next(); err != nil || v {
		t.Fail()
	}
}
