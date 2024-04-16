package storage

import (
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

type memoryStorage struct {
	storage map[string]string
}

func NewMemoryStorage() contract.Storage {
	return &memoryStorage{
		storage: make(map[string]string),
	}
}

func (s *memoryStorage) GetByHash(key string) (*entity.URL, error) {
	if v, ok := s.storage[key]; ok {
		return &entity.URL{
			Short:    key,
			Original: v,
		}, nil
	}

	return nil, nil //nolint:nilnil
}

func (s *memoryStorage) GetByURL(val string) (*entity.URL, error) {
	for k, v := range s.storage {
		if val == v {
			return &entity.URL{
				Short:    k,
				Original: v,
			}, nil
		}
	}

	return nil, nil //nolint:nilnil
}

func (s *memoryStorage) Add(hash string, url string) (*entity.URL, error) {
	for k, v := range s.storage {
		if k == hash || v == url {
			return nil, contract.ErrAlreadyExistsInStorage
		}
	}
	s.storage[hash] = url

	return &entity.URL{
		Short:    hash,
		Original: url,
	}, nil
}

func (s *memoryStorage) GetAll() ([]entity.URL, error) {
	res := make([]entity.URL, len(s.storage))
	i := 0
	for k, v := range s.storage {
		res[i] = entity.URL{
			Short:    k,
			Original: v,
		}
		i++
	}

	return res, nil
}

func (s *memoryStorage) AddBatch(b []entity.URL) (int, error) {
	for _, v := range b {
		s.storage[v.Short] = v.Original
	}

	return len(b), nil
}

func (s *memoryStorage) Truncate() {
	clear(s.storage)
}

func (s *memoryStorage) Close() error {
	return nil
}
