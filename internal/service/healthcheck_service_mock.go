package service

import (
	"context"

	"github.com/vagafonov/shortener/internal/contract"
)

// mock.
type HealthCheckServiceMock struct{}

// NewHealthCheckServiceMock Constructor for HealthCheckServiceMock.
func NewHealthCheckServiceMock() contract.ServiceHealthCheck {
	return &HealthCheckServiceMock{}
}

// mock.
func (s *HealthCheckServiceMock) Ping(ctx context.Context) error {
	return nil
}
