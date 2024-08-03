package contract

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
)

// ErrURLAlreadyExists url for already exists.
var ErrURLAlreadyExists = errors.New("url already exists")

// Storage abstract interface for Service.
// TODO Rename.
type Service interface {
	MakeShortURL(ctx context.Context, url string, length int, userID uuid.UUID) (*entity.URL, error)
	MakeShortURLBatch(ctx context.Context, URLs []*entity.URL, baseURL string) ([]response.ShortenBatchResponse, error) //nolint:lll
	GetShortURL(ctx context.Context, url string) (*entity.URL, error)
	RestoreURLs(ctx context.Context, fileName string) (int, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error)
	DeleteUserURLs(ctx context.Context, userID uuid.UUID, shortURLs []string, batchSize int, jobsCount int) error
	GetStat(ctx context.Context) (*entity.Stat, error)
}
