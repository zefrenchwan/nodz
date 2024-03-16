package graphs

// Link is a link between nodes, may have values, directed or undirected, or... something else.
// No matter if the node has a direction or not, we use source and destination.
//
// A note about link equality.
// For directed links, it is implicit that if two links are indeed the same,
// they have same source and same destination.
// For undirected links, it is implicit that if two links are indeed the same,
// their sources and destination match.
// It means that either sources and destinations are the same,
// or the source of one is the destination of the other, and vice versa.
// It will not be tested for every implementation, it is up to you to keep that in mind.
type Link[N Node] interface {
	// Test if link is the same as another one.
	SameLink(other Link[N]) bool
	// Source of the link
	Source() N
	// Destination of the link
	Destination() N
	// IsDirected returns true for directed links, false for undirected
	IsDirected() bool
}

// LinksIterator defines iterator over links.
// In some circumstances, some graphs are too huge to just return a []Link.
type LinksIterator[N Node, L Link[N]] GeneralIterator[L]

// FollowLink returns the destination of a link from a source, if possible.
// Formally, if link is directed link, it is possible to go from source to destination.
// If link is undiregcted, it is possible to go from an endpoint to another.
// If parameter is neither source nor destination, no move is possible
func FollowLink[N Node, L Link[N]](from N, link L) (bool, N) {
	var empty N

	if link.Source().SameNode(from) {
		return true, link.Destination()
	} else if link.Destination().SameNode(from) {
		if link.IsDirected() {
			return false, empty
		} else {
			return true, link.Source()
		}
	}

	return false, empty
}
