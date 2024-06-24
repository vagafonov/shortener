package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vagafonov/shortener/internal/contract"
)

// TODO rename.
type healtCheckService struct {
	logger      *zerolog.Logger
	mainStorage contract.Storage
}

// NewHealtCheckService Constructor for HealtCheckService.
func NewHealtCheckService(
	logger *zerolog.Logger,
	mainStorage contract.Storage,
) contract.ServiceHealthCheck {
	return &healtCheckService{
		logger:      logger,
		mainStorage: mainStorage,
	}
}

// Ping.
func (s *healtCheckService) Ping(ctx context.Context) error {
	return s.mainStorage.Ping(ctx)
}
