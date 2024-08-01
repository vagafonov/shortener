package config

import "github.com/rs/zerolog"

const shortURLLength = 8

// application modes.
const (
	ModeProd Mode = "prod"
	ModeDev  Mode = "dev"
	ModeTest Mode = "test"
)

// Mode application mode.
type Mode string

// Config.
type Config struct {
	ServerURL           string
	ResultURL           string
	ShortURLLength      int
	LogLevel            zerolog.Level
	FileStoragePath     string
	DatabaseDSN         string
	EnableHTTPS         string
	CryptoKey           []byte
	DeleteURLsBatchSize int
	DeleteURLsJobsCount int
	Mode                Mode
}

// Constructor for Config.
func NewConfig(
	serverURL string,
	resultURL string,
	fileStoragePath string,
	databaseDSN string,
	enableHTTPS string,
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
		EnableHTTPS:         enableHTTPS,
		CryptoKey:           cryptoKey,
		DeleteURLsBatchSize: deleteURLsBatchSize,
		DeleteURLsJobsCount: deleteURLsJobsCount,
		Mode:                mode,
	}
}
