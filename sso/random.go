package sso

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomToken generates a random string of n characters,
// alphanumeric plus - and _.
func RandomToken(n int) string {
	// b is used for raw bytes, which are then base64-encoded into s.
	b := make([]byte, (n+1)*3/4)

	_, err := rand.Read(b)
	if err != nil {
		// An error here means something seriously wrong with our runtime system.
		panic(err)
	}

	s := make([]byte, n+1)
	base64.RawURLEncoding.Encode(s, b)
	return string(s[:n])
}
