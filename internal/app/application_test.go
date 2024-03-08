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
		request      string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, request: "/", body: "https://practicum.yandex.ru", expectedCode: http.StatusCreated, expectedBody: ""},
		{method: http.MethodPost, request: "/", body: "https://practicum.yandex.ru", expectedCode: http.StatusCreated, expectedBody: ""},
		{method: http.MethodPost, request: "/", body: "", expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			s.app.Route(w, r)

			s.Require().Equal(test.expectedCode, w.Code)
			if test.expectedBody != "" {
				s.Require().JSONEq(test.expectedBody, w.Body.String())
			}
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
		{method: http.MethodGet, request: "/test", expectedCode: http.StatusTemporaryRedirect, expectedBody: ""},
		{method: http.MethodGet, request: "/", expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	s.st.Set("test", "test")
	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, test.request, nil)
			w := httptest.NewRecorder()
			s.app.Route(w, r)
			s.Require().Equal(test.expectedCode, w.Code)
		})
	}
}

func (s *FunctionalTestSuite) TestUndefinedURL() {
	s.Run("undefined method PUT", func() {
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		w := httptest.NewRecorder()
		s.app.Route(w, r)
		s.Require().Equal(http.StatusMethodNotAllowed, w.Code)
	})
}
