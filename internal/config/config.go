package config

import "github.com/rs/zerolog"

const shortURLLength = 8

type Config struct {
	ServerURL       string
	ResultURL       string
	ShortURLLength  int
	LogLevel        zerolog.Level
	FileStoragePath string
	DatabaseDSN     string
	CryptoKey       []byte
}

func NewConfig(
	serverURL string,
	resultURL string,
	fileStoragePath string,
	databaseDSN string,
	cryptoKey []byte,
) *Config {
	return &Config{
		ServerURL:       serverURL,
		ResultURL:       resultURL,
		ShortURLLength:  shortURLLength,
		LogLevel:        zerolog.DebugLevel,
		FileStoragePath: fileStoragePath,
		DatabaseDSN:     databaseDSN,
		CryptoKey:       cryptoKey,
	}
}
