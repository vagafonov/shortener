package service

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/storage"
	"github.com/vagafonov/shortener/pkg/entity"
	hasher "github.com/vagafonov/shortener/pkg/hasher"
)

type ServiceDBSuite struct {
	suite.Suite
	service       contract.Service
	cnt           *container.Container
	backupStorage *storage.FileSystemStorageMock
}

func TestServiceDBSuite(t *testing.T) {
	suite.Run(t, new(ServiceDBSuite))
}

func (s *ServiceDBSuite) SetupSuite() {
	fss, err := storage.NewFileSystemStorage(fileName)
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig("test", "http://test:8080", fileName, "test")
	lr := logger.CreateLogger(cfg.LogLevel)
	hr := hasher.NewMockHasher()
	s.cnt = container.NewContainer(
		cfg,
		nil,
		fss,
		hr,
		lr,
		nil,
	)

	backupStorage, err := storage.StorageFactory(s.cnt, "fs-mock")
	if err != nil {
		log.Fatal(err)
	}
	s.backupStorage, _ = backupStorage.(*storage.FileSystemStorageMock)
	s.cnt.SetBackupStorage(s.backupStorage)
}

func (s *ServiceDBSuite) TestGetShortURL() {
	s.Run("get short url successfully", func() {
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "****",
			Original: "some_url",
		}
		m.EXPECT().GetByHash("some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)
		e, err := s.service.GetShortURL("some_url")
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceDBSuite) TestMakeShortURL() {
	s.Run("make short url with empty storage", func() {
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
		}
		m.EXPECT().GetByURL("some_url").Return(nil, nil)
		m.EXPECT().Add("*****", "some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		e, err := s.service.MakeShortURL("some_url", 5)
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})

	s.Run("make short url with already exist urls", func() {
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
		}
		m.EXPECT().GetByURL("some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		e, err := s.service.MakeShortURL("some_url", 5)
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceDBSuite) TestRestoreURLs() {
	s.Run("restore all urls", func() {
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
		}
		s.backupStorage.SetGetAllResponse([]entity.URL{*expEntity}, nil)
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)

		m.EXPECT().Add("*****", "some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		totalRestored, err := s.service.RestoreURLs(fileName)
		s.Require().NoError(err)
		s.Require().Equal(1, totalRestored)
	})
}
