package main

import (
	"github.com/vagafonov/shrinkr/config"
	"github.com/vagafonov/shrinkr/internal/application"
	"github.com/vagafonov/shrinkr/pkg/storage"
)

func main() {
	parseFlags()
	cnt := application.NewContainer(
		config.NewConfig(flags.serverURL, flags.resultURL),
		storage.NewMemoryStorage(),
	)
	app := application.NewApplication(cnt)
	app.Serve()
}
