package service

import (
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/request"
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
}

func NewURLServiceMock() contract.Service {
	return &URLServiceMock{}
}

func (s *URLServiceMock) MakeShortURL(url string, length int) (*entity.URL, error) {
	return s.makeShortURLEntity, s.makeShortURLError
}

func (s *URLServiceMock) SetMakeShortURLResult(e *entity.URL, err error) {
	s.makeShortURLEntity = e
	s.makeShortURLError = err
}

func (s *URLServiceMock) GetShortURL(url string) (*entity.URL, error) {
	return s.getShortURLEntity, s.getShortURLError
}

func (s *URLServiceMock) SetGetShortURLResult(e *entity.URL, err error) {
	s.getShortURLEntity = e
	s.getShortURLError = err
}

func (s *URLServiceMock) RestoreURLs(fileName string) (int, error) {
	return 0, nil
}

func (s *URLServiceMock) MakeShortURLBatch(
	req []request.ShortenBatchRequest,
	length int,
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
