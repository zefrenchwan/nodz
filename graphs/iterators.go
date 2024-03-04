package graphs

import "errors"

// GeneralIterator is an abstract definition of an iterator.
// Once created, an iterator is set before the first element, if any.
// It means that value should not respond if called before the Next function.
type GeneralIterator[T any] interface {
	// Next moves to the next element if any, and returns true if there is a next element.
	// It is initially set before first element, to deal with empty iterations.
	// For distant data sources, it also includes the ability to return an error
	// (for instance if said data source is not available)
	Next() (bool, error)
	// Value returns the current value if any, or an error.
	// Implementations may return an error of any kind (data could not be loaded)
	Value() (T, error)
}

// LocalIterator is an in-memory implementation of iterators.
type LocalIterator[T any] struct {
	// values to return during iteration
	Values []T
	// Index of the NEXT element to return
	Index int
}

// NewLocalIterator returns an iterator over values
func NewLocalIterator[T any](values []T) LocalIterator[T] {
	return LocalIterator[T]{
		Values: values,
		Index:  -1,
	}
}

// Next returns true for the next value if any, false otherwise.
// It also gets to the next element
func (i *LocalIterator[T]) Next() (bool, error) {
	if i == nil || i.Values == nil {
		return false, nil
	}

	if i.Index <= len(i.Values) {
		i.Index = i.Index + 1
	}

	return i.Index < len(i.Values), nil
}

// Value returns the current value in the iterator, an error if there is no value
func (i *LocalIterator[T]) Value() (T, error) {
	var defaultValue T
	if i == nil || i.Index < 0 || i.Index >= len(i.Values) {
		return defaultValue, errors.New("no value to return")
	}

	return i.Values[i.Index], nil
}
