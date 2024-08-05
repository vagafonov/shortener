package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/vagafonov/shortener/internal/config"
)

//nolint:tagliatelle
type conf struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
	Protocol        string `json:"protocol"`
}

func parseConfFile(opt *options) *config.Config { //nolint:cyclop
	if opt == nil {
		log.Fatal("no configuration file provided")
	}

	var enableHTTPS bool
	var err error
	if opt.EnableHTTPS != "" {
		if enableHTTPS, err = strconv.ParseBool(opt.EnableHTTPS); err != nil {
			log.Fatal(err)
		}
	}
	cfg := config.NewConfig(
		opt.ServerURL,
		opt.ResultURL,
		opt.FileStoragePath,
		opt.DatabaseDSN,
		enableHTTPS,
		[]byte("0123456789abcdef"),
		10, //nolint:mnd,gomnd
		2,  //nolint:mnd,gomnd
		config.ModeDev,
		opt.TrustedSubnet,
		opt.Protocol,
	)

	if opt.ConfigFile == "" {
		return cfg
	}

	fBytes, err := os.ReadFile(opt.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := &conf{}
	if err := json.Unmarshal(fBytes, conf); err != nil {
		log.Fatal(err)
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = conf.ServerAddress
	}
	if cfg.ResultURL == "" {
		cfg.ResultURL = conf.BaseURL
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = conf.FileStoragePath
	}
	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = conf.DatabaseDSN
	}
	if !cfg.EnableHTTPS {
		cfg.EnableHTTPS = conf.EnableHTTPS
	}

	return cfg
}
