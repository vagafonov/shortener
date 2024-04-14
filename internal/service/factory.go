package service

import (
	"errors"

	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/contract"
)

var ErrUndefinedServiceType = errors.New("undefined service type")

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
