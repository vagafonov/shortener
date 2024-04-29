package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/pkg/entity"
)

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
}

func NewMemoryStorageMock() contract.Storage {
	return &MemoryStorageMock{}
}

func (s *MemoryStorageMock) GetByHash(ctx context.Context, hash string) (*entity.URL, error) {
	return s.getByHashResponseEntity, s.getByHashResponseError
}

func (s *MemoryStorageMock) SetGetByHashResponse(e *entity.URL, err error) {
	s.getByHashResponseEntity = e
	s.getByHashResponseError = err
}

func (s *MemoryStorageMock) GetByURL(ctx context.Context, url string) (*entity.URL, error) {
	return s.getByURLResponseEntity, s.getByURLResponseError
}

func (s *MemoryStorageMock) SetGetByURLResponse(e *entity.URL, err error) {
	s.getByURLResponseEntity = e
	s.getByURLResponseError = err
}

func (s *MemoryStorageMock) Add(ctx context.Context, hash string, url string, userID uuid.UUID) (*entity.URL, error) {
	return s.addResponseEntity, s.addResponseError
}

func (s *MemoryStorageMock) SetAddResponse(e *entity.URL, err error) {
	s.addResponseEntity = e
	s.addResponseError = err
}

func (s *MemoryStorageMock) GetAll(ctx context.Context) ([]*entity.URL, error) {
	return s.getAllResponseEntity, s.getAllResponseError
}

func (s *MemoryStorageMock) AddBatch(ctx context.Context, b []*entity.URL) (int, error) {
	return s.getAddBatchResponseTotalCreated, s.getAddBatchResponseError
}

func (s *MemoryStorageMock) SetAddBatchResponse(totalCreated int, err error) {
	s.getAddBatchResponseTotalCreated = totalCreated
	s.getAddBatchResponseError = err
}

func (s *MemoryStorageMock) GetAllURLsByUser(
	ctx context.Context,
	userID uuid.UUID,
	baseURL string,
) ([]*entity.URL, error) {
	return s.getAllURLsByUserEntity, s.getAllURLsByUserError
}

func (s *MemoryStorageMock) SetGetAllURLsByUserResponse(u []*entity.URL, err error) {
	s.getAllURLsByUserEntity = u
	s.getAllURLsByUserError = err
}

func (s *MemoryStorageMock) Ping(ctx context.Context) error { return nil }

func (s *MemoryStorageMock) Truncate() {
}

func (s *MemoryStorageMock) Close() error {
	return nil
}
