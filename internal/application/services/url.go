package services

import (
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/pkg/entity"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

type service struct {
	storage storage.Storage
	hasher  hash.Hasher
}

func NewService(
	storage storage.Storage,
	hasher hash.Hasher,
) Service {
	return &service{
		storage: storage,
		hasher:  hasher,
	}
}

func (s *service) MakeShortURL(url string, length int) (*entity.URL, error) {
	if shortURL := s.storage.GetByValue(url); shortURL != nil {
		return shortURL, nil
	}

	shortURL, err := s.storage.Add(s.hasher.Hash(length), url)
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *service) GetShortURL(url string) *entity.URL {
	return s.storage.GetByHash(url)
}
