package tests

import (
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shrinkr/internal/app"
	"github.com/vagafonov/shrinkr/pkg/storage"
	"testing"
)

type FunctionalTestSuite struct {
	suite.Suite
	app *app.Application
	st  storage.Storage
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	s.st = storage.NewMemoryStorage()
	s.app = app.NewApplication(app.NewContainer(s.st))
}
