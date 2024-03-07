package local

import (
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
)

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

// size returns the size of the mapping
func (im *increasingMapping[V]) size() int {
	return len(im.values)
}

// toIncreasingIndexes returns the sorted indexes of the values as a map.
// Key is the source index (for instance 10) and value is increasing int from 0 to size - 1
// Basically, it is the only increasing mapping from values to {0, size -1}
func (im *increasingMapping[V]) toIncreasingIndexes() map[int]int {
	if im == nil || im.values == nil {
		return nil
	}

	indexes := make([]int, len(im.values))
	index := 0
	for k := range im.values {
		indexes[index] = k
		index++
	}

	slices.Sort(indexes)

	result := make(map[int]int)
	for index, value := range indexes {
		result[value] = index
	}

	return result
}

// toIncreasingValues returns the values in a slice, forming an increasing mapping from im.values.
// It means that if a,A and b,B are in im.values and a < b, then A is before B in the result.
func (im *increasingMapping[V]) toIncreasingValues() []V {
	if im == nil || im.values == nil {
		return nil
	}

	// To do so, we get the indexes of the elements,
	// we sort them to get the final order, so we have indexes in the right order.
	// And then we get elements in the same order as the sorted indexes
	indexes := im.toIncreasingIndexes()

	result := make([]V, len(indexes))
	for i, k := range indexes {
		result[i] = im.values[k]
	}

	return result
}
