package graphs

// Link is a link between nodes, may have values, directed or undirected, or... something else.
// No matter if the node has a direction or not, we use source and destination.
// Dealing with direction is a question of neighborhood, not links per se.
// For instance:
// Given the undirected link a - link - b, then b is the neighbor of a, and vice versa.
// Given the directed link a - link -> b, b is in the neighborhood of a, but not the other way around.
type Link[N Node] interface {
	// Source of the link
	Source() N
	// Destination of the link
	Destination() N
}
