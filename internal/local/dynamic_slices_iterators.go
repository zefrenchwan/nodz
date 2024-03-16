package local

import "errors"

// DynamicSliceIterator is a slice based iterator.
// Plan is to add values at the head or tail of the slice
type DynamicSliceIterator[T any] struct {
	// values is the content of the iterator. We always read the first element if any
	values []T
	// opened is false for the first next, true otherwise
	opened bool
}

// NewDynamicSlicesIterator returns a new empty iterator
func NewDynamicSlicesIterator[T any]() DynamicSliceIterator[T] {
	var result DynamicSliceIterator[T]
	result.values = make([]T, 0)
	return result
}

// Next goes to the next element if any
func (dit *DynamicSliceIterator[T]) Next() (bool, error) {
	if dit == nil || dit.values == nil {
		return false, errors.New("nil iterator")
	}

	if len(dit.values) == 0 {
		return false, nil
	}

	if !dit.opened {
		dit.opened = true
	} else {
		dit.values = dit.values[1:]
	}

	return len(dit.values) != 0, nil
}

// Value returns current value if any
func (dit *DynamicSliceIterator[T]) Value() (T, error) {
	var empty T
	if dit == nil || dit.values == nil {
		return empty, errors.New("nil iterator")
	}

	if len(dit.values) == 0 {
		return empty, errors.New("empty iterator")
	}

	return dit.values[0], nil
}

// AddNextValue adds the next value to process
func (dit *DynamicSliceIterator[T]) AddNextValue(value T) error {
	if dit == nil {
		return errors.New("nil iterator")
	}

	dit.values = append(dit.values, value)
	copy(dit.values[1:], dit.values)
	dit.values[0] = value

	return nil
}

// AddLastValue adds the last element to process
func (dit *DynamicSliceIterator[T]) AddLastValue(value T) error {
	if dit == nil {
		return errors.New("nil iterator")
	}

	dit.values = append(dit.values, value)

	return nil
}

// Halt stops immediatly the iterator
func (dit *DynamicSliceIterator[T]) Halt() error {
	if dit == nil || dit.values == nil {
		return nil
	}

	dit.values = make([]T, 0)

	return nil
}
