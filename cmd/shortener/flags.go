package main

import (
	"flag"
)

var flags struct {
	serverURL string
	resultURL string
}

func parseFlags() {
	flag.StringVar(&flags.serverURL, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flags.resultURL, "b", "http://localhost:8080", "address and port for result short url")
	flag.Parse()
}
