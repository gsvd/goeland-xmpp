package id

import (
	"github.com/google/uuid"
)

// Centralize ID generation,
// This function returns a new UUID string.
// Keeping it simple for the moment.
func New() string {
	return uuid.NewString()
}
