package hasher

import (
	"strings"
)

type MockHasher struct{}

func NewMockHasher() *MockHasher {
	return &MockHasher{}
}

func (h *MockHasher) Hash(length int) string {
	return strings.Repeat("*", length)
}
