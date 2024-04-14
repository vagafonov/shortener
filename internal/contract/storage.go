package contract

import (
	"errors"

	"github.com/vagafonov/shortener/pkg/entity"
)

var ErrAlreadyExists = errors.New("already exists")

type Storage interface {
	GetByHash(hash string) (*entity.URL, error)
	GetByURL(url string) (*entity.URL, error)
	Add(hash string, url string) (*entity.URL, error)
	GetAll() ([]entity.URL, error) // todo use pointer
	Truncate()
	Close() error
}
