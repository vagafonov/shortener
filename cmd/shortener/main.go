package main

import (
	"log"

	"github.com/vagafonov/shortener/internal/application"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
	hasher "github.com/vagafonov/shortener/pkg/hasher"
)

type options struct {
	ServerURL string `env:"SERVER_ADDRESS"`
	ResultURL string `env:"BASE_URL"`
}

func main() {
	opt := &options{
		ServerURL: "",
		ResultURL: "",
	}
	parseFlags(opt)
	parseEnv(opt)

	cnt := application.NewContainer(
		config.NewConfig(opt.ServerURL, opt.ResultURL),
		storage.NewMemoryStorage(),
		hasher.NewRandHasher(),
	)
	app := application.NewApplication(cnt)
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
