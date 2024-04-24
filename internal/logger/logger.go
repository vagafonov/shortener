package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func CreateLogger(l zerolog.Level) *zerolog.Logger {
	// Инициализация логгера zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// human-friendly и цветной output
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:exhaustruct
	// Уровень логирования
	zerolog.SetGlobalLevel(l)

	return &logger
}
