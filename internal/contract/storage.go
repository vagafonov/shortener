package contract

import (
	"context"

	"github.com/vagafonov/shortener/pkg/entity"
)

type Storage interface {
	GetByHash(ctx context.Context, hash string) (*entity.URL, error)
	GetByURL(ctx context.Context, url string) (*entity.URL, error)
	Add(ctx context.Context, hash string, url string) (*entity.URL, error)
	AddBatch(ctx context.Context, URLs []entity.URL) (int, error)
	GetAll(ctx context.Context) ([]entity.URL, error) // todo use pointer
	Ping(ctx context.Context) error
	Truncate()
	Close() error
}
