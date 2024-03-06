package storage

type MemoryStorage struct {
	storage map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]string),
	}
}

func (s *MemoryStorage) GetByHash(key string) string {
	if v, ok := s.storage[key]; ok {
		return v
	}

	return ""
}

func (s *MemoryStorage) GetByValue(key string) string {
	for k, v := range s.storage {
		if key == v {
			return k
		}
	}
	return ""
}

func (s *MemoryStorage) Set(key string, value string) error {
	if v := s.GetByHash(key); v != "" {
		return ErrAlreadyExists
	}

	s.storage[key] = value
	return nil
}
