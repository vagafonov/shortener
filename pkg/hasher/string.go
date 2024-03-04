package hash

import (
	"math/rand"
)

type String struct {
}

func NewStringHasher() *String {
	return &String{}
}

func (h *String) Hash(length int) string {

	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(shortKey)
}
