package service

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
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
	cfg := config.NewConfig(
		"test",
		"http://test:8080",
		fileName,
		"test",
		false,
		[]byte("0123456789abcdef"),
		10,
		3,
		config.ModeTest,
		"",
		string(config.ProtocolHTTP),
	)
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
	ctx := context.Background()
	s.Run("get short url successfully", func() {
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "****",
			Original: "some_url",
		}
		m.EXPECT().GetByHash(ctx, "some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)
		e, err := s.service.GetShortURL(ctx, "some_url")
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceDBSuite) TestMakeShortURL() {
	s.Run("make short url with empty storage", func() {
		ctx := context.Background()
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
		}
		m.EXPECT().GetByURL(ctx, "some_url").Return(nil, nil)
		userID := uuid.Must(uuid.NewUUID())
		m.EXPECT().Add(ctx, "*****", "some_url", userID).Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		e, err := s.service.MakeShortURL(ctx, "some_url", 5, userID)
		s.Require().NoError(err)
		s.Require().Equal(expEntity, e)
	})

	s.Run("make short url with already exist urls", func() {
		ctx := context.Background()
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
		}
		m.EXPECT().GetByURL(ctx, "some_url").Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		userID := uuid.Must(uuid.NewUUID())
		e, err := s.service.MakeShortURL(ctx, "some_url", 5, userID)
		s.Require().Error(err)
		s.Require().Equal(expEntity, e)
	})
}

func (s *ServiceDBSuite) TestRestoreURLs() {
	ctx := context.Background()
	s.Run("restore all urls", func() {
		userID := uuid.New()
		expEntity := &entity.URL{
			UUID:     uuid.UUID{},
			Short:    "*****",
			Original: "some_url",
			UserID:   userID,
		}
		s.backupStorage.SetGetAllResponse([]*entity.URL{expEntity}, nil)
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)

		m.EXPECT().Add(ctx, "*****", "some_url", userID).Return(expEntity, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		totalRestored, err := s.service.RestoreURLs(ctx, fileName)
		s.Require().NoError(err)
		s.Require().Equal(1, totalRestored)
	})
}

func (s *ServiceDBSuite) TestAddBatch() {
	ctx := context.Background()
	s.Run("add batch", func() {
		ctrl := gomock.NewController(s.T())
		defer ctrl.Finish()
		m := storage.NewMockStorage(ctrl)

		newEntities := []*entity.URL{
			{
				ID:       "1",
				Short:    "*****",
				Original: "aaa",
			},
			{
				ID:       "2",
				Short:    "*****",
				Original: "bbb",
			},
		}

		m.EXPECT().AddBatch(ctx, newEntities).Return(2, nil)
		s.service = NewURLService(
			s.cnt.GetLogger(),
			m,
			s.cnt.GetBackupStorage(),
			s.cnt.GetHasher(),
		)

		req := []request.ShortenBatchRequest{
			{
				CorrelationID: "1",
				OriginalURL:   "aaa",
			},
			{
				CorrelationID: "2",
				OriginalURL:   "bbb",
			},
		}

		URLs := make([]*entity.URL, len(req))
		for k, v := range req {
			URLs[k] = &entity.URL{
				ID:       v.CorrelationID,
				Short:    s.cnt.GetHasher().Hash(5),
				Original: v.OriginalURL,
			}
		}
		resp, err := s.service.MakeShortURLBatch(ctx, URLs, "url")
		s.Require().NoError(err)
		respExp := []response.ShortenBatchResponse{
			{
				CorrelationID: "1",
				ShortURL:      "url/*****",
			},
			{
				CorrelationID: "2",
				ShortURL:      "url/*****",
			},
		}
		s.Require().Equal(respExp, resp)
	})
}
