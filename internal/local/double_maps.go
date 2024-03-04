package local

type doubleMap[K comparable, V any] map[K]map[K]V

func newDoubleMap[K comparable, V any]() doubleMap[K, V] {
	return make(map[K]map[K]V)
}

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

func (dm doubleMap[K, V]) putElement(a K) {
	rows, found := dm[a]
	if !found || rows == nil {
		dm[a] = make(map[K]V)
	}
}

func (dm doubleMap[K, V]) getElement(a K) map[K]V {
	if dm == nil {
		return nil
	}

	return dm[a]
}

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

func (dm doubleMap[K, V]) removeValue(a, b K) {
	if dm == nil {
		return
	}

	rows, found := dm[a]
	if found && rows != nil {
		delete(rows, b)
	}
}
