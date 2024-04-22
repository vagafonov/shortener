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

	resultURL, err := fss.Add(ctx, "1", "2")
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
	}, *urlActual)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestGetAll() {
	ctx := context.Background()
	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	_, err = addTestURLToFile(uuid1, "short1", "full1")
	s.Require().NoError(err)
	_, err = addTestURLToFile(uuid2, "short2", "full2")
	s.Require().NoError(err)
	resultURLs, err := fss.GetAll(ctx)
	s.Require().NoError(err)
	exp := []entity.URL{
		{
			UUID:     uuid1,
			Short:    "short1",
			Original: "full1",
		},
		{
			UUID:     uuid2,
			Short:    "short2",
			Original: "full2",
		},
	}
	s.Require().Equal(exp, resultURLs)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestGetByHash() {
	ctx := context.Background()
	_, err := addTestURLToFile(uuid.New(), "short1", "full")
	s.Require().NoError(err)
	u2, err := addTestURLToFile(uuid.New(), "short2", "full")
	s.Require().NoError(err)
	_, err = addTestURLToFile(uuid.New(), "short3", "full")
	s.Require().NoError(err)

	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	url, err := fss.GetByHash(ctx, "short2")
	s.Require().NoError(err)

	s.Require().Equal(&entity.URL{
		UUID:     u2.UUID,
		Short:    "short2",
		Original: "full",
	}, url)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestGetByURL() {
	ctx := context.Background()
	_, err := addTestURLToFile(uuid.New(), "short", "full1")
	s.Require().NoError(err)
	u2, err := addTestURLToFile(uuid.New(), "short", "full2")
	s.Require().NoError(err)
	_, err = addTestURLToFile(uuid.New(), "short", "full3")
	s.Require().NoError(err)

	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	url, err := fss.GetByURL(ctx, "full2")
	s.Require().NoError(err)

	s.Require().Equal(&entity.URL{
		UUID:     u2.UUID,
		Short:    "short",
		Original: "full2",
	}, url)
	os.Remove(fileName)
}

func (s *FileSystemStorageTestSuite) TestAddBatch() {
	ctx := context.Background()
	fss, err := NewFileSystemStorage(fileName)
	s.Require().NoError(err)
	defer fss.Close()

	URLs := []entity.URL{
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

func addTestURLToFile(id uuid.UUID, shortURL string, fullURL string) (*entity.URL, error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	urlEntity := &entity.URL{
		UUID:     id,
		Short:    shortURL,
		Original: fullURL,
	}
	encoder := json.NewEncoder(f)

	return urlEntity, encoder.Encode(urlEntity)
}
