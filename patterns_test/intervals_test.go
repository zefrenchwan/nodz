package patterns_test

import (
	"testing"

	"github.com/zefrenchwan/nodz.git/patterns"
)

func TestIntervalsCompare(t *testing.T) {
	comparator := patterns.NewIntComparator()
	var a, b patterns.Interval[int]
	// empty test
	a = comparator.NewEmptyInterval()
	b = comparator.NewEmptyInterval()

	if !a.IsEmpty() {
		t.Fail()
	}

	if comparator.CompareInterval(a, b) != 0 {
		t.Error("failed empty equals empty")
	}

	a = comparator.NewFullInterval()
	if comparator.CompareInterval(a, b) > 0 {
		t.Error("failed empty is less than anything")
	}

	// full is more than anything
	a = comparator.NewLeftInfiniteInterval(10, false)
	b = comparator.NewFullInterval()

	if a.IsFull() || !b.IsFull() {
		t.Fail()
	}

	if comparator.CompareInterval(b, b) != 0 {
		t.Error("failed test on fulll is full")
	} else if comparator.CompareInterval(b, a) > 0 {
		t.Error("failed full is more than anything")
	} else if comparator.CompareInterval(a, b) < 0 {
		t.Error("failed full is more than anything")
	}

	// test cases on left infinite
	a = comparator.NewLeftInfiniteInterval(5, false)
	b = comparator.NewLeftInfiniteInterval(5, true)

	if a.IsEmpty() || b.IsEmpty() || a.IsFull() || b.IsFull() {
		t.Fail()
	}

	if comparator.CompareInterval(a, a) != 0 {
		t.Error("test equality failure")
	} else if comparator.CompareInterval(a, b) >= 0 {
		t.Error("failed test on left infinite: check right comparison")
	} else if comparator.CompareInterval(b, a) <= 0 {
		t.Error("failed test on left infinite: check right comparison")
	}

	a = comparator.NewLeftInfiniteInterval(2, true)
	b = comparator.NewLeftInfiniteInterval(5, true)
	if comparator.CompareInterval(a, b) >= 0 {
		t.Error("failed test on left infinite: check right value comparison")
	} else if comparator.CompareInterval(b, a) <= 0 {
		t.Error("failed test on left infinite: check right value comparison")
	}

	// test cases on right infinite
	a = comparator.NewRightInfiniteInterval(10, false)
	b = comparator.NewRightInfiniteInterval(10, true)

	if a.IsEmpty() || b.IsEmpty() || a.IsFull() || b.IsFull() {
		t.Fail()
	}

	if comparator.CompareInterval(a, a) != 0 {
		t.Error("equality failure")
	} else if comparator.CompareInterval(b, b) != 0 {
		t.Error("equality failure")
	} else if comparator.CompareInterval(a, b) >= 0 {
		t.Error("check left comparison")
	} else if comparator.CompareInterval(b, a) <= 0 {
		t.Error("check left comparison")
	}

	a = comparator.NewRightInfiniteInterval(1, true)
	if comparator.CompareInterval(a, b) >= 0 {
		t.Error("check left comparison")
	} else if comparator.CompareInterval(b, a) <= 0 {
		t.Error("check left comparison")
	}

	// test finite impossible intervals
	if _, err := comparator.NewFiniteInterval(10, 2, false, false); err == nil {
		t.Fail()
	} else if _, err := comparator.NewFiniteInterval(10, 10, false, true); err == nil {
		t.Fail()
	} else if _, err := comparator.NewFiniteInterval(10, 10, false, true); err == nil {
		t.Fail()
	}

	// test many combinations
	a = comparator.NewLeftInfiniteInterval(1, false)
	b = comparator.NewRightInfiniteInterval(10, true)
	if comparator.CompareInterval(a, b) >= 0 {
		t.Error("mixed failure for ]-oo, 1[ and [10, +oo[")
	} else if comparator.CompareInterval(b, a) <= 0 {
		t.Error("mixed failure for ]-oo, 1[ and [10, +oo[")
	}
}

