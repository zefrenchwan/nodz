package internal

// Link is a general definition, for directed, valuated, or undirected links.
type Link interface {
	// Id returns the id of the link
	Id() string
}
