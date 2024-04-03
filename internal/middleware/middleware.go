package middleware

import "github.com/rs/zerolog"

type middleware struct {
	logger zerolog.Logger
}

func NewMiddleware(l zerolog.Logger) *middleware {
	return &middleware{logger: l}
}
