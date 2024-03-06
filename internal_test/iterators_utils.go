package internal_test

import (
	"slices"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// CompareIteratorWithSlice returns true if comparison was ok, false otherwise. If any error appears, return the error.
// Comparison for values may depend on an order.
// If order matters, set to true and compare values in the same order.
// Otherwise, slices are considered as sets.
func CompareIteratorWithSlice[N any](it graphs.GeneralIterator[N], slice []N, equals func(N, N) bool, orderMatters bool) (bool, error) {
	if it == nil {
		return len(slice) == 0, nil
	}

	// read values
	values := make([]N, 0)
	for has, err := it.Next(); has || err != nil; has, err = it.Next() {
		if err != nil {
			return false, err
		}

		if v, errV := it.Value(); errV != nil {
			return false, errV
		} else {
			values = append(values, v)
		}
	}

	if len(values) != len(slice) {
		return false, nil
	}

	// order test is same value at the same index
	if orderMatters {
		index := 0
		for {
			if index >= len(values) {
				break
			}

			if !equals(values[index], slice[index]) {
				return false, nil
			}

			index = index + 1
		}

		return true, nil
	}

	// no order is contains based
	for _, value := range values {
		localCompare := func(other N) bool {
			return equals(value, other)
		}

		if !slices.ContainsFunc(slice, localCompare) {
			return false, nil
		}
	}

	return true, nil
}
