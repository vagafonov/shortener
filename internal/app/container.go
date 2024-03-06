package app

import "github.com/vagafonov/shrinkr/pkg/storage"

type Container struct {
	storage storage.Storage
}

func NewContainer(s storage.Storage) *Container {
	return &Container{
		storage: s,
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.storage != nil {
		return c.storage
	}

	return storage.NewMemoryStorage()
}
