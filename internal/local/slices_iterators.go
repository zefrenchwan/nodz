package local

import "errors"

// SlicesIterator is an in-memory implementation of iterators.
type SlicesIterator[T any] struct {
	// values to return during iteration
	Values []T
	// Index of the NEXT element to return
	Index int
}

// NewSlicesIterator returns an iterator over values
func NewSlicesIterator[T any](values []T) SlicesIterator[T] {
	return SlicesIterator[T]{
		Values: values,
		Index:  -1,
	}
}

// Next returns true for the next value if any, false otherwise.
// It also gets to the next element
func (i *SlicesIterator[T]) Next() (bool, error) {
	if i == nil || i.Values == nil {
		return false, nil
	}

	if i.Index <= len(i.Values) {
		i.Index = i.Index + 1
	}

	return i.Index < len(i.Values), nil
}

// Value returns the current value in the iterator, an error if there is no value
func (i *SlicesIterator[T]) Value() (T, error) {
	var defaultValue T
	if i == nil || i.Index < 0 || i.Index >= len(i.Values) {
		return defaultValue, errors.New("no value to return")
	}

	return i.Values[i.Index], nil
}
