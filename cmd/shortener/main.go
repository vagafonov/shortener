package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/vagafonov/shrinkr/config"
	"github.com/vagafonov/shrinkr/internal/application"
	"github.com/vagafonov/shrinkr/pkg/storage"
	"log"
)

var options struct {
	ServerURL string `env:"SERVER_ADDRESS"`
	ResultURL string `env:"BASE_URL"`
}

func main() {
	parseFlags()
	parseEnv()
	cnt := application.NewContainer(
		config.NewConfig(options.ServerURL, options.ResultURL),
		storage.NewMemoryStorage(),
	)
	app := application.NewApplication(cnt)
	app.Serve()
}

func parseFlags() {
	flag.StringVar(&options.ServerURL, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&options.ResultURL, "b", "http://localhost:8080", "address and port for result short url")
	flag.Parse()
}

func parseEnv() {
	err := env.Parse(&options)
	fmt.Println(options)
	if err != nil {
		log.Fatal(err)
	}
}
