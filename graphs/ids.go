package graphs

import (
	"github.com/google/uuid"
)

// WithId defines an object that may be identifiable within a id.
// An Id is not something to change, even if it stays unique over time.
type WithId interface {
	// Id returns the id of the object, stable over time
	Id() string
}

// NewUniqueId returns an UUID
func NewUniqueId() string {
	return uuid.NewString()
}
