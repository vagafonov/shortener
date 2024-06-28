package service

import (
	"errors"

	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
)

// ErrUndefinedServiceType error for undefined service type.
var ErrUndefinedServiceType = errors.New("undefined service type")

// ServiceURLFactory return concrete service URL.
func ServiceURLFactory(cnt *container.Container, t string) (contract.Service, error) {
	// TODO use enum
	switch t {
	case "real":
		return NewURLService(
			cnt.GetLogger(),
			cnt.GetMainStorage(),
			cnt.GetBackupStorage(),
			cnt.GetHasher(),
		), nil
	case "mock":
		return NewURLServiceMock(), nil
	default:
		return nil, ErrUndefinedServiceType
	}
}

// ServiceHealthCheckFactory return concrete service health.
func ServiceHealthCheckFactory(cnt *container.Container, t string) (contract.ServiceHealthCheck, error) {
	// TODO use enum
	switch t {
	case "real":
		return NewHealtCheckService(
			cnt.GetLogger(),
			cnt.GetMainStorage(),
		), nil
	case "mock":
		return NewHealthCheckServiceMock(), nil
	default:
		return nil, ErrUndefinedServiceType
	}
}
