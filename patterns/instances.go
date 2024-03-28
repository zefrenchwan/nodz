package patterns

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FormalInstance is an instance of the class.
// It has an id, info about instantiating class to find back which class instantiated it.
// Plus, values store time dependent values.
// An entity follows no formal definition, its attributes may be anything.
// But an instance is an instance of a class and then follows the class definition.
type FormalInstance struct {
	// id of the instance, should be unique
	id string
	// metadata is the link to the instantiating class
	metadata *FormalClass
	// values contains attributes and their time dependent values
	values TimeValues
}

// NewFormalInstance returns a new formal instance of a given class
func NewFormalInstance(instantiatingClass *FormalClass) (FormalInstance, error) {
	var result FormalInstance

	if instantiatingClass == nil {
		return result, errors.New("nil class to build instance")
	}

	result = FormalInstance{
		id:       uuid.NewString(),
		metadata: instantiatingClass,
		values:   NewTimeValues(),
	}

	return result, nil
}

// Attributes returns the attributes of the instance as a sorted slice.
// If current instance is nil, it raises an error.
func (i *FormalInstance) Attributes() ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if i.metadata == nil {
		return nil, errors.New("nil instantiating class")
	}

	return i.metadata.ListAttributes(), nil
}

// SetValue sets the value for an attribute during the full period.
func (i *FormalInstance) SetValue(attribute string, value string) error {
	if i == nil {
		return errors.New("nil instance")
	} else if i.metadata == nil {
		return errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	}

	return i.values.SetValue(attribute, value)
}

// AddValue sets the value of an attribute during a given period.
// It updates the periods of the other values (for the same attribute) accordingly.
func (i *FormalInstance) AddValue(attribute string, value string, validity Period) error {
	// nil should return an error, empty period should change nothing
	if i == nil {
		return errors.New("nil instance")
	} else if validity.IsEmptyPeriod() {
		return nil
	} else if i.metadata == nil {
		return errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	}

	return i.values.AddValue(attribute, value, validity)
}

// ValuesForAttribute returns the values for an attribute as a sorted slice
func (i *FormalInstance) ValuesForAttribute(attribute string) ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if i.metadata == nil {
		return nil, errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	}

	return i.values.ValuesForAttribute(attribute)
}

// TimeValuesForAttribute returns, for each value of the attribute, the matching time intervals
func (i *FormalInstance) TimeValuesForAttribute(attribute string) (map[string][]Interval[time.Time], error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if i.metadata == nil {
		return nil, errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	}

	return i.values.TimeValuesForAttribute(attribute)
}
