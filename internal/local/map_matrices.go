package local

import (
	"errors"

	"github.com/zefrenchwan/nodz.git/graphs"
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
}

// NewMapMatrix creates a matrix with that size and defaultValue if no previous value was set
func NewMapMatrix[V any](size int, defaultValue V) (MapMatrix[V], error) {
	var result MapMatrix[V]
	if size <= 0 {
		return result, errors.New("invalid matric size")
	}

	result.size = size
	result.defaultValue = defaultValue
	result.values = newDoubleMap[int, V]()
	return result, nil
}

// Size returns the size of the matrix
func (sm *MapMatrix[V]) Size() int {
	if sm == nil {
		return 0
	}

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
	var empty V
	if i >= sm.size || i < 0 {
		return empty, false, errors.New("invalid index")
	}

	if j >= sm.size || j < 0 {
		return empty, false, errors.New("invalid index")
	}

	result, has := sm.values.getValue(i, j)
	if !has {
		result = sm.defaultValue
	}

	return result, true, nil
}

// Line returns line i as an iterator
func (sm *MapMatrix[V]) Line(i int) (graphs.GeneralIterator[V], error) {
	if sm == nil {
		return nil, nil
	}

	if i < 0 || i >= sm.size {
		return nil, errors.New("invalid index")
	}

	line := sm.values.getElement(i)
	it, errIt := NewMapIterator(sm.size-1, line, sm.defaultValue)

	return &it, errIt
}

// Column returns column j as an iterator
func (sm *MapMatrix[V]) Column(j int) (graphs.GeneralIterator[V], error) {
	if sm == nil {
		return nil, nil
	}

	if j < 0 || j >= sm.size {
		return nil, errors.New("invalid index")
	}

	column := sm.values.getElementsLinkedToSecondaryKey(j)
	it, errIt := NewMapIterator(sm.size-1, column, sm.defaultValue)

	return &it, errIt
}
