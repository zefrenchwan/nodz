package patterns

import (
	"github.com/google/uuid"
)

// Entity defines a unique "real world object" with time dependent traits
type Entity struct {
	// id of the entity, should be unique
	id string
}

// NewEntity returns an empty entity with a generated id
func NewEntity() Entity {
	return Entity{
		id: uuid.NewString(),
	}
}

// Id returns the id of an entity
func (e *Entity) Id() string {
	return e.id
}
