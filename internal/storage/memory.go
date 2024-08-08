package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/pkg/entity"
)

// memoryStorage store URLs in memory.
type memoryStorage struct {
	storage map[string]*entity.URL
}

// NewMemoryStorage Constructor for MemoryStorage.
func NewMemoryStorage() contract.Storage {
	return &memoryStorage{
		storage: make(map[string]*entity.URL),
	}
}

// GetByHash get short URLs by hash.
func (s *memoryStorage) GetByHash(ctx context.Context, key string) (*entity.URL, error) {
	if v, ok := s.storage[key]; ok {
		return v, nil
	}

	return nil, nil //nolint:nilnil
}

// GetByURL get short URLs by url.
func (s *memoryStorage) GetByURL(ctx context.Context, val string) (*entity.URL, error) {
	for _, v := range s.storage {
		if val == v.Original {
			return v, nil
		}
	}

	return nil, nil //nolint:nilnil
}

// Add create new short URL in memory.
func (s *memoryStorage) Add(ctx context.Context, hash string, url string, userID uuid.UUID) (*entity.URL, error) {
	for k, v := range s.storage {
		if k == hash || v.Original == url {
			return nil, customerror.ErrAlreadyExistsInStorage
		}
	}
	s.storage[hash] = &entity.URL{
		ID:       "",
		UUID:     uuid.Must(uuid.NewUUID()),
		Short:    hash,
		Original: url,
		UserID:   userID,
	}

	return s.storage[hash], nil
}

// GetAll get all short urls from memory.
func (s *memoryStorage) GetAll(ctx context.Context) ([]*entity.URL, error) {
	res := make([]*entity.URL, len(s.storage))
	i := 0
	for _, v := range s.storage {
		res[i] = v
		i++
	}

	return res, nil
}

// AddBatch add multiple short URLs.
func (s *memoryStorage) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	for _, v := range b {
		s.storage[v.Short] = v
	}

	return len(b), nil
}

// GetAllURLsByUser get all URLs ny user.
func (s *memoryStorage) GetAllURLsByUser(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	res := make([]*entity.URL, 0)
	for _, v := range s.storage {
		if v.UserID != userID {
			continue
		}
		res = append(res, v)
	}

	return res, nil
}

// DeleteURLsByUser not implemented.
func (s *memoryStorage) DeleteURLsByUser(ctx context.Context, userID uuid.UUID, batch []string) error {
	return nil
}

// GetStat TODO need implement.
func (s *memoryStorage) GetStat(ctx context.Context) (*entity.Stat, error) {
	return nil, nil //nolint:nilnil
}

// Ping not implemented.
func (s *memoryStorage) Ping(ctx context.Context) error {
	return nil
}

// Truncate clear memory storage.
func (s *memoryStorage) Truncate() {
	clear(s.storage)
}

// Close not implemented.
func (s *memoryStorage) Close() error {
	return nil
}
