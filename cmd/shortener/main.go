package main

import (
	"github.com/vagafonov/shrinkr/internal/app"
	"github.com/vagafonov/shrinkr/pkg/storage"
)

func main() {
	cnt := app.NewContainer(storage.NewMemoryStorage())
	app := app.NewApplication(cnt)
	app.Serve()
}
