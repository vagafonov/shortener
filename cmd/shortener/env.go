package main

import (
	"log"

	"github.com/caarlos0/env/v6"
)

func parseEnv(opt *options) {
	if err := env.Parse(opt); err != nil {
		log.Fatal(err)
	}
}
