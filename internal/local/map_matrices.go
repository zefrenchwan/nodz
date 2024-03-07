package local

import (
	"errors"
)

// MapMatrix is a map implementation of a matrix, in memory
type MapMatrix[V any] struct {
	// size of the matrix.
	// Remember this is a local implementation that may run out of memory for large sizes
	size int
	// values contain the values for the indexes if set.
	// If not, we return the defaultValue
	values doubleMap[int, V]
	// defaultValue is the value to return if indexes are valid but have no value
	defaultValue V
	// empty is just the golang empty value for v (to return for an error)
	empty V
}

// NewMapMatrix creates a matrix with that size and defaultValue if no previous value was set
func NewMapMatrix[V any](size int, defaultValue V) MapMatrix[V] {
	var result MapMatrix[V]
	var empty V
	result.size = size
	result.empty = empty
	result.defaultValue = defaultValue
	result.values = newDoubleMap[int, V]()
	return result
}

// Size returns the size of the matrix
func (sm *MapMatrix[V]) Size() int {
	return sm.size
}

// SetValue sets the value at a given position
func (sm *MapMatrix[V]) SetValue(i, j int, value V) error {
	if i >= sm.size || i < 0 {
		return errors.New("invalid index")
	}

	if j >= sm.size || j < 0 {
		return errors.New("invalid index")
	}

	sm.values.putValue(i, j, value)
	return nil
}

// GetValue returns the value at a given position if set, otherwise the default value
func (sm *MapMatrix[V]) GetValue(i, j int) (V, bool, error) {
	if i >= sm.size || i < 0 {
		return sm.empty, false, errors.New("invalid index")
	}

	if j >= sm.size || j < 0 {
		return sm.empty, false, errors.New("invalid index")
	}

	result, has := sm.values.getValue(i, j)
	if !has {
		result = sm.defaultValue
	}

	return result, true, nil
}
