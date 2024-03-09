package application

import (
	"github.com/vagafonov/shrinkr/config"
	"github.com/vagafonov/shrinkr/pkg/storage"
)

type Container struct {
	config  *config.Config
	storage storage.Storage
}

func NewContainer(cfg *config.Config, s storage.Storage) *Container {
	return &Container{
		config:  cfg,
		storage: s,
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.storage != nil {
		return c.storage
	}

	return storage.NewMemoryStorage()
}
