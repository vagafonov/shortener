package hasher

// Hasher hashing interface.
type Hasher interface {
	Hash(length int) string
}
