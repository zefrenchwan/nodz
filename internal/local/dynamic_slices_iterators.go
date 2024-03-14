package local

import (
	"errors"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// DynamicSlicesIterator is the local implementation of a dynamic iterator.
// Basically, it looks like a list of iterators to take in order
type DynamicSlicesIterator[T any] struct {
	// currentIterator is the source iterator to get data from
	currentIterator graphs.GeneralIterator[T]
	// nextIterators are the iterators to read once currentIterator is over
	nextIterators []graphs.GeneralIterator[T]
}

// NewDynamicSlicesIterator returns an iterator over parameter it, with the ability to add iterators after
func NewDynamicSlicesIterator[T any](it graphs.GeneralIterator[T]) DynamicSlicesIterator[T] {
	first := it
	if it == nil {
		empty := NewSlicesIterator([]T{})
		first = &empty
	}

	var result DynamicSlicesIterator[T]
	result.currentIterator = first
	result.nextIterators = make([]graphs.GeneralIterator[T], 0)
	return result
}

// Next moves to the next element, if any.
// This method implements iterator
func (it *DynamicSlicesIterator[T]) Next() (bool, error) {
	if it == nil || it.currentIterator == nil {
		return false, nil
	}

	// try to pick next from current iterator
	if has, err := it.currentIterator.Next(); err != nil {
		return false, err
	} else if has {
		return true, nil
	}

	var globalErr error
	// find first valid iterator and let it be the new current iterator
	firstValidIteratorIndex := -1
	for index, iterator := range it.nextIterators {
		if iterator == nil {
			continue
		}

		if has, err := iterator.Next(); err != nil {
			globalErr = errors.Join(globalErr, err)
		} else if has {
			firstValidIteratorIndex = index
			break
		}
	}

	// attention: globalErr may not be nil, but a valid iterator may exist
	switch firstValidIteratorIndex {
	case -1:
		// no matching iterator
		return false, globalErr
	case 0:
		// first element, just leave it out
		it.currentIterator = it.nextIterators[firstValidIteratorIndex]
		it.nextIterators = it.nextIterators[1:]
	case len(it.nextIterators) - 1:
		// last one matches, clean nextIterators
		it.currentIterator = it.nextIterators[firstValidIteratorIndex]
		it.nextIterators = nil
	default:
		// matching iterator, remaining iterators
		it.currentIterator = it.nextIterators[firstValidIteratorIndex]
		it.nextIterators = it.nextIterators[firstValidIteratorIndex+1:]
	}

	return true, globalErr
}

// Value returns the current value, if any.
// This method implements iterator
func (it *DynamicSlicesIterator[T]) Value() (T, error) {
	var defaultValue T
	if it == nil || it.currentIterator == nil {
		return defaultValue, errors.New("empty iterator, no value")
	}

	return it.currentIterator.Value()
}

// ForceCurrent forces newIterator as the current iterator
func (it *DynamicSlicesIterator[T]) ForceCurrent(newIterator graphs.GeneralIterator[T]) error {
	if it == nil || newIterator == nil {
		return errors.New("empty iterator")
	}

	it.currentIterator = newIterator
	return nil
}

// PostponeCurrent passes current iterator as the next one, and runs newIterator first.
func (it *DynamicSlicesIterator[T]) PostponeCurrent(newIterator graphs.GeneralIterator[T]) error {
	if it == nil || newIterator == nil {
		return errors.New("empty iterator")
	}

	it.nextIterators = append(it.nextIterators, nil)
	copy(it.nextIterators[1:], it.nextIterators)
	it.nextIterators[0] = it.currentIterator
	it.currentIterator = newIterator
	return nil
}

// AddNext adds newIterator as the next one to run
func (it *DynamicSlicesIterator[T]) AddNext(newIterator graphs.GeneralIterator[T]) error {
	if it == nil || newIterator == nil {
		return errors.New("empty iterator")
	}

	it.nextIterators = append(it.nextIterators, nil)
	copy(it.nextIterators[1:], it.nextIterators)
	it.nextIterators[0] = newIterator
	return nil
}

// AddLast adds newIterator as the last iterator to run
func (it *DynamicSlicesIterator[T]) AddLast(newIterator graphs.GeneralIterator[T]) error {
	if it == nil || newIterator == nil {
		return errors.New("empty iterator")
	}

	it.nextIterators = append(it.nextIterators, newIterator)
	return nil
}
