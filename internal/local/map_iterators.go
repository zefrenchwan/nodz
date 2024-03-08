package local

import "errors"

// MapIterator walks through a map using increasing indexes from 0 to a max,
// and returns default value if no value was in the map.
// It is a sort of iterator from a map dealing with non existing values.
type MapIterator[V any] struct {
	// minIndex (0)
	minIndex int
	// index is the current index of the value to find in the map, if any
	index int
	// maxIndex is the max index to iterate through
	maxIndex int
	// values is the core map to find values in
	values map[int]V
	// defaultValue is the value to return if no value was in the map
	defaultValue V
}

// NewMapIterator returns a new map iterator over a given map, until a maxIndex (included)
func NewMapIterator[V any](maxIndex int, originalValues map[int]V, defaultValue V) (MapIterator[V], error) {
	var result MapIterator[V]

	minIndex := 0
	if minIndex > maxIndex {
		return result, errors.New("minIndex should be less than maxIndex")
	}

	result.minIndex = minIndex
	result.maxIndex = maxIndex
	result.index = minIndex - 1
	result.values = originalValues
	result.defaultValue = defaultValue

	return result, nil
}

// Next moves to the next element if any (and return true), or returns false
func (mi *MapIterator[V]) Next() (bool, error) {
	if mi == nil || mi.index >= mi.maxIndex {
		return false, nil
	}

	mi.index = mi.index + 1
	return true, nil
}

// Value returns current value if any, error otherwise
func (mi *MapIterator[V]) Value() (V, error) {
	var empty V
	if mi == nil || mi.index > mi.maxIndex {
		return empty, errors.New("no value from iterator")
	}

	if mi.values == nil {
		return mi.defaultValue, nil
	}

	result, found := mi.values[mi.index]
	if !found {
		result = mi.defaultValue
	}

	return result, nil
}
