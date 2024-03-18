package config

const shortURLLength = 8

type Config struct {
	ServerURL      string
	ResultURL      string
	ShortURLLength int
}

func NewConfig(serverURL string, resultURL string) *Config {
	return &Config{
		ServerURL:      serverURL,
		ResultURL:      resultURL,
		ShortURLLength: shortURLLength,
	}
}
