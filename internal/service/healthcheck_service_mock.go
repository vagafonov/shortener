package service

import (
	"context"

	"github.com/vagafonov/shortener/internal/contract"
)

type HealthCheckServiceMock struct{}

func NewHealthCheckServiceMock() contract.ServiceHealthCheck {
	return &HealthCheckServiceMock{}
}

func (s *HealthCheckServiceMock) Ping(ctx context.Context) error {
	return nil
}
