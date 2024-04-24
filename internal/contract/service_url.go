package contract

import (
	"context"
	"errors"

	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
)

var ErrURLAlreadyExists = errors.New("url already exists")

// TODO Rename.
type Service interface {
	MakeShortURL(ctx context.Context, url string, length int) (*entity.URL, error)
	MakeShortURLBatch(ctx context.Context, URLs []entity.URL, baseURL string) ([]response.ShortenBatchResponse, error) //nolint:lll
	GetShortURL(ctx context.Context, url string) (*entity.URL, error)
	RestoreURLs(ctx context.Context, fileName string) (int, error)
}
