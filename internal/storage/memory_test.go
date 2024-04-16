package storage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/pkg/entity"
)

type MemoryStorageTestSuite struct {
	suite.Suite
}

func TestMemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}

func (s *MemoryStorageTestSuite) TestAdd() {
	s.T().Skip("create later")
}

func (s *MemoryStorageTestSuite) TestGetAll() {
	s.T().Skip("create later")
}

func (s *MemoryStorageTestSuite) TestGetByHash() {
	s.T().Skip("create later")
}

func (s *MemoryStorageTestSuite) TestGetByURL() {
	s.T().Skip("create later")
}

func (s *MemoryStorageTestSuite) TestAddBatch() {
	ms := NewMemoryStorage()
	batchURLs := []entity.URL{
		{
			UUID:     uuid.UUID{},
			Short:    "a",
			Original: "aaa",
		},
		{
			UUID:     uuid.UUID{},
			Short:    "b",
			Original: "bbb",
		},
	}
	tc, err := ms.AddBatch(batchURLs)
	s.Require().NoError(err)
	s.Require().Equal(2, tc)

	allURLs, err := ms.GetAll()
	s.Require().Equal(batchURLs, allURLs)
	s.Require().NoError(err)
}
