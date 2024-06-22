package container

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/contract"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

type Container struct {
	cfg                *config.Config
	mainStorage        contract.Storage
	backupStorage      contract.Storage
	hasher             hash.Hasher
	logger             *zerolog.Logger
	db                 *sql.DB
	serviceURL         contract.Service
	serviceHealthCheck contract.ServiceHealthCheck
}

func NewContainer(
	cfg *config.Config,
	mainStorage contract.Storage,
	backupStorage contract.Storage,
	hasher hash.Hasher,
	logger *zerolog.Logger,
	db *sql.DB,
) *Container {
	return &Container{
		cfg:           cfg,
		mainStorage:   mainStorage,
		backupStorage: backupStorage,
		hasher:        hasher,
		logger:        logger,
		db:            db,
	}
}

func (c *Container) GetConfig() *config.Config {
	return c.cfg
}

func (c *Container) GetMainStorage() contract.Storage {
	return c.mainStorage
}

func (c *Container) SetMainStorage(s contract.Storage) {
	c.mainStorage = s
}

func (c *Container) GetBackupStorage() contract.Storage {
	return c.backupStorage
}

func (c *Container) SetBackupStorage(s contract.Storage) {
	c.backupStorage = s
}

func (c *Container) GetHasher() hash.Hasher {
	if c.hasher != nil {
		return c.hasher
	}

	return hash.NewRandHasher(hash.Alphabet)
}

func (c *Container) GetLogger() *zerolog.Logger {
	return c.logger
}

func (c *Container) GetDB() *sql.DB {
	return c.db
}

func (c *Container) GetServiceURL() contract.Service {
	return c.serviceURL
}

func (c *Container) SetServiceURL(s contract.Service) {
	c.serviceURL = s
}

func (c *Container) GetServiceHealthCheck() contract.ServiceHealthCheck {
	return c.serviceHealthCheck
}

func (c *Container) SetServiceHealthCheck(s contract.ServiceHealthCheck) {
	c.serviceHealthCheck = s
}
