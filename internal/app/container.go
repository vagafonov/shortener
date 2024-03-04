package app

import "github.com/vagafonov/shrinkr/pkg/storage"

type Container struct {
	storage storage.Storage
}

func NewContainer() *Container {
	return &Container{
		storage: storage.NewMemoryStorage(),
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.storage != nil {
		return c.storage
	}

	return storage.NewMemoryStorage()
}
