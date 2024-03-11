package application

import (
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
)

type Container struct {
	cfg *config.Config
	srg storage.Storage
}

func NewContainer(cfg *config.Config, s storage.Storage) *Container {
	return &Container{
		cfg: cfg,
		srg: s,
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.srg != nil {
		return c.srg
	}

	return storage.NewMemoryStorage()
}
