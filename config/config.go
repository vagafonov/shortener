package config

type Config struct {
	ServerURL      string
	ResultURL      string
	ShortURLLength int
}

func NewConfig(serverURL string, resultURL string) *Config {
	return &Config{
		ServerURL:      serverURL,
		ResultURL:      resultURL,
		ShortURLLength: 8,
	}
}
