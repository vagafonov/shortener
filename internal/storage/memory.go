package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/pkg/entity"
)

type memoryStorage struct {
	storage map[string]*entity.URL
}

func NewMemoryStorage() contract.Storage {
	return &memoryStorage{
		storage: make(map[string]*entity.URL),
	}
}

func (s *memoryStorage) GetByHash(ctx context.Context, key string) (*entity.URL, error) {
	if v, ok := s.storage[key]; ok {
		return v, nil
	}

	return nil, nil //nolint:nilnil
}

func (s *memoryStorage) GetByURL(ctx context.Context, val string) (*entity.URL, error) {
	for _, v := range s.storage {
		if val == v.Original {
			return v, nil
		}
	}

	return nil, nil //nolint:nilnil
}

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

func (s *memoryStorage) GetAll(ctx context.Context) ([]*entity.URL, error) {
	res := make([]*entity.URL, len(s.storage))
	i := 0
	for _, v := range s.storage {
		res[i] = v
		i++
	}

	return res, nil
}

func (s *memoryStorage) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	for _, v := range b {
		s.storage[v.Short] = v
	}

	return len(b), nil
}

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

func (s *memoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *memoryStorage) Truncate() {
	clear(s.storage)
}

func (s *memoryStorage) Close() error {
	return nil
}
