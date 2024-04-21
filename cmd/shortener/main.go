package main

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/vagafonov/shortener/internal/application"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/service"
	"github.com/vagafonov/shortener/internal/storage"
	hasher "github.com/vagafonov/shortener/pkg/hasher"
)

type options struct {
	ServerURL       string `env:"SERVER_ADDRESS"`
	ResultURL       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

//nolint:funlen
func main() {
	opt := &options{
		ServerURL:       "",
		ResultURL:       "",
		FileStoragePath: "",
	}
	parseFlags(opt)
	parseEnv(opt)

	cfg := config.NewConfig(opt.ServerURL, opt.ResultURL, opt.FileStoragePath, opt.DatabaseDSN)
	lr := logger.CreateLogger(cfg.LogLevel)
	var strg contract.Storage
	var err error
	var db *sql.DB
	var fstorage contract.Storage

	if cfg.DatabaseDSN != "" {
		if db, err = createConnect(cfg.DatabaseDSN); err != nil {
			lr.Err(err).Send()
		}
		err := runMigrations(db)

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			lr.Err(err).Send()
		}
		strg = storage.NewDBStorage(db)
	} else {
		strg = storage.NewMemoryStorage()
		if fstorage, err = storage.NewFileSystemStorage(cfg.FileStoragePath); err != nil {
			lr.Err(err).Send()
		}
	}

	hr := hasher.NewRandHasher()
	cnt := container.NewContainer(
		cfg,
		strg,
		fstorage,
		hr,
		lr,
		db,
	)

	backupStorage, err := storage.StorageFactory(cnt, "fs")
	if err != nil {
		lr.Err(err).Send()
	}

	cnt.SetBackupStorage(backupStorage)

	servURL, err := service.ServiceURLFactory(cnt, "real")
	if err != nil {
		lr.Err(err).Send()
	}
	cnt.SetServiceURL(servURL)

	servHealthcheck, err := service.ServiceHealthCheckFactory(cnt, "real")
	if err != nil {
		lr.Err(err).Send()
	}
	cnt.SetServiceHealthCheck(servHealthcheck)

	app := application.NewApplication(cnt)
	if err := app.Serve(); err != nil {
		lr.Err(err).Send()
	}
}
