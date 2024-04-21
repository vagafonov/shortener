package contract

import (
	"context"
	"errors"

	"github.com/vagafonov/shortener/pkg/entity"
)

var (
	ErrAlreadyExistsInStorage = errors.New("already exists")
	ErrURLNotAdded            = errors.New("url not added")
)

type Storage interface {
	GetByHash(hash string) (*entity.URL, error)
	GetByURL(url string) (*entity.URL, error)
	Add(hash string, url string) (*entity.URL, error)
	AddBatch(URLs []entity.URL) (int, error)
	GetAll() ([]entity.URL, error) // todo use pointer
	Ping(ctx context.Context) error
	Truncate()
	Close() error
}
