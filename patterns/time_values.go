package patterns

import (
	"errors"
	"slices"
	"time"
)

// TimeValues represents attributes with time dependent values.
// Keys are attributes name, values are the values per period.
type TimeValues map[string]map[string]*Period

// NewTimeValues returns a new empty TimeValues
func NewTimeValues() TimeValues {
	return make(map[string]map[string]*Period)
}

// ContainsAttribute returns true if receiver is not nil and it contains a non nil entry with that key
func (i TimeValues) ContainsAttribute(attr string) bool {
	switch value, found := i["attr"]; found {
	case true:
		return value != nil
	default:
		return false
	}
}

// Attributes returns the sorted slice of all attributes
func (i TimeValues) Attributes() ([]string, error) {
	result := make([]string, len(i))
	index := 0
	for k := range i {
		result[index] = k
		index++
	}

	slices.Sort(result)
	return result, nil

}

// SetValue sets a value for an attribute, for the full period.
func (i TimeValues) SetValue(attribute string, value string) error {
	var matchingAttributeMap map[string]*Period

	// find matching map for this attribute, if any
	if i == nil {
		return errors.New("nil instance")
	} else if value, found := i[attribute]; !found {
		// not found, allocate
		i[attribute] = make(map[string]*Period)
		matchingAttributeMap = i[attribute]
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
func (i TimeValues) AddValue(attribute string, value string, validity Period) error {
	// nil should return an error, empty period should change nothing
	if i == nil {
		return errors.New("nil instance")
	} else if validity.IsEmptyPeriod() {
		return nil
	}

	// find matching attribute map if any
	var matchingAttributeMap map[string]*Period
	if value, found := i[attribute]; !found {
		// not found, allocate
		i[attribute] = make(map[string]*Period)
		matchingAttributeMap = i[attribute]
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
func (i TimeValues) ValuesForAttribute(attribute string) ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if attributeValues, found := i[attribute]; !found {
		return nil, nil
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
func (i TimeValues) TimeValuesForAttribute(attribute string) (map[string][]Interval[time.Time], error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if attributeValues, found := i[attribute]; !found {
		return nil, nil
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
