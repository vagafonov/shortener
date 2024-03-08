package app

import (
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shrinkr/pkg/storage"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		method string
		URL    string
		body   string
		code   int
		result string
	}{
		{method: http.MethodPost, URL: "/", body: "https://practicum.yandex.ru", code: http.StatusCreated, result: "1"},
		{method: http.MethodPost, URL: "/", body: "https://practicum.yandex.ru", code: http.StatusCreated, result: "2"},
		{method: http.MethodPost, URL: "/", body: "", code: http.StatusBadRequest},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()

	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, test.URL, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.createShortURL(w, r)
			s.Require().Equal(test.code, w.Code)
			if test.result != "" {
				u, err := url.Parse(w.Body.String())
				s.Require().NoError(err)
				s.Require().Len(strings.Trim(u.Path, "/"), ShortURLLength)
			}
		})
	}
	s.Require().Len(s.st.GetAll(), 1, "exists doubles for same url")
	// TODO move to tearDown
	s.st.Truncate()
}

func (s *FunctionalTestSuite) TestGetURL() {
	tests := []struct {
		method   string
		URL      string
		code     int
		location string
	}{
		{method: http.MethodGet, URL: "/test", code: http.StatusTemporaryRedirect, location: "https://practicum.yandex.ru"},
		{method: http.MethodGet, URL: "/", code: http.StatusMethodNotAllowed, location: ""},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()
	// TODO use dummy page
	_, err := s.st.Add("test", "https://practicum.yandex.ru")
	s.Require().NoError(err)
	for _, test := range tests {
		s.Run(test.method, func() {
			r, err := http.NewRequest(test.method, ts.URL+test.URL, nil)
			s.Require().NoError(err)

			cli := ts.Client()
			cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
			resp, err := cli.Do(r)
			s.Require().NoError(err)
			defer resp.Body.Close()
			s.Require().Equal(test.code, resp.StatusCode)
			s.Require().Equal(test.location, resp.Header.Get("Location"))
		})
	}
	// TODO move to tearDown
	s.st.Truncate()
}
