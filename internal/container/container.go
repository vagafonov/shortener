package container

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/contract"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

// Container store dependencies.
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

// NewContainer Constructor for Container.
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

// GetConfig return config from container.
func (c *Container) GetConfig() *config.Config {
	return c.cfg
}

// GetMainStorage return main storage from container.
func (c *Container) GetMainStorage() contract.Storage {
	return c.mainStorage
}

// SetMainStorage set main storage to container.
func (c *Container) SetMainStorage(s contract.Storage) {
	c.mainStorage = s
}

// GetBackupStorage return backup storage from container.
func (c *Container) GetBackupStorage() contract.Storage {
	return c.backupStorage
}

// SetBackupStorage set backup storage to container.
func (c *Container) SetBackupStorage(s contract.Storage) {
	c.backupStorage = s
}

// GetHasher return hasher from container.
func (c *Container) GetHasher() hash.Hasher {
	if c.hasher != nil {
		return c.hasher
	}

	return hash.NewRandHasher(hash.Alphabet)
}

// GetLogger return logger from container.
func (c *Container) GetLogger() *zerolog.Logger {
	return c.logger
}

// GetDB return DB instance from container.
func (c *Container) GetDB() *sql.DB {
	return c.db
}

// GetServiceURL return service URL from container.
func (c *Container) GetServiceURL() contract.Service {
	return c.serviceURL
}

// SetServiceURL set service URL to container.
func (c *Container) SetServiceURL(s contract.Service) {
	c.serviceURL = s
}

// GetServiceHealthCheck return service HealthCheck from container.
func (c *Container) GetServiceHealthCheck() contract.ServiceHealthCheck {
	return c.serviceHealthCheck
}

// SetServiceHealthCheck set ServiceHealthCheck to container.
func (c *Container) SetServiceHealthCheck(s contract.ServiceHealthCheck) {
	c.serviceHealthCheck = s
}
