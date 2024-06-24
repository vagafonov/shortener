package hasher

import (
	"math/rand"
)

// Alphabet default available alphabet of symbols from which a hash is created.
const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandHasher creates hashes.
type RandHasher struct {
	alphabet string
}

// Constructor for RandHasher.
func NewRandHasher(a string) *RandHasher {
	return &RandHasher{
		alphabet: a,
	}
}

// Hash make hash.
func (h *RandHasher) Hash(length int) string {
	if length <= 0 {
		return ""
	}
	shortKey := make([]byte, length)
	num := len(h.alphabet) - 1
	for i := range shortKey {
		shortKey[i] = h.alphabet[rand.Intn(num)+1] //nolint:gosec
	}

	return string(shortKey)
}
