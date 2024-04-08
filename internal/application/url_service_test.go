package application

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

const fileName = "test-file-db"

type URLSuite struct {
	suite.Suite
	service Service
	cnt     *Container
}

func TestURLSuite(t *testing.T) {
	suite.Run(t, new(URLSuite))
}

func (s *URLSuite) SetupSuite() {
	ms := storage.NewMemoryStorage()
	fss, err := storage.NewFileSystemStorage(fileName)
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig("test", "http://test:8080", fileName, "")
	hasher := hash.NewMockHasher()
	s.cnt = NewContainer(
		cfg,
		ms,
		fss,
		hasher,
	)
	s.service = NewService(
		s.cnt.GetLogger(),
		s.cnt.GetStorage(),
		s.cnt.GetBackupStorage(),
		s.cnt.GetHasher(),
	)
}

func (s *URLSuite) TestRestoreURLs() {
	_, err := s.cnt.GetBackupStorage().Add("short1", "full1")
	s.Require().NoError(err)
	_, err = s.cnt.GetBackupStorage().Add("short2", "full2")
	s.Require().NoError(err)

	s.Run("restore all urls", func() {
		URLs, err := s.cnt.GetStorage().GetAll()
		s.Require().NoError(err)
		s.Require().Empty(URLs)
		totalRestored, err := s.service.RestoreURLs(fileName)
		s.Require().NoError(err)
		s.Require().Equal(2, totalRestored)
		URLs, err = s.cnt.GetStorage().GetAll()
		s.Require().NoError(err)
		s.Require().Len(URLs, 2)
		os.Remove(fileName)
	})
}
