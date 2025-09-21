package crypto

import "crypto/rand"

// RandBytes returns a random byte slice of the specified length.
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
