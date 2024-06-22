package hasher

type Hasher interface {
	Hash(length int) string
}
