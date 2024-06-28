package hasher

import (
	"math/rand"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandHasher struct {
	alphabet string
}

func NewRandHasher(a string) *RandHasher {
	return &RandHasher{
		alphabet: a,
	}
}

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
