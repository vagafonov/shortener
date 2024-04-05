package application

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/pkg/entity"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

// TODO rename.
type urlService struct {
	logger        zerolog.Logger
	storage       storage.Storage
	backupStorage storage.Storage
	hasher        hash.Hasher
}

func NewService(
	logger zerolog.Logger,
	storage storage.Storage,
	backupStorage storage.Storage,
	hasher hash.Hasher,
) Service {
	return &urlService{
		logger:        logger,
		storage:       storage,
		backupStorage: backupStorage,
		hasher:        hasher,
	}
}

func (s *urlService) MakeShortURL(url string, length int) (*entity.URL, error) {
	shortURL, err := s.storage.GetByURL(url)
	if err != nil {
		return nil, err
	}
	if shortURL != nil {
		return shortURL, nil
	}
	hashShortURL := s.hasher.Hash(length)
	shortURL, err = s.storage.Add(hashShortURL, url)
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

	return s.storage.GetByHash(url)
}

func (s *urlService) RestoreURLs(fileName string) (int, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0666) //nolint:gofumpt, gomnd
	scanner := bufio.NewScanner(f)
	if err != nil {
		return 0, err
	}

	var e entity.URL
	var count int
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return 0, err
		}
		if _, err = s.storage.Add(e.Short, e.Full); err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}
