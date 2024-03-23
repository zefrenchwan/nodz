package patterns

// SeparatedIntervals is the union of separated intervals.
// None of them is empty.
type SeparatedIntervals[T any] struct {
	// elements are the interval.
	// Invariant is: none empty, if (a,b) in elements then a and b are separated
	elements []Interval[T]
}

// NewSeparatedIntervals returns a new separated intervals with initial value
func NewSeparatedIntervals[T any](initialValue Interval[T]) SeparatedIntervals[T] {
	var result SeparatedIntervals[T]
	result.elements = make([]Interval[T], 0)
	if !initialValue.IsEmpty() {
		result.elements = append(result.elements, initialValue)
	}

	return result
}
