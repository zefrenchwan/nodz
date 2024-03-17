package local

import (
	"errors"
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// SlicesSet is a set based on slices
type SlicesSet[T any] struct {
	// equality returns true if two instances of T are the same (for set equality)
	equality graphs.SetEqualsFunction[T]
	// elements is a slice of non equal elements of the set
	elements []T
}

// NewSlicesSet returns a new empty slices set
func NewSlicesSet[T any](f func(T, T) bool) SlicesSet[T] {
	var result SlicesSet[T]
	result.equality = f
	result.elements = make([]T, 0)
	return result
}

// ToIterator returns an iterator over the elements of the set
func (s *SlicesSet[T]) ToIterator() (graphs.GeneralIterator[T], error) {
	var result SlicesIterator[T]
	if s == nil {
		result = NewSlicesIterator([]T{})
	} else {
		result = NewSlicesIterator(s.elements)
	}

	return &result, nil
}

// IsEmpty returns true for nil or empty set, false otherwise
func (s *SlicesSet[T]) IsEmpty() (bool, error) {
	return s == nil || len(s.elements) == 0, nil
}

// Add the element if not in the set
func (s *SlicesSet[T]) Add(element T) error {
	if s == nil {
		return errors.New("nil set")
	} else if s.equality == nil {
		return errors.New("nil test function")
	}

	if len(s.elements) == 0 {
		s.elements = []T{element}
		return nil
	}

	for _, elt := range s.elements {
		if s.equality(elt, element) {
			return nil
		}
	}

	s.elements = append(s.elements, element)

	return nil
}

// Has returns true if the set has an equal value in it
func (s *SlicesSet[T]) Has(element T) (bool, error) {
	if s == nil {
		return false, nil
	} else if s.equality == nil {
		return false, errors.New("nil test function")
	}

	return slices.ContainsFunc(s.elements, func(elt T) bool { return s.equality(elt, element) }), nil
}

// Remove excludes the first value equals (as defined in the set) to element
func (s *SlicesSet[T]) Remove(element T) error {
	if s == nil {
		return nil
	} else if s.equality == nil {
		return errors.New("nil test function")
	}

	s.elements = slices.DeleteFunc(s.elements, func(elt T) bool { return s.equality(elt, element) })
	return nil
}

// Peek gets an element in the set.
func (s *SlicesSet[T]) Peek() (T, error) {
	var empty T
	if s == nil || s.elements == nil {
		return empty, errors.New("nil set")
	}

	if len(s.elements) == 0 {
		return empty, errors.New("empty set")
	}

	return s.elements[0], nil
}

// Size returns the size of the set
func (s *SlicesSet[T]) Size() int64 {
	if s == nil {
		return 0
	}

	return int64(len(s.elements))
}
