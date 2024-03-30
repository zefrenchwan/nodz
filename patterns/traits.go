package patterns

import (
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// Trait represent a behavior or capacity for a given period.
// For instance, an entity (let us say a person)
// may be considered as a trait (student) for a given period.
// For instance, an entity (the Eiffel tower) may be seen as
// a meeting point during a specific period and a monument since its creation.
// Both are traits (meeting point, monument) for the same entity.
type Trait struct {
	// id of the trait, should be unique
	id string
	// entityId is the id of the entity that owns the trait
	entityId string
	// labels are case insensitive information about the trait
	labels []string
	// content is the activity of the trait and its attributes
	content ActiveTimeValues
}

// NewTrait builds a trait for that entity with given labels
func NewTrait(entityId string, labels []string) Trait {
	res, _ := NewTraitDuring(entityId, labels, NewFullPeriod())
	return res
}

// NewTraitDuring builds a trait for that entity with given labels.
// Trait is active only during given period.
// If period is empty, it returns an error
func NewTraitDuring(entityId string, labels []string, period Period) (Trait, error) {
	var result Trait
	if period.IsEmptyPeriod() {
		return result, errors.New("empty period for trait")
	}

	values := NewActiveTimeValues()
	values.SetActivity(period)
	result.content = values
	result.entityId = entityId
	result.id = uuid.NewString()
	result.labels = append(result.labels, labels...)

	return result, nil
}

// AddLabel appends a label to the set of labels (no duplicate)
func (t *Trait) AddLabel(label string) {
	if t == nil {
		return
	}

	t.labels = append(t.labels, label)
}

// Labels returns the labels as a sorted slice and nil for nil
// Sorting allows to compare labels using the order.
func (t *Trait) Labels() []string {
	if t == nil {
		return nil
	} else if len(t.labels) == 0 {
		return nil
	}

	values := make([]string, len(t.labels))
	copy(values, t.labels)
	slices.Sort(values)
	return values
}

// RemoveLabel removes, if any, the label key
func (t *Trait) RemoveLabel(key string) {
	if t == nil {
		return
	}

	if !slices.ContainsFunc(t.labels, func(a string) bool { return strings.EqualFold(a, key) }) {
		t.labels = append(t.labels, key)
	}
}
