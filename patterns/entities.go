package patterns

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Entity defines a unique "real world object" with time dependent values
type Entity struct {
	// id of the entity, should be unique
	id string
	// values are the time dependent values of the entity
	values TimeValues
}

// NewEntity returns an empty entity with a generated id
func NewEntity() Entity {
	return Entity{
		id:     uuid.NewString(),
		values: NewTimeValues(),
	}
}

// Attributes returns the attributes of the entity with at least a value
func (e *Entity) Attributes() ([]string, error) {
	if e == nil {
		return nil, errors.New("nil entity")
	} else if e.values == nil {
		return nil, nil
	}

	return e.values.Attributes()
}

// SetValue sets the value for an attribute during the full period.
func (e *Entity) SetValue(attribute string, value string) error {
	if e == nil {
		return errors.New("nil entity")
	} else if e.values == nil {
		e.values = NewTimeValues()
	}

	return e.values.SetValue(attribute, value)
}

// AddValue sets the value of an attribute during a given period.
// It updates the periods of the other values (for the same attribute) accordingly.
func (e *Entity) AddValue(attribute string, value string, validity Period) error {
	// nil should return an error, empty period should change nothing
	if e == nil {
		return errors.New("nil entity")
	} else if validity.IsEmptyPeriod() {
		return nil
	} else if e.values == nil {
		e.values = NewTimeValues()
	}

	return e.values.AddValue(attribute, value, validity)
}

// ValuesForAttribute returns the values for an attribute as a sorted slice
func (e *Entity) ValuesForAttribute(attribute string) ([]string, error) {
	if e == nil {
		return nil, errors.New("nil entity")
	} else if e.values == nil {
		return nil, nil
	}

	return e.values.ValuesForAttribute(attribute)
}

// TimeValuesForAttribute returns, for each value of the attribute, the matching time intervals
func (e *Entity) TimeValuesForAttribute(attribute string) (map[string][]Interval[time.Time], error) {
	if e == nil {
		return nil, errors.New("nil entity")
	} else if e.values == nil {
		return nil, nil
	}

	return e.values.TimeValuesForAttribute(attribute)
}
