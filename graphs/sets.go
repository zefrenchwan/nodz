package graphs

// AbstractSet defines what a generic set (distributed maybe) should do
type AbstractSet[T any] interface {
	// ToIterator returns a new iterator over the elements of the set
	ToIterator() (GeneralIterator[T], error)
	// IsEmpty returns true for an empty set, false otherwise
	IsEmpty() (bool, error)
	// Add a value in the set
	Add(T) error
	// Has returns true if the value is in the set
	Has(T) (bool, error)
	// Remove an element
	Remove(T) error
	// Peek return an element if the set is not empty, raises an error if set if nil or empty.
	// It does not remove the element in the set, though
	Peek() (T, error)
	// Size returns the size of the set
	Size() int64
}

// SetEqualsFunction defines what is equality for the set
type SetEqualsFunction[T any] func(a, b T) bool

// AbstractSetBuilder builds a new empty set, with the equality
type AbstractSetBuilder[T any] func(SetEqualsFunction[T]) (AbstractSet[T], error)
