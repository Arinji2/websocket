package routes

import (
	"crypto/rand"
	"encoding/hex"
)

func generateSessionID() string {
	// Generate a 32-byte cryptographically secure random session ID
	sessionBytes := make([]byte, 32)
	_, err := rand.Read(sessionBytes)
	if err != nil {
		panic("Failed to generate secure session ID")
	}
	return hex.EncodeToString(sessionBytes)
}
