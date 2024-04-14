package service

import (
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

type URLServiceMock struct {
	makeShortURLEntity *entity.URL
	makeShortURLError  error
	getShortURLEntity  *entity.URL
	getShortURLError   error
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
