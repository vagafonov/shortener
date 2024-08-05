package main

import "flag"

func parseFlags(opt *options) {
	flag.StringVar(&opt.ServerURL, "a", "127.0.0.1:8080", "address and port to run server")
	flag.StringVar(&opt.ResultURL, "b", "http://localhost:8080", "address and port for result short url")
	flag.StringVar(&opt.FileStoragePath, "f", "/tmp/short-url-db.json", "file where stored short urls. if not specified, no data is saved") //nolint:lll
	flag.StringVar(&opt.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&opt.EnableHTTPS, "s", "", "enable https")
	flag.StringVar(&opt.ConfigFile, "c", "", "config file path")
	flag.StringVar(&opt.ConfigFile, "t", "", "trusted subnet")
	flag.StringVar(&opt.Protocol, "p", "", "protocol")
	flag.Parse()
}
