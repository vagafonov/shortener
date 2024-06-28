package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

// MemoryStorageMock.
type MemoryStorageMock struct {
	getAllResponseEntity []*entity.URL
	getAllResponseError  error

	addResponseEntity *entity.URL
	addResponseError  error

	getByURLResponseEntity *entity.URL
	getByURLResponseError  error

	getByHashResponseEntity *entity.URL
	getByHashResponseError  error

	getAddBatchResponseTotalCreated int
	getAddBatchResponseError        error

	getAllURLsByUserEntity []*entity.URL
	getAllURLsByUserError  error

	deleteURLsByUserError error
}

// Constructor for MemoryStorageMock.
func NewMemoryStorageMock() contract.Storage {
	return &MemoryStorageMock{}
}

// GetByHash.
func (s *MemoryStorageMock) GetByHash(ctx context.Context, hash string) (*entity.URL, error) {
	return s.getByHashResponseEntity, s.getByHashResponseError
}

// SetGetByHashResponse.
func (s *MemoryStorageMock) SetGetByHashResponse(e *entity.URL, err error) {
	s.getByHashResponseEntity = e
	s.getByHashResponseError = err
}

// GetByURL.
func (s *MemoryStorageMock) GetByURL(ctx context.Context, url string) (*entity.URL, error) {
	return s.getByURLResponseEntity, s.getByURLResponseError
}

// SetGetByURLResponse.
func (s *MemoryStorageMock) SetGetByURLResponse(e *entity.URL, err error) {
	s.getByURLResponseEntity = e
	s.getByURLResponseError = err
}

// Add.
func (s *MemoryStorageMock) Add(ctx context.Context, hash string, url string, userID uuid.UUID) (*entity.URL, error) {
	return s.addResponseEntity, s.addResponseError
}

// SetAddResponse.
func (s *MemoryStorageMock) SetAddResponse(e *entity.URL, err error) {
	s.addResponseEntity = e
	s.addResponseError = err
}

// GetAll.
func (s *MemoryStorageMock) GetAll(ctx context.Context) ([]*entity.URL, error) {
	return s.getAllResponseEntity, s.getAllResponseError
}

// AddBatch.
func (s *MemoryStorageMock) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	return s.getAddBatchResponseTotalCreated, s.getAddBatchResponseError
}

// SetAddBatchResponse.
func (s *MemoryStorageMock) SetAddBatchResponse(totalCreated int, err error) {
	s.getAddBatchResponseTotalCreated = totalCreated
	s.getAddBatchResponseError = err
}

// GetAllURLsByUser.
func (s *MemoryStorageMock) GetAllURLsByUser(
	ctx context.Context,
	userID uuid.UUID,
	baseURL string,
) ([]*entity.URL, error) {
	return s.getAllURLsByUserEntity, s.getAllURLsByUserError
}

// SetGetAllURLsByUserResponse.
func (s *MemoryStorageMock) SetGetAllURLsByUserResponse(u []*entity.URL, err error) {
	s.getAllURLsByUserEntity = u
	s.getAllURLsByUserError = err
}

// DeleteURLsByUser.
func (s *MemoryStorageMock) DeleteURLsByUser(ctx context.Context, userID uuid.UUID, batch []string) error {
	return nil
}

// SetDeleteURLsByUserResponse.
func (s *MemoryStorageMock) SetDeleteURLsByUserResponse(err error) {
	s.deleteURLsByUserError = err
}

// Ping.
func (s *MemoryStorageMock) Ping(ctx context.Context) error { return nil }

// Truncate.
func (s *MemoryStorageMock) Truncate() {
}

// Close.
func (s *MemoryStorageMock) Close() error {
	return nil
}
