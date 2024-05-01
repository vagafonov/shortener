package config

import "github.com/rs/zerolog"

const shortURLLength = 8

type Config struct {
	ServerURL           string
	ResultURL           string
	ShortURLLength      int
	LogLevel            zerolog.Level
	FileStoragePath     string
	DatabaseDSN         string
	CryptoKey           []byte
	DeleteURLsBatchSize int
	DeleteURLsJobsCount int
}

func NewConfig(
	serverURL string,
	resultURL string,
	fileStoragePath string,
	databaseDSN string,
	cryptoKey []byte,
	deleteURLsBatchSize int,
	deleteURLsJobsCount int,
) *Config {
	return &Config{
		ServerURL:           serverURL,
		ResultURL:           resultURL,
		ShortURLLength:      shortURLLength,
		LogLevel:            zerolog.DebugLevel,
		FileStoragePath:     fileStoragePath,
		DatabaseDSN:         databaseDSN,
		CryptoKey:           cryptoKey,
		DeleteURLsBatchSize: deleteURLsBatchSize,
		DeleteURLsJobsCount: deleteURLsJobsCount,
	}
}
