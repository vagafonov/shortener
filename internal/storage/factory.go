package storage

import (
	"errors"

	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
)

var ErrUndefinedStorageType = errors.New("undefined storage type")

func StorageFactory(cnt *container.Container, t string) (contract.Storage, error) {
	// TODO use enum
	switch t {
	case "db":
		return NewDBStorage(cnt.GetDB()), nil
	case "fs":
		return NewFileSystemStorage(cnt.GetConfig().FileStoragePath)
	case "fs-mock":
		return NewFileSystemStorageMock(), nil
	case "memory-mock":
		return NewMemoryStorageMock(), nil
	default:
		return nil, ErrUndefinedStorageType
	}
}