func TestIntervalComplement(t *testing.T) {
	comparator := patterns.NewIntComparator()
	var a patterns.Interval[int]
	// empty test
	a = comparator.NewEmptyInterval()
	result := comparator.Complement(a)
	if len(result) != 1 || !result[0].IsFull() {
		t.Error("complement of empty should be full")
	}

	// full test
	a = comparator.NewFullInterval()
	result = comparator.Complement(a)
	if len(result) != 1 || !result[0].IsEmpty() {
		t.Error("complement of full should be empty")
	}

	// semi bounded intervals
	a = comparator.NewLeftInfiniteInterval(10, false)
	expected := comparator.NewRightInfiniteInterval(10, true)
	result = comparator.Complement(a)
	if len(result) != 1 || comparator.CompareInterval(expected, result[0]) != 0 {
		t.Error("complement failure for semi bounded intervals")
	}

	a = comparator.NewRightInfiniteInterval(1, false)
	expected = comparator.NewLeftInfiniteInterval(1, true)
	result = comparator.Complement(a)
	if len(result) != 1 || comparator.CompareInterval(expected, result[0]) != 0 {
		t.Error("complement failure for semi bounded intervals")
	}

	// bounded intervals
	a, errInterval := comparator.NewFiniteInterval(1, 10, true, false)
	if errInterval != nil {
		t.Fail()
	}

	result = comparator.Complement(a)
	if len(result) != 2 {
		t.Error("complement failure for semi bounded intervals")
	}

	exp1 := result[0]
	exp2 := result[1]

	if comparator.CompareInterval(exp1, exp2) > 0 {
		exp1, exp2 = exp2, exp1
	}

	if comparator.CompareInterval(exp1, comparator.NewLeftInfiniteInterval(1, false)) != 0 {
		t.Error("complement failure for semi bounded intervals")
	} else if comparator.CompareInterval(exp2, comparator.NewRightInfiniteInterval(10, true)) != 0 {
		t.Error("complement failure for semi bounded intervals")
	}
}

func TestIntervalIntersection(t *testing.T) {
	comparator := patterns.NewIntComparator()
	var a, b, result, expected patterns.Interval[int]

	// test empty and full
	a = comparator.NewFullInterval()
	b = comparator.NewEmptyInterval()
	result = comparator.Intersection(a, b)
	if !result.IsEmpty() {
		t.Fail()
	}

	result = comparator.Intersection(b, a)
	if !result.IsEmpty() {
		t.Fail()
	}

	result = comparator.Intersection(a, a)
	if !result.IsFull() {
		t.Fail()
	}

	// test semi bounded
	a = comparator.NewLeftInfiniteInterval(10, true)
	b, _ = comparator.NewFiniteInterval(0, 20, true, false)
	expected, _ = comparator.NewFiniteInterval(0, 10, true, true)
	result = comparator.Intersection(a, b)
	if comparator.CompareInterval(expected, result) != 0 {
		t.Fail()
	}

	result = comparator.Intersection(b, a)
	if comparator.CompareInterval(expected, result) != 0 {
		t.Fail()
	}

	a = comparator.NewLeftInfiniteInterval(10, true)
	b = comparator.NewLeftInfiniteInterval(50, false)
	expected = comparator.NewLeftInfiniteInterval(10, true)
	result = comparator.Intersection(a, b)
	if comparator.CompareInterval(expected, result) != 0 {
		t.Fail()
	}

	a = comparator.NewLeftInfiniteInterval(10, true)
	b = comparator.NewRightInfiniteInterval(0, false)
	expected, _ = comparator.NewFiniteInterval(0, 10, false, true)
	result = comparator.Intersection(a, b)
	if comparator.CompareInterval(expected, result) != 0 {
		t.Fail()
	}

	// test bounded
	a, _ = comparator.NewFiniteInterval(0, 5, true, false)
	b, _ = comparator.NewFiniteInterval(100, 105, true, false)
	result = comparator.Intersection(a, b)
	if !result.IsEmpty() {
		t.Fail()
	}

	a, _ = comparator.NewFiniteInterval(0, 102, true, false)
	b, _ = comparator.NewFiniteInterval(100, 105, true, false)
	result = comparator.Intersection(a, b)
	expected, _ = comparator.NewFiniteInterval(100, 102, true, false)
	if comparator.CompareInterval(result, expected) != 0 {
		t.Fail()
	}
}
