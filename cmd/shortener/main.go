package main

import (
	"github.com/vagafonov/shrinkr/internal/app"
)

func main() {
	cnt := app.NewContainer()
	app := app.NewApplication(cnt)
	app.Serve()
}
