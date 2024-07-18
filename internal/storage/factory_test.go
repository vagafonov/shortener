package storage

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/application"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/service"
)

const fileStoragePath = "/dev/null"

type FactoryTestSuite struct {
	suite.Suite
	cnt                *container.Container
	app                *application.Application
	serviceURL         *service.URLServiceMock
	serviceHealthCheck *service.HealthCheckServiceMock
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}

func (s *FactoryTestSuite) SetupSuite() {
	fss, err := NewFileSystemStorage(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig(
		"test",
		"http://test:8080",
		fileStoragePath,
		"test",
		"",
		[]byte("0123456789abcdef"),
		10,
		3,
		config.ModeTest,
	)
	lr := logger.CreateLogger(cfg.LogLevel)
	s.cnt = container.NewContainer(
		cfg,
		nil,
		fss,
		nil,
		lr,
		nil,
	)
	servURL, err := service.ServiceURLFactory(s.cnt, "mock")
	if err != nil {
		log.Fatal(err)
	}
	s.serviceURL, _ = servURL.(*service.URLServiceMock)
	s.cnt.SetServiceURL(s.serviceURL)

	servHealthcheck, err := service.ServiceHealthCheckFactory(s.cnt, "mock")
	if err != nil {
		log.Fatal(err)
	}
	s.serviceHealthCheck, _ = servHealthcheck.(*service.HealthCheckServiceMock)
	s.cnt.SetServiceHealthCheck(s.serviceHealthCheck)

	s.app = application.NewApplication(
		s.cnt,
	)
}

func (s *FactoryTestSuite) TestDB() {
	_, err := StorageFactory(s.cnt, "db")
	s.Require().NoError(err)
}

func (s *FactoryTestSuite) TestFS() {
	_, err := StorageFactory(s.cnt, "fs")
	s.Require().NoError(err)
}

func (s *FactoryTestSuite) TestFSMock() {
	_, err := StorageFactory(s.cnt, "fs-mock")
	s.Require().NoError(err)
}

func (s *FactoryTestSuite) TestMemoryMock() {
	_, err := StorageFactory(s.cnt, "memory-mock")
	s.Require().NoError(err)
}
