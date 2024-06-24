package encrypting

import "crypto/rand"

// GenerateRandom generate random byte slice.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
