package contract

import (
	"github.com/vagafonov/shortener/internal/request"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
)

type Service interface {
	MakeShortURL(url string, length int) (*entity.URL, error)
	MakeShortURLBatch(req []request.ShortenBatchRequest, length int, baseURL string) ([]response.ShortenBatchResponse, error) //nolint:lll
	GetShortURL(url string) (*entity.URL, error)
	RestoreURLs(fileName string) (int, error)
}
