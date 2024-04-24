package storage

import (
	"context"
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
	ctx := context.Background()
	ms := NewMemoryStorage()
	batchURLs := []entity.URL{
		{
			UUID:     uuid.UUID{},
			Short:    "a",
			Original: "aaa",
		},
	}
	tc, err := ms.AddBatch(ctx, batchURLs)
	s.Require().NoError(err)
	s.Require().Equal(1, tc)

	allURLs, err := ms.GetAll(ctx)
	s.Require().Equal(batchURLs, allURLs)
	s.Require().NoError(err)
}
