package main

import "flag"

func parseFlags(opt *options) {
	flag.StringVar(&opt.ServerURL, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&opt.ResultURL, "b", "http://localhost:8080", "address and port for result short url")
	flag.Parse()
}
