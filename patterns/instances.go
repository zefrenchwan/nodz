package patterns

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FormalInstance is an instance of the class to decorate an entity.
// It "forgets" some attributes of the entity to make it appear as an instance of a class.
// Still, when changing a value, if the class allows it, then it changes also the decorated entity.
// Same when reading a value, this value comes from the decorated entity if the class has the attribute.
// It has its own id, and info about instantiating class to find back which class instantiated it.
// An entity follows no formal definition, its attributes may be anything.
// But an instance is an instance of a class and then follows the class definition.
type FormalInstance struct {
	// id of the instance, should be unique
	id string
	// metadata is the link to the instantiating class
	metadata *FormalClass
	// decoratedEntity is the entity to consider as an instance of a class
	decoratedEntity *Entity
}

// NewFormalInstance returns a new formal instance of a given class from a specific entity.
// If class is nil, it fails.
// If entity is nil, it builds a new one.
func NewFormalInstance(instantiatingClass *FormalClass, entity *Entity) (FormalInstance, error) {
	var result FormalInstance

	if instantiatingClass == nil {
		return result, errors.New("nil class to build instance")
	} else if entity == nil {
		result.decoratedEntity = new(Entity)
		*result.decoratedEntity = NewEntity()
	} else {
		result.decoratedEntity = entity
	}

	result.id = uuid.NewString()
	result.metadata = instantiatingClass

	return result, nil
}

// GetDecoratedEntity returns the decorated entity
func (i *FormalInstance) GetDecoratedEntity() *Entity {
	if i == nil {
		return nil
	}

	return i.decoratedEntity
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
	} else if i.decoratedEntity == nil {
		return errors.New("no decorated entity")
	}

	return i.decoratedEntity.SetValue(attribute, value)
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
	} else if i.decoratedEntity == nil {
		return errors.New("no decorated entity")
	}

	return i.decoratedEntity.AddValue(attribute, value, validity)
}

// ValuesForAttribute returns the values for an attribute as a sorted slice
func (i *FormalInstance) ValuesForAttribute(attribute string) ([]string, error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if i.metadata == nil {
		return nil, errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	} else if i.decoratedEntity == nil {
		return nil, errors.New("no decorated entity")
	}

	return i.decoratedEntity.ValuesForAttribute(attribute)
}

// TimeValuesForAttribute returns, for each value of the attribute, the matching time intervals
func (i *FormalInstance) TimeValuesForAttribute(attribute string) (map[string][]Interval[time.Time], error) {
	if i == nil {
		return nil, errors.New("nil instance")
	} else if i.metadata == nil {
		return nil, errors.New("nil metadata for instance")
	} else if !i.metadata.hasAttribute(attribute) {
		return nil, fmt.Errorf("no attribute %s in class %s", attribute, i.metadata.name)
	} else if i.decoratedEntity == nil {
		return nil, errors.New("no decorated entity")
	}

	return i.decoratedEntity.TimeValuesForAttribute(attribute)
}
