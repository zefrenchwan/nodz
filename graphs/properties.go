package graphs

// WithProperties is a simple properties implementation, with no type (book)
type WithProperties interface {
	// GetProperty returns the value of a property, if any (or empty, false)
	GetProperty(key string) (string, bool)
	// SetProperty sets the value for a property
	SetProperty(key, value string)
	// RemoveProperty removes a property by key
	RemoveProperty(key string)
	// PropertyKeys returns the keys of the available properties
	PropertyKeys() []string
}
