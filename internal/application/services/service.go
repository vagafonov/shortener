package services

import (
	"github.com/vagafonov/shortener/pkg/entity"
)

type Service interface {
	MakeShortURL(url string, length int) (*entity.URL, error)
	GetShortURL(url string) *entity.URL
}
