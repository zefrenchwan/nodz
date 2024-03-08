package local

// doubleMap is a map of maps, useful to link two instances of K with a value
type doubleMap[K comparable, V any] map[K]map[K]V

// newDoubleMap constructs a new empty double map
func newDoubleMap[K comparable, V any]() doubleMap[K, V] {
	return make(map[K]map[K]V)
}

// putValue sets the value for first key a and then second key b
func (dm doubleMap[K, V]) putValue(a, b K, value V) {
	if dm == nil {
		return
	}

	rows, found := dm[a]
	if !found || rows == nil {
		dm[a] = make(map[K]V)
	}

	dm[a][b] = value
}

// putElement adds an empty map for given first key a
func (dm doubleMap[K, V]) putElement(a K) {
	rows, found := dm[a]
	if !found || rows == nil {
		dm[a] = make(map[K]V)
	}
}

// getElement returns the map for first key a
func (dm doubleMap[K, V]) getElement(a K) map[K]V {
	if dm == nil {
		return nil
	}

	return dm[a]
}

// getValue for keys a and then b returns the value, if any, set for those keys.
// If no value was set, returns default value and false
func (dm doubleMap[K, V]) getValue(a, b K) (V, bool) {
	var result V
	if dm == nil {
		return result, false
	}

	rows, found := dm[a]
	if !found || rows == nil {
		return result, false
	}

	result, found = dm[a][b]
	return result, found
}

// removeElement removes all instances of a in secondary maps and primary map
func (dm doubleMap[K, V]) removeElement(a K) {
	if dm == nil {
		return
	}

	for _, v := range dm {
		if v != nil {
			delete(v, a)
		}
	}

	delete(dm, a)
}

// removeValue removes value for first key a and second key b
func (dm doubleMap[K, V]) removeValue(a, b K) {
	if dm == nil {
		return
	}

	rows, found := dm[a]
	if found && rows != nil {
		delete(rows, b)
	}
}

// getElementsLinkedToSecondaryKey is the map of primary keys and secondary values, with b as secondary keys
func (dm doubleMap[K, V]) getElementsLinkedToSecondaryKey(b K) map[K]V {
	result := make(map[K]V)
	for index, secValues := range dm {
		secValue, found := secValues[b]
		if found {
			result[index] = secValue
		}
	}

	return result
}
