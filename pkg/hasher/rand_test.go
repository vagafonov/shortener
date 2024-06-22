package hasher

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RandHasherTestSuite struct {
	suite.Suite
}

func TestRandHasherTestSuite(t *testing.T) {
	suite.Run(t, new(RandHasherTestSuite))
}

func (s *RandHasherTestSuite) TestCreateURL() {

	testCases := []struct {
		name     string
		inputLen int
		expected string
	}{
		{
			name:     "zero length",
			inputLen: 0,
			expected: "",
		},
		{
			name:     "zero length",
			inputLen: 1,
			expected: "a",
		},
		{
			name:     "zero length",
			inputLen: 5,
			expected: "aaaaa",
		},
	}

	rh := NewRandHasher("aaa")

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(tc.expected, rh.Hash(tc.inputLen))
		})
	}
}

func BenchmarkRandHash(b *testing.B) {
	rh := NewRandHasher(Alphabet)

	for i := 0; i < b.N; i++ {
		rh.Hash(10)
	}
}
