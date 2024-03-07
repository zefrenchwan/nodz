package graphs

// Matrix is the definition of a square matrix of size Size.
type Matrix[V any] interface {
	// Size returns the number of lines = number of columns of the matrix
	Size() int
	// SetValue sets the value at a given position (line i, column j).
	// An error is allowed for implementations
	SetValue(i, j int, value V) error
	// GetValue returns the value at position i,j (line i, column j).
	// If there is no value, returns the default value, false.
	// An error is allowed for implementations
	GetValue(i, j int) (V, bool, error)
}
