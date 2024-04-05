package storage

import (
	"github.com/vagafonov/shortener/pkg/entity"
)

type memoryStorage struct {
	storage map[string]string
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		storage: make(map[string]string),
	}
}

func (s *memoryStorage) GetByHash(key string) (*entity.URL, error) {
	if v, ok := s.storage[key]; ok {
		return &entity.URL{
			Short: key,
			Full:  v,
		}, nil
	}

	return nil, nil //nolint:nilnil
}

func (s *memoryStorage) GetByURL(val string) (*entity.URL, error) {
	for k, v := range s.storage {
		if val == v {
			return &entity.URL{
				Short: k,
				Full:  v,
			}, nil
		}
	}

	return nil, nil //nolint:nilnil
}

func (s *memoryStorage) Add(hash string, url string) (*entity.URL, error) {
	for k, v := range s.storage {
		if k == hash || v == url {
			return nil, ErrAlreadyExists
		}
	}
	s.storage[hash] = url

	return &entity.URL{
		Short: hash,
		Full:  url,
	}, nil
}

func (s *memoryStorage) GetAll() ([]entity.URL, error) {
	res := make([]entity.URL, len(s.storage))
	i := 0
	for k, v := range s.storage {
		res[i] = entity.URL{
			Short: k,
			Full:  v,
		}
		i++
	}

	return res, nil
}

func (s *memoryStorage) Truncate() {
	clear(s.storage)
}

func (s *memoryStorage) Close() error {
	return nil
}
