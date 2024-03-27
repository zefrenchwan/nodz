package patterns

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

// FormalInstance is an instance of the class.
// It has an id, and a class name to find back which class instantiated it.
// Plus, attributes have time dependant values.
// As an example, assume a person is represented by a class Person, with attributes birth date, etc.
// Those values will be set once and will not vary over time.
// But some will, for instance, address.
// This code uses time dependant values for this purpose.
// You create an instance, then add a value for an attribute during a given period.
type FormalInstance struct {
	// id of the instance, should be unique
	id string
	// className is the name of the instantiating class
	className string
	// values is a map of attributes, and the value for that attribute during a period.
	values map[string]map[string]*Period
}

// NewFormalInstance returns a new formal instance of a given class
func NewFormalInstance(instantiatingClass FormalClass) FormalInstance {
	result := FormalInstance{
		id:        uuid.NewString(),
		className: instantiatingClass.name,
		values:    make(map[string]map[string]*Period),
	}

	for _, attribute := range instantiatingClass.attributes {
		result.values[attribute] = make(map[string]*Period)
	}

	return result
}

// Attributes returns the attributes of the instance as a sorted slice.
// If current instance is nil, it raises an error.
func (i *FormalInstance) Attributes() ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if len(i.values) == 0 {
		return nil, nil
	}

	result := make([]string, len(i.values))
	index := 0
	for k := range i.values {
		result[index] = k
		index++
	}

	slices.Sort(result)
	return result, nil

}

// SetValue sets the value for an attribute during the full period.
func (i *FormalInstance) SetValue(attribute string, value string) error {
	var matchingAttributeMap map[string]*Period

	// find matching map for this attribute, if any
	if i == nil {
		return errors.New("nil instance")
	} else if value, found := i.values[attribute]; !found {
		return fmt.Errorf("no attribute %s in class %s", attribute, i.className)
	} else {
		matchingAttributeMap = value
	}

	// clean the map
	for k := range matchingAttributeMap {
		delete(matchingAttributeMap, k)
	}

	// add value -> full
	period := NewFullPeriod()
	matchingAttributeMap[value] = &period

	// no error
	return nil
}

// AddValue sets the value of an attribute during a given period.
// It updates the periods of the other values (for the same attribute) accordingly.
func (i *FormalInstance) AddValue(attribute string, value string, validity Period) error {
	// nil should return an error, empty period should change nothing
	if i == nil {
		return errors.New("nil instance")
	} else if validity.IsEmptyPeriod() {
		return nil
	}

	// find matching attribute map if any
	var matchingAttributeMap map[string]*Period
	if value, found := i.values[attribute]; !found {
		return fmt.Errorf("no attribute %s in class %s", attribute, i.className)
	} else {
		matchingAttributeMap = value
	}

	// for each attribute value different than parameter, get the intersection with the validity
	for valueForAttribute, matchingPeriod := range matchingAttributeMap {
		// will change value later
		if valueForAttribute == value {
			continue
		}

		// remove the period for the other attribute.
		// And if it is empty, value should be removed
		matchingPeriod.Remove(validity)
		if matchingPeriod.IsEmptyPeriod() {
			delete(matchingAttributeMap, value)
		}
	}

	// and set the value
	if matchingPeriod, found := matchingAttributeMap[value]; found {
		matchingPeriod.Add(validity)
	} else {
		copyOfPeriod := NewPeriodCopy(validity)
		matchingAttributeMap[value] = &copyOfPeriod
	}

	return nil
}

// ValuesForAttribute returns the values for an attribute as a sorted slice
func (i *FormalInstance) ValuesForAttribute(attribute string) ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if attributeValues, found := i.values[attribute]; !found {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.className)
	} else if len(attributeValues) == 0 {
		return nil, nil
	} else {
		result := make([]string, len(attributeValues))

		index := 0
		for k := range attributeValues {
			result[index] = k
			index++
		}

		slices.Sort(result)
		return result, nil
	}
}

// TimeValuesForAttribute returns, for each value of the attribute, the matching time intervals
func (i *FormalInstance) TimeValuesForAttribute(attribute string) (map[string][]Interval[time.Time], error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if attributeValues, found := i.values[attribute]; !found {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.className)
	} else if len(attributeValues) == 0 {
		return nil, nil
	} else {
		result := make(map[string][]Interval[time.Time])

		for value, period := range attributeValues {
			// should not happen
			if period == nil {
				continue
			}

			result[value] = period.AsIntervals()
		}

		return result, nil
	}
}
