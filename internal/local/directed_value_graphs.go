package local

import "errors"

type DirectedValuesGraph[NV comparable, LV any] struct {
	content doubleMap[NV, LV]
}

func NewDirectedValuesGraph[NV comparable, LV any]() DirectedValuesGraph[NV, LV] {
	var impl doubleMap[NV, LV] = newDoubleMap[NV, LV]()
	var result DirectedValuesGraph[NV, LV] = DirectedValuesGraph[NV, LV]{
		content: impl,
	}

	return result
}

func (d DirectedValuesGraph[NV, LV]) SetLink(source, destination NV, value LV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.putValue(source, destination, value)

	return nil
}

func (d DirectedValuesGraph[NV, LV]) RemoveLink(source, destination NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.removeValue(source, destination)

	return nil
}

func (d DirectedValuesGraph[NV, LV]) AddNode(source NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.putElement(source)

	return nil
}

func (d DirectedValuesGraph[NV, LV]) RemoveNode(source NV) error {
	if d.content == nil {
		return errors.New("nil content, graph initialization failure")
	}

	d.content.removeElement(source)

	return nil
}

func (d DirectedValuesGraph[NV, LV]) Neighbors(source NV) (map[NV]LV, error) {
	if d.content == nil {
		return nil, nil
	}

	rows, found := d.content[source]
	if !found || len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func (d DirectedValuesGraph[NV, LV]) LinkValue(source, destination NV) (LV, bool, error) {
	var result LV

	if d.content == nil {
		return result, false, errors.New("nil content, graph initialization failure")
	}

	result, found := d.content.getValue(source, destination)
	return result, found, nil
}
