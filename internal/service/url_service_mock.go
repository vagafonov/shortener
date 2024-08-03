package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/pkg/entity"
)

// URLServiceMock mock.
type URLServiceMock struct {
	makeShortURLEntity        *entity.URL
	makeShortURLError         error
	getShortURLEntity         *entity.URL
	getShortURLError          error
	makeShortURLBatchResponse []response.ShortenBatchResponse
	makeShortURLBatchError    error
	getUserURLsEntities       []*entity.URL
	getUserURLsError          error
	deleteUserURLsError       error
	getStatEntity             *entity.Stat
	getStatError              error
}

// NewURLServiceMock Constructor for URLServiceMock.
func NewURLServiceMock() contract.Service {
	return &URLServiceMock{}
}

// MakeShortURL mock.
func (s *URLServiceMock) MakeShortURL(
	ctx context.Context,
	url string,
	length int,
	userID uuid.UUID,
) (*entity.URL, error) {
	return s.makeShortURLEntity, s.makeShortURLError
}

// SetMakeShortURLResult mock.
func (s *URLServiceMock) SetMakeShortURLResult(e *entity.URL, err error) {
	s.makeShortURLEntity = e
	s.makeShortURLError = err
}

// GetShortURL mock.
func (s *URLServiceMock) GetShortURL(ctx context.Context, url string) (*entity.URL, error) {
	return s.getShortURLEntity, s.getShortURLError
}

// SetGetShortURLResult mock.
func (s *URLServiceMock) SetGetShortURLResult(e *entity.URL, err error) {
	s.getShortURLEntity = e
	s.getShortURLError = err
}

// RestoreURLs mock.
func (s *URLServiceMock) RestoreURLs(ctx context.Context, fileName string) (int, error) {
	return 0, nil
}

// MakeShortURLBatch mock.
func (s *URLServiceMock) MakeShortURLBatch(
	ctx context.Context,
	urls []*entity.URL,
	baseURL string,
) (
	[]response.ShortenBatchResponse, error,
) {
	return s.makeShortURLBatchResponse, s.makeShortURLBatchError
}

// SetMakeShortURLBatchResult mock.
func (s *URLServiceMock) SetMakeShortURLBatchResult(resp []response.ShortenBatchResponse, err error) {
	s.makeShortURLBatchResponse = resp
	s.makeShortURLBatchError = err
}

// GetUserURLs mock.
func (s *URLServiceMock) GetUserURLs(ctx context.Context, userID uuid.UUID, baseURL string) ([]*entity.URL, error) {
	return s.getUserURLsEntities, s.getUserURLsError
}

// SetGetUserURLsResult mock.
func (s *URLServiceMock) SetGetUserURLsResult(e []*entity.URL, err error) {
	s.getUserURLsEntities = e
	s.getUserURLsError = err
}

// DeleteUserURLs mock.
func (s *URLServiceMock) DeleteUserURLs(
	ctx context.Context,
	userID uuid.UUID,
	shortURLs []string,
	batchSize int,
	jobsCount int,
) error {
	return s.deleteUserURLsError
}

// SetDeleteUserURLsResult mock.
func (s *URLServiceMock) SetDeleteUserURLsResult(err error) {
	s.deleteUserURLsError = err
}

// DeleteUserURLs mock.
func (s *URLServiceMock) GetStat(ctx context.Context) (*entity.Stat, error) {
	return s.getStatEntity, s.getStatError
}

// SetDeleteUserURLsResult mock.
func (s *URLServiceMock) SetGetStatResult(stat *entity.Stat, err error) {
	s.getStatEntity = stat
	s.getStatError = err
}
