package application

import (
	"github.com/vagafonov/shrinkr/pkg/entity"
	hash "github.com/vagafonov/shrinkr/pkg/hasher"
	"github.com/vagafonov/shrinkr/pkg/storage"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) MakeShortURL(url string, length int) (*entity.URL, error) {
	if shortURL := s.storage.GetByValue(url); shortURL != nil {
		return shortURL, nil
	}

	h := hash.NewStringHasher().Hash(length)
	shortURL, err := s.storage.Add(h, url)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}

func (s *Service) GetShortURL(url string) *entity.URL {
	return s.storage.GetByHash(url)
}
