package internal

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

// TypePropertiesLink is a link with a type and properties (Neo4j model).
type TypePropertiesLink[N Node] struct {
	// nodeSource is the source of the link
	linkSource N
	// nodeDestination is the destination of the link
	linkDestination N
	// linkType is the type of link
	linkType string
	// linkProperties are the properties of the link
	linkProperties map[string]string
}

// Neo4jLink is a type alias to ease use of Neo4j
type Neo4jLink TypePropertiesLink[*LabelPropertiesNode]

// NewTypePropertiesLink returns a new node with a given type, from a source to a destination
func NewTypePropertiesLink[N Node](linkType string, source, destination N) TypePropertiesLink[N] {
	return TypePropertiesLink[N]{
		linkSource:      source,
		linkDestination: destination,
		linkType:        linkType,
		linkProperties:  make(map[string]string),
	}
}

// LinkType returns the type of the node
func (tpl *TypePropertiesLink[N]) LinkType() string {
	var result string
	if tpl != nil {
		result = tpl.linkType
	}

	return result
}

// Source returns the source of the link, assuming receiver is not nil
func (tpl *TypePropertiesLink[N]) Source() N {
	return tpl.linkSource
}

// Destination returns the destination of the link, assuming receiver is not nil
func (tpl *TypePropertiesLink[N]) Destination() N {
	return tpl.linkDestination
}

// SameLink returns true for same type, same source and destination, false for any other case
func (tpl *TypePropertiesLink[N]) SameLink(link TypePropertiesLink[N]) bool {
	if tpl == nil {
		return false
	}

	if link.linkType != tpl.linkType {
		return false
	}

	return tpl.linkSource.SameNode(link.linkSource) && tpl.linkDestination.SameNode(link.linkDestination)
}

// Properties returns a copy of the properties of the link, nil for nil
func (tpl *TypePropertiesLink[N]) Properties() map[string]string {
	if tpl == nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range tpl.linkProperties {
		result[k] = v
	}

	return result
}

// SetProperty adds (or changes) this property of a link
func (tpl *TypePropertiesLink[N]) SetProperty(key, value string) {
	if tpl != nil {
		tpl.linkProperties[key] = value
	}
}

// GetProperty returns the value if any, true if link has the value and false otherwise
func (tpl *TypePropertiesLink[N]) GetProperty(key string) (string, bool) {
	if tpl == nil {
		return "", false
	}

	value, has := tpl.linkProperties[key]
	return value, has
}

// RemoveProperty removes a property per key
func (tpl *TypePropertiesLink[N]) RemoveProperty(key string) {
	if tpl != nil {
		delete(tpl.linkProperties, key)
	}
}
