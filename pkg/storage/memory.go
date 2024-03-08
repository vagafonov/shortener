package storage

import "github.com/vagafonov/shrinkr/pkg/entity"

type MemoryStorage struct {
	storage map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]string),
	}
}

func (s *MemoryStorage) GetByHash(key string) *entity.URL {
	if v, ok := s.storage[key]; ok {
		return &entity.URL{
			Short: key,
			Full:  v,
		}
	}

	return nil
}

func (s *MemoryStorage) GetByValue(val string) *entity.URL {
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

func (s *MemoryStorage) Add(key string, value string) (*entity.URL, error) {
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

func (s *MemoryStorage) GetAll() []entity.URL {
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

func (s *MemoryStorage) Truncate() {
	clear(s.storage)
}
