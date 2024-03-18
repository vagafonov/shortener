package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	storage2 "github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
)

type FunctionalTestSuite struct {
	suite.Suite
	app *Application
	st  storage2.Storage
	cfg *config.Config
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	s.st = storage2.NewMemoryStorage()
	s.cfg = config.NewConfig("test", "http://test:8080")
	s.app = NewApplication(NewContainer(s.cfg, s.st))
}

func (s *FunctionalTestSuite) TestCreateURL() {
	tests := []struct {
		method string
		body   string
		code   int
	}{
		{method: http.MethodPost, body: "https://practicum.yandex.ru", code: http.StatusCreated},
		{method: http.MethodPost, body: "https://practicum.yandex.ru", code: http.StatusCreated},
		{method: http.MethodPost, body: "", code: http.StatusBadRequest},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()

	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.createShortURL(w, r)
			s.Require().Equal(test.code, w.Code)
			if test.code == http.StatusCreated {
				u, err := url.Parse(w.Body.String())
				s.Require().NoError(err)
				s.Require().Len(strings.Trim(u.Path, "/"), s.cfg.ShortURLLength)
			}
		})
	}
	s.Require().Len(s.st.GetAll(), 1, "exists doubles for same url")
	// TODO move to tearDown
	s.st.Truncate()
}

func (s *FunctionalTestSuite) TestGetURL() {
	ctx := context.Background()
	tests := []struct {
		method   string
		URL      string
		code     int
		location string
	}{
		{method: http.MethodGet, URL: "/test", code: http.StatusTemporaryRedirect, location: "https://practicum.yandex.ru"},
		{method: http.MethodGet, URL: "/", code: http.StatusMethodNotAllowed, location: ""},
		{method: http.MethodGet, URL: "/undefined-short-url", code: http.StatusNotFound, location: ""},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()
	// TODO use dummy page
	_, err := s.st.Add("test", "https://practicum.yandex.ru")
	s.Require().NoError(err)
	for _, test := range tests {
		s.Run(test.method, func() {
			r, err := http.NewRequestWithContext(ctx, test.method, ts.URL+test.URL, nil)
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
