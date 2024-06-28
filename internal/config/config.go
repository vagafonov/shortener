package config

import "github.com/rs/zerolog"

const shortURLLength = 8

const (
	ModeProd Mode = "prod"
	ModeDev  Mode = "dev"
	ModeTest Mode = "test"
)

type Mode string

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
	Mode                Mode
}

func NewConfig(
	serverURL string,
	resultURL string,
	fileStoragePath string,
	databaseDSN string,
	cryptoKey []byte,
	deleteURLsBatchSize int,
	deleteURLsJobsCount int,
	mode Mode,
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
		Mode:                mode,
	}
}
