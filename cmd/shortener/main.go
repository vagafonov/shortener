package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/application"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/service"
	"github.com/vagafonov/shortener/internal/storage"
	"github.com/vagafonov/shortener/pkg/hasher"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

type options struct {
	ServerURL       string `env:"SERVER_ADDRESS"`
	ResultURL       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	EnableHTTPS     string `env:"ENABLE_HTTPS"`
	ConfigFile      string `env:"CONFIG_FILE"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}

func main() {
	ctx := context.Background()
	printBuildInfo()
	opt := &options{
		ServerURL:       "",
		ResultURL:       "",
		FileStoragePath: "",
	}
	parseFlags(opt)
	parseEnv(opt)
	cfg := parseConfFile(opt)

	lr := logger.CreateLogger(cfg.LogLevel)
	var strg contract.Storage
	var err error
	var db *sql.DB
	var fstorage contract.Storage

	if cfg.DatabaseDSN != "" {
		db, err = prepareDB(lr, cfg)
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

	hr := hasher.NewRandHasher(hasher.Alphabet)
	cnt := container.NewContainer(
		cfg,
		strg,
		fstorage,
		hr,
		lr,
		db,
	)

	setBackupStorage(cnt, lr)
	setServiceStorage(cnt, lr)
	setHealthCheckService(cnt, lr)

	app := application.NewApplication(cnt)
	err = runServer(ctx, cfg.EnableHTTPS, app)
	if err != nil {
		lr.Err(err).Send()
	}
}

func prepareDB(lr *zerolog.Logger, cfg *config.Config) (*sql.DB, error) {
	db, err := createConnect(cfg.DatabaseDSN)
	if err != nil {
		lr.Err(err).Send()
	}

	return db, runMigrations(db)
}

func setBackupStorage(cnt *container.Container, lr *zerolog.Logger) {
	backupStorage, err := storage.StorageFactory(cnt, "fs")
	if err != nil {
		lr.Err(err).Send()
	}
	cnt.SetBackupStorage(backupStorage)
}

func setHealthCheckService(cnt *container.Container, lr *zerolog.Logger) {
	servHealthcheck, err := service.ServiceHealthCheckFactory(cnt, "real")
	if err != nil {
		lr.Err(err).Send()
	}
	cnt.SetServiceHealthCheck(servHealthcheck)
}

func setServiceStorage(cnt *container.Container, lr *zerolog.Logger) {
	servURL, err := service.ServiceURLFactory(cnt, "real")
	if err != nil {
		lr.Err(err).Send()
	}
	cnt.SetServiceURL(servURL)
}

//nolint:forbidigo
func printBuildInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func runServer(ctx context.Context, isEnableHTTPS bool, app *application.Application) error {
	if isEnableHTTPS {
		return app.ServeHTTPS(ctx)
	}

	return app.Serve(ctx)
}
