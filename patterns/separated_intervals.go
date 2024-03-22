package patterns

// SeparatedIntervals is the union of separated intervals.
// None of them is empty.
type SeparatedIntervals[T any] struct {
	// elements are the interval.
	// Invariant is: none empty, if (a,b) in elements then a and b are separated
	elements []Interval[T]
}
