package local

import "github.com/zefrenchwan/nodz.git/graphs"

// increasingMapping is a mapping from V to int with increasing values
type increasingMapping[V any] struct {
	// maxIndex is the index of the next element to insert
	maxIndex int
	// values is the mapping per se, but V is not comparable, so we use int as keys
	values map[int]V
	// equals tests if two elements are the same
	equals func(V, V) bool
}

// newIncreasingMapping returns a new empty mapping. Elements are compared with equals
func newIncreasingMapping[V any](equalsFn func(V, V) bool) increasingMapping[V] {
	var result increasingMapping[V]
	result.values = make(map[int]V)
	result.equals = equalsFn
	result.maxIndex = 0
	return result
}

// addValue adds a value if not already there, and returns the index of the value
func (im *increasingMapping[V]) addValue(value V) int {
	for k, v := range im.values {
		if im.equals(v, value) {
			return k
		}
	}

	index := im.maxIndex
	im.values[index] = value
	im.maxIndex = index + 1

	return index
}

// getValue returns the index of the element if found, 0 and false
func (im *increasingMapping[V]) getValue(value V) (int, bool) {
	for k, v := range im.values {
		if im.equals(v, value) {
			return k, true
		}
	}

	return 0, false
}

// removeValue removes the value if any, it does not affect the mapping of the other elements
func (im *increasingMapping[V]) removeValue(value V) {
	index, found := im.getValue(value)
	if found {
		delete(im.values, index)
	}
}

// toIterator returns an iterator over the values of V
func (im *increasingMapping[V]) toIterator() graphs.GeneralIterator[V] {
	values := make([]V, 0)
	for _, v := range im.values {
		values = append(values, v)
	}

	it := NewSlicesIterator(values)
	return &it
}
