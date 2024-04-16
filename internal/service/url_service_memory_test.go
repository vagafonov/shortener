package service

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/request"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/internal/storage"
	"github.com/vagafonov/shortener/pkg/entity"
	hasher "github.com/vagafonov/shortener/pkg/hasher"
)

const fileName = "test-file-db"

var ErrEmpty = errors.New("")

type ServiceURLMemorySuite struct {
	suite.Suite
	service       contract.Service
	cnt           *container.Container
	backupStorage *storage.FileSystemStorageMock
	mainStorage   *storage.MemoryStorageMock
}

func TestServiceURLSuite(t *testing.T) {
	suite.Run(t, new(ServiceURLMemorySuite))
}

func (s *ServiceURLMemorySuite) TearDownSuite() {
	os.Remove(fileName)
}

func (s *ServiceURLMemorySuite) SetupSuite() {
	fss, err := storage.NewFileSystemStorage(fileName)
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig("test", "http://test:8080", fileName, "test")
	lr := logger.CreateLogger(cfg.LogLevel)
	s.cnt = container.NewContainer(
		cfg,
		nil,
		fss,
		hasher.NewMockHasher(),
		lr,
		nil,
	)
	mainStorage, err := storage.StorageFactory(s.cnt, "memory-mock")
	if err != nil {
		log.Fatal(err)
	}
	s.mainStorage, _ = mainStorage.(*storage.MemoryStorageMock)
	s.cnt.SetMainStorage(s.mainStorage)

	backupStorage, err := storage.StorageFactory(s.cnt, "fs-mock")
	if err != nil {
		log.Fatal(err)
	}
	s.backupStorage, _ = backupStorage.(*storage.FileSystemStorageMock)
	s.cnt.SetBackupStorage(s.backupStorage)

	s.service = NewURLService(
		s.cnt.GetLogger(),
		s.cnt.GetMainStorage(),
		s.cnt.GetBackupStorage(),
		s.cnt.GetHasher(),
	)
}

func (s *ServiceURLMemorySuite) TestGetShortURL() {
	s.Run("get short url successfully", func() {
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "****",
			Original: "some_url",
		}
		s.mainStorage.SetGetByHashResponse(expEntity, nil)
		e, err := s.service.GetShortURL("some_url")
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceURLMemorySuite) TestMakeShortURL() {
	s.Run("make short url with empty storage", func() {
		s.mainStorage.SetGetByURLResponse(nil, nil)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "****",
			Original: "some_url",
		}
		s.mainStorage.SetAddResponse(expEntity, nil)
		e, err := s.service.MakeShortURL("some_url", 5)
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})

	s.Run("make short url with already exist urls", func() {
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "****",
			Original: "some_url",
		}
		s.mainStorage.SetGetByURLResponse(expEntity, nil)
		e, err := s.service.MakeShortURL("some_url", 5)
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceURLMemorySuite) TestRestoreURLs() {
	s.Run("restore all urls", func() {
		s.backupStorage.SetGetAllResponse([]entity.URL{
			{
				UUID:     uuid.UUID{},
				Short:    "",
				Original: "",
			},
		}, nil)
		s.mainStorage.SetAddResponse(&entity.URL{
			UUID:     uuid.UUID{},
			Short:    "",
			Original: "",
		}, nil)

		totalRestored, err := s.service.RestoreURLs(fileName)
		s.Require().NoError(err)
		s.Require().Equal(1, totalRestored)
	})

	s.Run("get all URLs failed", func() {
		s.backupStorage.SetGetAllResponse(nil, ErrEmpty)
		_, err := s.service.RestoreURLs(fileName)
		s.Require().Error(err)
	})

	s.Run("add URL failed", func() {
		s.backupStorage.SetGetAllResponse([]entity.URL{
			{
				UUID:     uuid.UUID{},
				Short:    "",
				Original: "",
			},
		}, nil)
		s.mainStorage.SetAddResponse(nil, ErrEmpty)
		_, err := s.service.RestoreURLs(fileName)
		s.Require().Error(err)
	})

	s.Run("add batch", func() {
		req := []request.ShortenBatchRequest{
			{
				CorrelationID: "1",
				OriginalURL:   "aaa",
			},
		}
		expResp := []response.ShortenBatchResponse{
			{
				CorrelationID: "1",
				ShortURL:      "url/*****",
			},
		}
		s.mainStorage.SetAddBatchResponse(1, nil)
		resp, err := s.service.MakeShortURLBatch(req, 5, "url")
		s.Require().Equal(expResp, resp)
		s.Require().NoError(err)
	})
}
