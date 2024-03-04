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
