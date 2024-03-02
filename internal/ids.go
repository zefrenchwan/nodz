package internal

import (
	"github.com/google/uuid"
)

// NewUniqueId returns an UUID
func NewUniqueId() string {
	return uuid.NewString()
}
