package config

import "github.com/rs/zerolog"

const shortURLLength = 8

type Config struct {
	ServerURL      string
	ResultURL      string
	ShortURLLength int
	LogLevel       zerolog.Level
}

func NewConfig(serverURL string, resultURL string) *Config {
	return &Config{
		ServerURL:      serverURL,
		ResultURL:      resultURL,
		ShortURLLength: shortURLLength,
		LogLevel:       zerolog.DebugLevel,
	}
}
