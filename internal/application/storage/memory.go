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

func (s *memoryStorage) GetByHash(key string) *entity.URL {
	if v, ok := s.storage[key]; ok {
		return &entity.URL{
			Short: key,
			Full:  v,
		}
	}

	return nil
}

func (s *memoryStorage) GetByValue(val string) *entity.URL {
	for k, v := range s.storage {
		if val == v {
			return &entity.URL{
				Short: k,
				Full:  v,
			}
		}
	}

	return nil
}

func (s *memoryStorage) Add(key string, value string) (*entity.URL, error) {
	if shortURL := s.GetByHash(key); shortURL != nil {
		return nil, ErrAlreadyExists
	}
	if shortURL := s.GetByValue(value); shortURL != nil {
		return nil, ErrAlreadyExists
	}
	s.storage[key] = value

	return &entity.URL{
		Short: key,
		Full:  value,
	}, nil
}

func (s *memoryStorage) GetAll() []entity.URL {
	res := make([]entity.URL, len(s.storage))
	i := 0
	for k, v := range s.storage {
		res[i] = entity.URL{
			Short: k,
			Full:  v,
		}
		i++
	}

	return res
}

func (s *memoryStorage) Truncate() {
	clear(s.storage)
}
