package hasher

import (
	"strings"
)

// MockHasher mock.
type MockHasher struct{}

// NewMockHasher mock.
func NewMockHasher() *MockHasher {
	return &MockHasher{}
}

// Hash mock.
func (h *MockHasher) Hash(length int) string {
	return strings.Repeat("*", length)
}
