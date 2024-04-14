package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func main() {
	opt := &options{
		ServerURL:       "",
		ResultURL:       "",
		FileStoragePath: "",
	}
	parseFlags(opt)
	parseEnv(opt)

	cfg := config.NewConfig(opt.ServerURL, opt.ResultURL, opt.FileStoragePath, opt.DatabaseDSN)
	var strg contract.Storage
	var err error
	var db *sql.DB
	var fstorage contract.Storage

	if cfg.DatabaseDSN != "" {
		if db, err = createConnect(cfg.DatabaseDSN); err != nil {
			log.Fatal(err)
		}
		err := runMigrations(db)

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(fmt.Errorf("cannot run migrations: %w", err))
		}
		strg = storage.NewDBStorage(db)
	} else {
		strg = storage.NewMemoryStorage()
		if fstorage, err = storage.NewFileSystemStorage(cfg.FileStoragePath); err != nil {
			log.Fatal(err)
		}
	}

	lr := logger.CreateLogger(cfg.LogLevel)
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
		log.Fatal(err)
	}

	cnt.SetBackupStorage(backupStorage)
	serv, err := service.ServiceURLFactory(cnt, "real")
	if err != nil {
		log.Fatal(err)
	}
	cnt.SetServiceURL(serv)

	app := application.NewApplication(cnt)
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
