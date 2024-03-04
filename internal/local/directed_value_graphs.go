package local

import "errors"

// DirectedValuesGraph is a directed graph, values based.
// NV is the type of the nodes, LV is the one for values.
// For instance, if your graph looks like "person A" - 10.0 -> "course about networks",
// then NV = string and LV = "float".
type DirectedValuesGraph[NV comparable, LV any] struct {
	// content is the inner content of the graph. nil value may raise errors, not recommended
	content doubleMap[NV, LV]
}

// NewDirectedValuesGraph returns a new empty graph
func NewDirectedValuesGraph[NV comparable, LV any]() DirectedValuesGraph[NV, LV] {
	var impl doubleMap[NV, LV] = newDoubleMap[NV, LV]()
	var result DirectedValuesGraph[NV, LV] = DirectedValuesGraph[NV, LV]{
		content: impl,
	}

	return result
}

// SetLink set values for a link, and adds missing nodes in the graph
func (d DirectedValuesGraph[NV, LV]) SetLink(source, destination NV, value LV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.putValue(source, destination, value)

	return nil
}

// RemoveLink removes the link, keeps the nodes
func (d DirectedValuesGraph[NV, LV]) RemoveLink(source, destination NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.removeValue(source, destination)

	return nil
}

// AddNode adds the node if not in the graph, but does nothing for an existing node
func (d DirectedValuesGraph[NV, LV]) AddNode(source NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.putElement(source)

	return nil
}

// RemoveNode removes the node and all links for that node
func (d DirectedValuesGraph[NV, LV]) RemoveNode(source NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.removeElement(source)

	return nil
}

// Neighbors returns, if any, neighbors from source as a map of destinations and values.
// It returns the same result (nil) if source is not in the graph and source has no neghbor.
func (d DirectedValuesGraph[NV, LV]) Neighbors(source NV) (map[NV]LV, error) {
	if d.content == nil {
		return nil, nil
	}

	rows := d.content.getElement(source)
	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

// LinkValue returns link value, if any, for link source -> destination
func (d DirectedValuesGraph[NV, LV]) LinkValue(source, destination NV) (LV, bool, error) {
	var result LV

	if d.content == nil {
		return result, false, errors.New("nil content, graph initialization failure")
	}

	result, found := d.content.getValue(source, destination)
	return result, found, nil
}
