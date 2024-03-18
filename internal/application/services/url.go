package services

import (
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/pkg/entity"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

type service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) MakeShortURL(url string, length int) (*entity.URL, error) {
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

func (s *service) GetShortURL(url string) *entity.URL {
	return s.storage.GetByHash(url)
}
