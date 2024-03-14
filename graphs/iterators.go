package graphs

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

// DynamicIterator is a composition of iterators.
// Typical use would be a graph walkthrough: neighborhoods are added as we run into a graph.
// But, in general, we may add iterators on the fly:
// * as the next iterator to run (depth first walkthrough)
// * as the last iterator (breadth first walkthrough)
type DynamicIterator[T any] interface {
	// A dynamic iterator is an iterator too
	GeneralIterator[T]
	// ForceCurrent replaces current iterator with the parameter (then, current iterator is lost)
	ForceCurrent(GeneralIterator[T]) error
	// PostponeCurrent sets current iterator as the one to process right now, and the old current as the next one
	PostponeCurrent(GeneralIterator[T]) error
	// AddNext adds the pararemeter as the next iterator (once current is done)
	AddNext(GeneralIterator[T]) error
	// AddLast adds the parameter as the last iterator to process so far
	AddLast(GeneralIterator[T]) error
}
