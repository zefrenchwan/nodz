package local_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/internal/local"
	"github.com/zefrenchwan/nodz.git/internal_test"
)

func TestMapMatrixt(t *testing.T) {
	matrix, errMatrix := local.NewMapMatrix(3, 1)
	if errMatrix != nil {
		t.Fail()
	}

	matrix.SetValue(0, 2, 20)
	matrix.SetValue(1, 1, 10)
	matrix.SetValue(2, 0, 20)

	if v, f, err := matrix.GetValue(0, 0); v != 1 || err != nil || !f {
		t.Error("default value failing")
	}

	if v, f, err := matrix.GetValue(0, 2); v != 20 || err != nil || !f {
		t.Error("set value failing")
	}

	it, errIt := matrix.Line(1)
	if errIt != nil {
		t.Fail()
	}

	res, errRes := internal_test.CompareIteratorWithSlice(
		it,
		[]int{1, 10, 1},
		func(a, b int) bool { return a == b },
		true,
	)
	if !res || errRes != nil {
		t.Error("line failure")
	}

	it, errIt = matrix.Column(2)
	if errIt != nil {
		t.Fail()
	}

	res, errRes = internal_test.CompareIteratorWithSlice(
		it,
		[]int{20, 1, 1},
		func(a, b int) bool { return a == b },
		true,
	)
	if !res || errRes != nil {
		t.Error("line failure")
	}
}
