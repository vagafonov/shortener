package application

import (
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

type Container struct {
	cfg    *config.Config
	srg    storage.Storage
	hasher hash.Hasher
}

func NewContainer(
	cfg *config.Config,
	s storage.Storage,
	hasher hash.Hasher,
) *Container {
	return &Container{
		cfg:    cfg,
		srg:    s,
		hasher: hasher,
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.srg != nil {
		return c.srg
	}

	return storage.NewMemoryStorage()
}

func (c *Container) GetHasher() hash.Hasher {
	if c.hasher != nil {
		return c.hasher
	}

	return hash.NewRandHasher()
}
