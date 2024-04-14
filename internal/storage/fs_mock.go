package storage

import (
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

type FileSystemStorageMock struct {
	getAllResponseEntity []entity.URL
	getAllResponseError  error

	addResponseEntity *entity.URL
	addResponseError  error
}

func NewFileSystemStorageMock() contract.Storage {
	return &FileSystemStorageMock{}
}

func (s *FileSystemStorageMock) GetByHash(hash string) (*entity.URL, error) {
	return nil, nil //nolint:nilnil
}

func (s *FileSystemStorageMock) GetByURL(url string) (*entity.URL, error) {
	return nil, nil //nolint:nilnil
}

func (s *FileSystemStorageMock) Add(hash string, url string) (*entity.URL, error) {
	return s.addResponseEntity, s.addResponseError
}

func (s *FileSystemStorageMock) SetAddResponse(e *entity.URL, err error) {
	s.addResponseEntity = e
	s.addResponseError = err
}

func (s *FileSystemStorageMock) GetAll() ([]entity.URL, error) {
	return s.getAllResponseEntity, s.getAllResponseError
}

func (s *FileSystemStorageMock) SetGetAllResponse(e []entity.URL, err error) {
	s.getAllResponseEntity = e
	s.getAllResponseError = err
}

func (s *FileSystemStorageMock) Truncate() {
}

func (s *FileSystemStorageMock) Close() error {
	return nil
}
