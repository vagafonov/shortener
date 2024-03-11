package storage

import (
	"errors"

	"github.com/vagafonov/shortener/pkg/entity"
)

var ErrAlreadyExists = errors.New("already exists")

type Storage interface {
	GetByHash(key string) *entity.URL
	GetByValue(key string) *entity.URL
	Add(key string, val string) (*entity.URL, error)
	GetAll() []entity.URL
	Truncate()
}
