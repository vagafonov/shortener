package application

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/application/storage"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/pkg/entity"
	hash "github.com/vagafonov/shortener/pkg/hasher"
)

const fileStoragePath = "short-url-db-test.json"

type FunctionalTestSuite struct {
	suite.Suite
	app    *Application
	st     storage.Storage
	fss    storage.Storage
	cfg    *config.Config
	hasher hash.Hasher
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	s.st = storage.NewMemoryStorage()
	var err error
	s.fss, err = storage.NewFileSystemStorage(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	s.cfg = config.NewConfig("test", "http://test:8080", fileStoragePath)
	s.hasher = hash.NewMockHasher()
	s.app = NewApplication(
		NewContainer(
			s.cfg,
			s.st,
			s.fss,
			s.hasher,
		),
	)
}

func (s *FunctionalTestSuite) TearDownSuite() {
	os.Remove(fileStoragePath)
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
	URLs, err := s.st.GetAll()
	s.Require().NoError(err)
	s.Require().Len(URLs, 1, "exists doubles for same url")
	// TODO move to tearDown
	s.st.Truncate()
}

func (s *FunctionalTestSuite) TestApiShorten() {
	tests := []struct {
		method string
		body   string
		code   int
	}{
		{method: http.MethodPost, body: `{"url":"https://practicum.yandex.ru"}`, code: http.StatusCreated},
		{method: http.MethodPost, body: `{"url":"https://practicum.yandex.ru"}`, code: http.StatusCreated},
		{method: http.MethodPost, body: "{,}", code: http.StatusBadRequest},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			r := httptest.NewRequest(test.method, "/api/shorten", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.shorten(w, r)
			s.Require().Equal(test.code, w.Code)

			if w.Body.Bytes() != nil {
				decoder := json.NewDecoder(w.Body)
				var resp entity.ShortenResponse
				err := decoder.Decode(&resp)
				s.Require().NoError(err)

				if test.code == http.StatusCreated {
					u, err := url.Parse(resp.Result)
					s.Require().NoError(err)
					s.Require().Len(strings.Trim(u.Path, "/"), s.cfg.ShortURLLength)
					s.Require().Equal("application/json", w.Header().Get("content-type"))
					s.Require().Empty(w.Header().Get("Content-Encoding"))
				}
			}
		})
	}
	URLs, err := s.st.GetAll()
	s.Require().NoError(err)
	s.Require().Len(URLs, 1, "exists doubles for same url")
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

func (s *FunctionalTestSuite) TestCompress() {
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("send encoded request", func() {
		requestBody := `{"url": "ya.ru"}`
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		s.Require().NoError(err)
		err = zb.Close()
		s.Require().NoError(err)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten", buf)

		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()
		zr, err := gzip.NewReader(resp.Body)
		s.Require().NoError(err)

		b, err := io.ReadAll(zr)
		s.Require().NoError(err)

		s.Require().Equal(`gzip`, resp.Header.Get("Content-Encoding"))
		s.Require().Equal(`application/json`, resp.Header.Get("Content-Type"))
		s.Require().JSONEq(`{"result":"http://test:8080/********"}`, string(b))
		// TODO move to tearDown
		s.st.Truncate()
	})
}
