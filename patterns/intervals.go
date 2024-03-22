package patterns

import "errors"

// TypedComparator defines a compare function over a type
type TypedComparator[T any] struct {
	comparator func(T, T) int
}

// NewTypedComparator returns an interval manager based on a compare function.
// Contract for compareFn(a, b) is:
// * if a < b, return a negative value
// * if a > b, return a positive value
// * if a == b, return 0
// * no error returned
// * comparison should be quick
func NewTypedComparator[T any](compareFn func(T, T) int) TypedComparator[T] {
	var result TypedComparator[T]
	result.comparator = compareFn
	return result
}

// Compare decorates the comparator function
func (t TypedComparator[T]) Compare(a, b T) int {
	return t.comparator(a, b)
}

// Min returns the min of values
func (t TypedComparator[T]) Min(a T, v ...T) T {
	min := a
	for _, other := range v {
		if t.Compare(other, min) < 0 {
			min = other
		}
	}

	return min
}

// Interval is the definition of intervals based on the comparator.
type Interval[T any] struct {
	// true for empty interval
	empty bool
	// true if interval is not left bounded
	minInfinite bool
	// min of the interval, if not minInfinite
	min T
	// if not minInfinite, true if the min is the interval, false otherwise
	minIncluded bool
	// true if interval is not right bounded, false otherwise
	maxInfinite bool
	// max of the interval, if not maxInfinite
	max T
	// if not maxInfinite, true if the max is the interval, false otherwise
	maxIncluded bool
}

// IsFull returns true for an unbounded interval
func (i Interval[T]) IsFull() bool {
	return i.maxInfinite && i.minInfinite
}

// IsEmpty is true for an empty interval, false otherwise
func (i Interval[T]) IsEmpty() bool {
	return i.empty
}

// NewEmptyInterval returns a new empty interval
func (t TypedComparator[T]) NewEmptyInterval() Interval[T] {
	var result Interval[T]
	result.empty = true
	return result
}

// NewFullInterval returns a full interval
func (t TypedComparator[T]) NewFullInterval() Interval[T] {
	var result Interval[T]
	result.minInfinite = true
	result.maxInfinite = true
	return result
}

// NewLeftInfiniteInterval returns an interval ] -oo, rightValue )
func (t TypedComparator[T]) NewLeftInfiniteInterval(rightValue T, rightIncluded bool) Interval[T] {
	var result Interval[T]
	result.max = rightValue
	result.maxIncluded = rightIncluded
	result.minInfinite = true
	return result
}

// NewRightInfiniteInterval returns an interval ( leftValue, +oo [
func (t TypedComparator[T]) NewRightInfiniteInterval(leftValue T, leftIncluded bool) Interval[T] {
	var result Interval[T]
	result.min = leftValue
	result.minIncluded = leftIncluded
	result.maxInfinite = true
	return result
}

// NewFiniteInterval returns a finite interval or an error if interval would be empty
func (t TypedComparator[T]) NewFiniteInterval(left, right T, leftIn, rightIn bool) (Interval[T], error) {
	var result Interval[T]
	comparison := t.Compare(left, right)
	if comparison > 0 || (comparison == 0 && !(leftIn && rightIn)) {
		return result, errors.New("interval parameters would make empty interval")
	}

	result.max = right
	result.min = left
	result.minIncluded = leftIn
	result.maxIncluded = rightIn

	return result, nil
}

// CompareInterval is an order based on the lexicographic order.
// Same sets are equals (return 0).
func (t TypedComparator[T]) CompareInterval(a, b Interval[T]) int {
	// deal with empty or full intervals
	switch {
	case a.empty:
		if b.empty {
			return 0
		}

		return 1
	case b.empty:
		return -1
	case a.maxInfinite && a.minInfinite:
		if b.minInfinite && b.maxInfinite {
			return 0
		}

		return -1
	}

	// deal with left boundaries
	switch {
	case a.minInfinite && !b.minInfinite:
		return -1
	case b.minInfinite && !a.minInfinite:
		return 1
	case !a.minInfinite && !b.minInfinite:
		minCompare := t.Compare(a.min, b.min)
		if minCompare != 0 {
			return minCompare
		} else if a.minIncluded && !b.minIncluded {
			return 1
		} else if !a.minIncluded && b.minIncluded {
			return -1
		}
	}

	// left boundaries are equals, result is now based on right boundaries
	switch {
	case a.maxInfinite && b.maxInfinite:
		return 0
	case a.maxInfinite:
		return -1
	case b.maxInfinite:
		return 1
	}

	comparison := t.Compare(a.max, b.max)
	if comparison != 0 {
		return comparison
	} else if a.maxIncluded == b.maxIncluded {
		return 0
	} else if a.maxIncluded {
		return 1
	} else {
		return -1
	}

}

// Complement returns the complement of the interval.
// It may be a single set (empty => full, full => empty, etc) or two (for finite intervals)
func (t TypedComparator[T]) Complement(i Interval[T]) []Interval[T] {
	var result Interval[T]
	switch {
	case i.empty:
		result.maxInfinite = true
		result.minInfinite = true
		return []Interval[T]{result}
	case i.minInfinite && i.maxInfinite:
		result.empty = true
		return []Interval[T]{result}
	case i.minInfinite:
		result.maxInfinite = true
		result.min = i.max
		result.minIncluded = !result.maxIncluded
		return []Interval[T]{result}
	case i.maxInfinite:
		result.minInfinite = true
		result.max = i.min
		result.maxIncluded = !i.minIncluded
		return []Interval[T]{result}
	}

	// remaining case is (a,b) with finite values
	// Then, result is ]-oo, a( and )b, +oo[
	var otherResult Interval[T]
	result.minInfinite = true
	result.max = i.min
	result.maxIncluded = !i.minIncluded
	otherResult.maxInfinite = true
	otherResult.minIncluded = !i.maxIncluded
	otherResult.min = i.max

	return []Interval[T]{result, otherResult}
}
