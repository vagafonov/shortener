package application

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

type Container struct {
	cfg     *config.Config
	strg    storage.Storage
	bcpStrg storage.Storage
	hasher  hash.Hasher
	logger  zerolog.Logger
	db      *sql.DB
}

func NewContainer(
	cfg *config.Config,
	strg storage.Storage,
	bcpStrg storage.Storage,
	hasher hash.Hasher,
) *Container {
	// Инициализация логгера zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// human-friendly и цветной output
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:exhaustruct
	// Уровень логирования
	zerolog.SetGlobalLevel(cfg.LogLevel)
	logger.Info().Msgf("DSN: %v", cfg.DatabaseDSN)
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	return &Container{
		cfg:     cfg,
		strg:    strg,
		bcpStrg: bcpStrg,
		hasher:  hasher,
		logger:  logger,
		db:      db,
	}
}

func (c *Container) GetStorage() storage.Storage {
	if c.strg != nil {
		return c.strg
	}

	return storage.NewMemoryStorage()
}

func (c *Container) GetBackupStorage() storage.Storage {
	return c.bcpStrg
}

func (c *Container) GetHasher() hash.Hasher {
	if c.hasher != nil {
		return c.hasher
	}

	return hash.NewRandHasher()
}

func (c *Container) GetLogger() zerolog.Logger {
	return c.logger
}
