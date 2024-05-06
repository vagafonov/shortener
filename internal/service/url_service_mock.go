package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
)

type URLServiceMock struct {
	makeShortURLEntity        *entity.URL
	makeShortURLError         error
	getShortURLEntity         *entity.URL
	getShortURLError          error
	makeShortURLBatchResponse []response.ShortenBatchResponse
	makeShortURLBatchError    error
	getUserURLsEntities       []*entity.URL
	getUserURLsError          error
}

func NewURLServiceMock() contract.Service {
	return &URLServiceMock{}
}

func (s *URLServiceMock) MakeShortURL(
	ctx context.Context,
	url string,
	length int,
	userID uuid.UUID,
) (*entity.URL, error) {
	return s.makeShortURLEntity, s.makeShortURLError
}

func (s *URLServiceMock) SetMakeShortURLResult(e *entity.URL, err error) {
	s.makeShortURLEntity = e
	s.makeShortURLError = err
}

func (s *URLServiceMock) GetShortURL(ctx context.Context, url string) (*entity.URL, error) {
	return s.getShortURLEntity, s.getShortURLError
}

func (s *URLServiceMock) SetGetShortURLResult(e *entity.URL, err error) {
	s.getShortURLEntity = e
	s.getShortURLError = err
}

func (s *URLServiceMock) RestoreURLs(ctx context.Context, fileName string) (int, error) {
	return 0, nil
}

func (s *URLServiceMock) MakeShortURLBatch(
	ctx context.Context,
	urls []*entity.URL,
	baseURL string,
) (
	[]response.ShortenBatchResponse, error,
) {
	return s.makeShortURLBatchResponse, s.makeShortURLBatchError
}

func (s *URLServiceMock) SetMakeShortURLBatchResult(resp []response.ShortenBatchResponse, err error) {
	s.makeShortURLBatchResponse = resp
	s.makeShortURLBatchError = err
}

func (s *URLServiceMock) GetUserURLs(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	return s.getUserURLsEntities, s.getUserURLsError
}

func (s *URLServiceMock) SetGetUserURLsResult(e []*entity.URL, err error) {
	s.getUserURLsEntities = e
	s.getUserURLsError = err
}
