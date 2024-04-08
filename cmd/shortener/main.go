package main

import (
	"log"

	"github.com/vagafonov/shortener/internal/application"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
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

	fss, err := storage.NewFileSystemStorage(opt.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	cnt := application.NewContainer(
		config.NewConfig(opt.ServerURL, opt.ResultURL, opt.FileStoragePath, opt.DatabaseDSN),
		storage.NewMemoryStorage(),
		fss,
		hasher.NewRandHasher(),
	)
	app := application.NewApplication(cnt)
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
