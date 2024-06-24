package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

// FileSystemStorageMock mock.
type FileSystemStorageMock struct {
	getAllResponseEntity []*entity.URL
	getAllResponseError  error

	addResponseEntity *entity.URL
	addResponseError  error

	addBatchResponseTotalCreated int
	addBatchResponseError        error

	deleteURLsByUserError error
}

// Constructor for FileSystemStorageMock.
func NewFileSystemStorageMock() contract.Storage {
	return &FileSystemStorageMock{}
}

// GetByHash mock.
func (s *FileSystemStorageMock) GetByHash(ctx context.Context, hash string) (*entity.URL, error) {
	return nil, nil //nolint:nilnil
}

// GetByURL mock.
func (s *FileSystemStorageMock) GetByURL(ctx context.Context, url string) (*entity.URL, error) {
	return nil, nil //nolint:nilnil
}

// Add mock.
func (s *FileSystemStorageMock) Add(
	ctx context.Context,
	hash string,
	url string,
	userID uuid.UUID,
) (*entity.URL, error) {
	return s.addResponseEntity, s.addResponseError
}

// SetAddResponse mock.
func (s *FileSystemStorageMock) SetAddResponse(e *entity.URL, err error) {
	s.addResponseEntity = e
	s.addResponseError = err
}

// GetAll mock.
func (s *FileSystemStorageMock) GetAll(ctx context.Context) ([]*entity.URL, error) {
	return s.getAllResponseEntity, s.getAllResponseError
}

// SetGetAllResponse mock.
func (s *FileSystemStorageMock) SetGetAllResponse(e []*entity.URL, err error) {
	s.getAllResponseEntity = e
	s.getAllResponseError = err
}

// AddBatch mock.
func (s *FileSystemStorageMock) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	return s.addBatchResponseTotalCreated, s.addBatchResponseError
}

// SetAddBatchResponse mock.
func (s *FileSystemStorageMock) SetAddBatchResponse(totalCreated int, err error) {
	s.addBatchResponseTotalCreated = totalCreated
	s.addBatchResponseError = err
}

// GetAllURLsByUser mock.
func (s *FileSystemStorageMock) GetAllURLsByUser(
	ctx context.Context,
	userID uuid.UUID,
	baseURL string,
) ([]*entity.URL, error) {
	return nil, nil
}

// DeleteURLsByUser mock.
func (s *FileSystemStorageMock) DeleteURLsByUser(ctx context.Context, userID uuid.UUID, batch []string) error {
	return nil
}

// SetDeleteURLsByUserResponse mock.
func (s *FileSystemStorageMock) SetDeleteURLsByUserResponse(err error) {
	s.deleteURLsByUserError = err
}

// Ping mock.
func (s *FileSystemStorageMock) Ping(ctx context.Context) error { return nil }

// Truncate mock.
func (s *FileSystemStorageMock) Truncate() {
}

// Close mock.
func (s *FileSystemStorageMock) Close() error {
	return nil
}
