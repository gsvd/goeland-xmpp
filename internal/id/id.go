package id

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

// This function returns a new UUID string.
// Keeping it simple for the moment.
func NewUUID() string {
	return uuid.NewString()
}

func NewShortHex() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
