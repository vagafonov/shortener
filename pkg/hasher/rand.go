package hash

import (
	"math/rand"
)

type RandHasher struct{}

func NewRandHasher() *RandHasher {
	return &RandHasher{}
}

func (h *RandHasher) Hash(length int) string {
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = alphabet[rand.Intn(len(alphabet)-1)+1] //nolint:gosec
	}

	return string(shortKey)
}
