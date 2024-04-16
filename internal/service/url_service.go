package service

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/request"
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

func (s *urlService) MakeShortURL(url string, length int) (*entity.URL, error) {
	shortURL, err := s.mainStorage.GetByURL(url)
	if err != nil {
		return nil, err
	}
	if shortURL != nil {
		return shortURL, nil
	}
	hashShortURL := s.hasher.Hash(length)
	shortURL, err = s.mainStorage.Add(hashShortURL, url)
	if err != nil {
		return nil, err
	}
	_, err = s.backupStorage.Add(hashShortURL, url)
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *urlService) GetShortURL(url string) (*entity.URL, error) {
	s.logger.Info().Str("url", url).Msg("GetShortURL")

	return s.mainStorage.GetByHash(url)
}

func (s *urlService) RestoreURLs(fileName string) (int, error) {
	// TODO need pagination

	URLs, err := s.backupStorage.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to get all URLs: %w", err)
	}

	for _, v := range URLs {
		// TODO handle id
		if _, err = s.mainStorage.Add(v.Short, v.Original); err != nil {
			return 0, fmt.Errorf("failed to add URL: %w", err)
		}
	}

	return len(URLs), err
}

func (s *urlService) MakeShortURLBatch(
	req []request.ShortenBatchRequest,
	length int,
	baseURL string,
) ([]response.ShortenBatchResponse, error) {
	URLs := make([]entity.URL, len(req))
	for k, v := range req {
		URLs[k] = entity.URL{
			Short:    s.hasher.Hash(length),
			Original: v.OriginalURL,
		}
	}

	totalCreated, err := s.mainStorage.AddBatch(URLs)
	if err != nil {
		return nil, fmt.Errorf("cannot add batch to main storage: %w", err)
	}

	resp := make([]response.ShortenBatchResponse, totalCreated)

	for k, v := range URLs {
		resp[k] = response.ShortenBatchResponse{
			CorrelationID: req[k].CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", baseURL, v.Short),
		}
	}

	_, err = s.backupStorage.AddBatch(URLs)
	if err != nil {
		return nil, fmt.Errorf("cannot add batch to backup storage: %w", err)
	}

	return resp, nil
}
