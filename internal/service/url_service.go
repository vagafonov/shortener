package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

// TODO rename.
type urlService struct {
	logger        *zerolog.Logger
	mainStorage   contract.Storage
	backupStorage contract.Storage
	hasher        hash.Hasher
}

func NewURLService(
	logger *zerolog.Logger,
	mainStorage contract.Storage,
	backupStorage contract.Storage,
	hasher hash.Hasher,
) contract.Service {
	return &urlService{
		logger:        logger,
		mainStorage:   mainStorage,
		backupStorage: backupStorage,
		hasher:        hasher,
	}
}

func (s *urlService) MakeShortURL(ctx context.Context, url string, length int, userID uuid.UUID) (*entity.URL, error) {
	shortURL, err := s.mainStorage.GetByURL(ctx, url)
	if err != nil {
		return nil, err
	}
	if shortURL != nil {
		return shortURL, customerror.ErrURLAlreadyExists
	}
	hashShortURL := s.hasher.Hash(length)
	shortURL, err = s.mainStorage.Add(ctx, hashShortURL, url, userID)
	if err != nil {
		return nil, err
	}
	_, err = s.backupStorage.Add(ctx, hashShortURL, url, userID)
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *urlService) GetShortURL(ctx context.Context, url string) (*entity.URL, error) {
	s.logger.Info().Str("url", url).Msg("GetShortURL")

	return s.mainStorage.GetByHash(ctx, url)
}

func (s *urlService) RestoreURLs(ctx context.Context, fileName string) (int, error) {
	// TODO need pagination

	URLs, err := s.backupStorage.GetAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get all URLs: %w", err)
	}

	for _, v := range URLs {
		// TODO handle id
		if _, err = s.mainStorage.Add(ctx, v.Short, v.Original, v.UserID); err != nil {
			return 0, fmt.Errorf("failed to add URL: %w", err)
		}
	}

	return len(URLs), err
}

func (s *urlService) MakeShortURLBatch(
	ctx context.Context,
	urls []*entity.URL,
	baseURL string,
) ([]response.ShortenBatchResponse, error) {
	totalCreated, err := s.mainStorage.AddBatch(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf("cannot add batch to main storage: %w", err)
	}

	resp := make([]response.ShortenBatchResponse, totalCreated)

	for k, v := range urls {
		resp[k] = response.ShortenBatchResponse{
			CorrelationID: v.ID,
			ShortURL:      fmt.Sprintf("%s/%s", baseURL, v.Short),
		}
	}

	_, err = s.backupStorage.AddBatch(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf("cannot add batch to backup storage: %w", err)
	}

	return resp, nil
}

func (s *urlService) GetUserURLs(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	return s.mainStorage.GetAllURLsByUser(ctx, userID, baseURL)
}
