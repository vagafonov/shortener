package app

import (
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

func (s *Service) MakeShortURL(url string) string {
	if val := s.storage.GetByHash(url); val != "" {
		return val
	}

	//fmt.Println("creating new url")
	h := hash.NewStringHasher().Hash(8)
	s.storage.Set(h, url)
	return h
}

func (s *Service) GetShortURL(url string) string {
	return s.storage.GetByHash(url)
}
