package middleware

import "github.com/rs/zerolog"

type middleware struct {
	logger *zerolog.Logger
}

// Constructor for Middleware.
func NewMiddleware(l *zerolog.Logger) *middleware {
	return &middleware{logger: l}
}
