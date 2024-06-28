package service

import (
	"context"
	"fmt"
	"sync"

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

// NewURLService Constructor for URLService.
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

// MakeShortURL make short url.
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

// GetShortURL get short url.
func (s *urlService) GetShortURL(ctx context.Context, url string) (*entity.URL, error) {
	s.logger.Info().Str("url", url).Msg("GetShortURL")

	return s.mainStorage.GetByHash(ctx, url)
}

// RestoreURLs restore short URLs.
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

// MakeShortURLBatch make short URL batch.
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

// GetUserURLs get user URLS.
func (s *urlService) GetUserURLs(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	return s.mainStorage.GetAllURLsByUser(ctx, userID, baseURL)
}

// DeleteUserURLs delete user URLS.
func (s *urlService) DeleteUserURLs(
	ctx context.Context,
	userID uuid.UUID,
	shortURLs []string,
	batchSize int,
	jobsCount int,
) error {
	ch := s.batchDeleteGenerator(shortURLs, batchSize)
	s.batchDeleteConsumer(ctx, ch, jobsCount, userID)

	return nil
}

func (s *urlService) batchDeleteGenerator(shortURLs []string, bSize int) chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)

		batch := make([]string, bSize)
		var skippedPosition int
		for k, v := range shortURLs {
			if k != 0 && k%bSize == 0 {
				s.logger.Debug().Strs("batch", batch).Msg("generator write batch")
				ch <- batch
				batch = make([]string, bSize)
				skippedPosition = k
			}
			batch[k%bSize] = v
		}
		ch <- shortURLs[skippedPosition:]
	}()

	return ch
}

func (s *urlService) batchDeleteConsumer(ctx context.Context, ch chan []string, jobsCount int, userID uuid.UUID) {
	s.logger.Debug().Msgf("batch delete consumer started with %v jobs", jobsCount)
	wg := sync.WaitGroup{}
	wg.Add(jobsCount)
	for i := 1; i <= jobsCount; i++ {
		go func(n int) {
			s.logger.Debug().Msgf("gorutine №%v started", n)
			for batch := range ch {
				if err := s.mainStorage.DeleteURLsByUser(ctx, userID, batch); err != nil {
					s.logger.Error().Err(err).Strs("batch", batch).Msg("failed delete batch urls in consumer")
				}
				s.logger.Debug().Strs("batch", batch).Msgf("gorutine №%v successfully handled batch in comsumer", n)
			}
			wg.Done()

			s.logger.Debug().Msgf("gorutine №%v completed", n)
		}(i)
	}
	wg.Wait()
	s.logger.Debug().Msg("batch delete consumer completed")
}
