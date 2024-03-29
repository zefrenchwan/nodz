package graphs

import "errors"

// GeneralIterator is an abstract definition of an iterator.
// Once created, an iterator is set before the first element, if any.
// It means that value should not respond if called before the Next function.
type GeneralIterator[T any] interface {
	// Next moves to the next element if any, and returns true if there is a next element.
	// It is initially set before first element, to deal with empty iterations.
	// For distant data sources, it also includes the ability to return an error
	// (for instance if said data source is not available)
	Next() (bool, error)
	// Value returns the current value if any, or an error.
	// Implementations may return an error of any kind (data could not be loaded)
	Value() (T, error)
}

// MapFilterIterator composes an iterator with a mapper and a filter.
// Mapper is not optional, but filter is (nil means no filter).
type MapFilterIterator[T any, S any] struct {
	// Iterator is the base iterator to map values for
	Iterator GeneralIterator[T]
	// Mapper is the function to pass from an instance of T to an instance of S
	Mapper func(T) (S, error)
	// Filter to exclude some S values
	Filter func(S) bool
}

// Next finds the next matching, if any
func (mi *MapFilterIterator[T, S]) Next() (bool, error) {
	if mi == nil {
		return false, errors.New("nil iterator")
	} else if mi.Iterator == nil {
		return false, nil
	}

	var globalErr error

	for {
		if has, err := mi.Iterator.Next(); err != nil {
			globalErr = errors.Join(globalErr, err)
			continue
		} else if !has {
			return false, globalErr
		} else if v, errV := mi.Iterator.Value(); errV != nil {
			globalErr = errors.Join(globalErr, errV)
			continue
		} else if mi.Filter == nil {
			return true, globalErr
		} else if res, errRes := mi.Mapper(v); errRes != nil {
			globalErr = errors.Join(globalErr, errRes)
			continue
		} else if mi.Filter(res) {
			return true, globalErr
		}
	}
}

// Value reads value from mi.Iterator and returns mapped value
func (mi *MapFilterIterator[T, S]) Value() (S, error) {
	var empty S

	if mi == nil || mi.Iterator == nil || mi.Mapper == nil {
		return empty, errors.New("nil iterator")
	}

	if value, err := mi.Iterator.Value(); err != nil {
		return empty, err
	} else {
		return mi.Mapper(value)
	}
}

// CompositeIterator is a composition of iterators.
// Once an iterator is complete, then move to the next one.
// Definition of "the next one" may vary over time
type CompositeIterator[T any] interface {
	// A composite iterator is an iterator too
	GeneralIterator[T]
	// ForceCurrent replaces current iterator with the parameter (then, current iterator is lost)
	ForceCurrent(GeneralIterator[T]) error
	// PostponeCurrent sets current iterator as the one to process right now, and the old current as the next one
	PostponeCurrent(GeneralIterator[T]) error
	// AddNext adds the pararemeter as the next iterator (once current is done)
	AddNext(GeneralIterator[T]) error
	// AddLast adds the parameter as the last iterator to process so far
	AddLast(GeneralIterator[T]) error
	// Halt immediatly stops the iterator : no more value, no more next
	Halt() error
}

// DynamicIterator is an iterator that allows to change values on the fly.
// Typical local implementation would be a double entries list (from head or tail)
type DynamicIterator[T any] interface {
	// A dynamic iterator is an iterator too
	GeneralIterator[T]
	// AddNextValue adds the parameter as the next value to read
	AddNextValue(T) error
	// AddLastValue adds the parameter as the last value to read
	AddLastValue(T) error
	// Halt stops iteration as soon as called
	Halt() error
}

// DynamicIteratorBuilder returns a new empty DynamicIterator
type DynamicIteratorBuilder[T any] func() (DynamicIterator[T], error)

// EmptyIterator is a commodity for an empty iterator of any type
type EmptyIterator[T any] struct{}

// Next returns false, and that is it (no move on empty iterator)
func (ei EmptyIterator[T]) Next() (bool, error) {
	return false, nil
}

// Value returns an error because an empty iterator has no content (by definition)
func (ei EmptyIterator[T]) Value() (T, error) {
	var defaultValue T
	return defaultValue, errors.New("no value for empty iterator")
}
