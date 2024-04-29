package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/pkg/entity"
)

type Storage interface {
	GetByHash(ctx context.Context, hash string) (*entity.URL, error)
	GetByURL(ctx context.Context, url string) (*entity.URL, error)
	Add(ctx context.Context, hash string, url string, userID uuid.UUID) (*entity.URL, error)
	AddBatch(ctx context.Context, URLs []*entity.URL) (int, error)
	GetAll(ctx context.Context) ([]*entity.URL, error)
	GetAllURLsByUser(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error)
	Ping(ctx context.Context) error
	Truncate()
	Close() error
}
