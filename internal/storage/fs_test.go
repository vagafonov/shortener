package storage

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/pkg/entity"
)

const fileName = "test"

type FileSystemStorageTestSuite struct {
	suite.Suite
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FileSystemStorageTestSuite))
}

func (s *FileSystemStorageTestSuite) SetupSuite() {
}

func (s *FileSystemStorageTestSuite) TestAdd() {
	ctx := context.Background()
	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	userID := uuid.Must(uuid.NewUUID())
	resultURL, err := fss.Add(ctx, "1", "2", userID)
	s.Require().NoError(err)
	data, err := os.ReadFile(fileName)
	s.Require().NoError(err)
	urlActual := &entity.URL{}
	err = json.Unmarshal(data, &urlActual)
	s.Require().NoError(err)

	s.Equal(entity.URL{
		UUID:     resultURL.UUID,
		Short:    "1",
		Original: "2",
		UserID:   userID,
	}, *urlActual)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestGetAll() {
	ctx := context.Background()
	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()
	entityURL := &entity.URL{
		ID:       "123456",
		UUID:     uuid.New(),
		Short:    "short1",
		Original: "full1",
		UserID:   uuid.New(),
	}
	_, err = addTestURLToFile(entityURL)
	s.Require().NoError(err)

	resultURLs, err := fss.GetAll(ctx)
	s.Require().NoError(err)
	exp := []*entity.URL{entityURL}
	s.Require().Equal(exp, resultURLs)
	os.Remove(fileName)
}

//nolint:dupl
func (s *FileSystemStorageTestSuite) TestGetByHash() {
	ctx := context.Background()
	_, err := addTestURLToFile(&entity.URL{
		ID:       "1",
		UUID:     uuid.New(),
		Short:    "short1",
		Original: "full1",
		UserID:   uuid.New(),
	})
	s.Require().NoError(err)

	entityURL := &entity.URL{
		ID:       "2",
		UUID:     uuid.New(),
		Short:    "short2",
		Original: "full2",
		UserID:   uuid.New(),
	}
	_, err = addTestURLToFile(entityURL)
	s.Require().NoError(err)

	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	url, err := fss.GetByHash(ctx, "short2")
	s.Require().NoError(err)

	s.Require().Equal(entityURL, url)
	os.Remove(fileName)
}

//nolint:dupl
func (s *FileSystemStorageTestSuite) TestGetByURL() {
	ctx := context.Background()
	_, err := addTestURLToFile(&entity.URL{
		ID:       "1",
		UUID:     uuid.New(),
		Short:    "short1",
		Original: "full1",
		UserID:   uuid.New(),
	})
	s.Require().NoError(err)

	entityURL := &entity.URL{
		ID:       "2",
		UUID:     uuid.New(),
		Short:    "short2",
		Original: "full2",
		UserID:   uuid.New(),
	}
	_, err = addTestURLToFile(entityURL)
	s.Require().NoError(err)

	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	url, err := fss.GetByURL(ctx, "full2")
	s.Require().NoError(err)

	s.Require().Equal(entityURL, url)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestAddBatch() {
	ctx := context.Background()
	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	URLs := []*entity.URL{
		{
			UUID:     uuid.UUID{},
			Short:    "",
			Original: "",
		},
	}
	totalCreated, err := fss.AddBatch(ctx, URLs)
	s.Require().NoError(err)
	s.Require().Equal(1, totalCreated)
	os.Remove(fileName)
}

//nolint:unparam
func addTestURLToFile(u *entity.URL) (*entity.URL, error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)

	return u, encoder.Encode(u)
}

func (s *FileSystemStorageTestSuite) TestGetUserURLs() {
	ctx := context.Background()

	s.Run("get one user URL", func() {
		userID := uuid.New()
		_, err := addTestURLToFile(&entity.URL{
			ID:       "2",
			UUID:     uuid.New(),
			Short:    "short2",
			Original: "full2",
			UserID:   uuid.New(),
		})
		s.Require().NoError(err)
		entityURL := &entity.URL{
			ID:       "2",
			UUID:     uuid.New(),
			Short:    "short2",
			Original: "full2",
			UserID:   userID,
		}
		_, err = addTestURLToFile(entityURL)
		s.Require().NoError(err)

		fss, err := NewFileSystemStorage(fileName)
		s.Require().NoError(err)
		defer fss.Close()

		url, err := fss.GetAllURLsByUser(ctx, userID, "")
		s.Require().NoError(err)

		s.Require().Equal([]*entity.URL{entityURL}, url)
		os.Remove(fileName)
	})
}
