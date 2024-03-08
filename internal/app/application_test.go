package app

import (
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shrinkr/pkg/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FunctionalTestSuite struct {
	suite.Suite
	app *Application
	st  storage.Storage
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	s.st = storage.NewMemoryStorage()
	s.app = NewApplication(NewContainer(s.st))
}

func (s *FunctionalTestSuite) TestCreateURL() {
	tests := []struct {
		method       string
		URL          string
		body         string
		expectedCode int
	}{
		{method: http.MethodPost, URL: "/", body: "https://practicum.yandex.ru", expectedCode: http.StatusCreated},
		{method: http.MethodPost, URL: "/", body: "https://practicum.yandex.ru", expectedCode: http.StatusCreated},
		{method: http.MethodPost, URL: "/", body: "", expectedCode: http.StatusBadRequest},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()

	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, test.URL, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.createShortURL(w, r)
			s.Require().Equal(test.expectedCode, w.Code)
		})
	}
}

func (s *FunctionalTestSuite) TestGetURL() {
	tests := []struct {
		method       string
		request      string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, request: "/test", expectedCode: http.StatusTemporaryRedirect},
		{method: http.MethodGet, request: "/", expectedCode: http.StatusBadRequest},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()
	s.st.Set("test", "test")
	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, test.request, nil)
			w := httptest.NewRecorder()
			s.app.getShortURL(w, r)
			s.Require().Equal(test.expectedCode, w.Code)
		})
	}
}
