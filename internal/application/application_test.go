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

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/internal/service"
	"github.com/vagafonov/shortener/internal/storage"
	"github.com/vagafonov/shortener/pkg/entity"
)

const fileStoragePath = "short-url-db-test.json"

type FunctionalTestSuite struct {
	suite.Suite
	cnt        *container.Container
	app        *Application
	serviceURL *service.URLServiceMock
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	fss, err := storage.NewFileSystemStorage(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig("test", "http://test:8080", fileStoragePath, "test")
	lr := logger.CreateLogger(cfg.LogLevel)
	s.cnt = container.NewContainer(
		cfg,
		nil,
		fss,
		nil,
		lr,
		nil,
	)
	serv, err := service.ServiceURLFactory(s.cnt, "mock")
	if err != nil {
		log.Fatal(err)
	}
	s.serviceURL, _ = serv.(*service.URLServiceMock)
	s.cnt.SetServiceURL(s.serviceURL)
	s.app = NewApplication(
		s.cnt,
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
		init   func(s *FunctionalTestSuite)
	}{
		{
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru",
			code:   http.StatusCreated,
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetMakeShortURLResult(&entity.URL{
					UUID:     uuid.UUID{},
					Short:    "********",
					Original: "2",
				}, nil)
			},
		},
		{
			method: http.MethodPost,
			body:   "",
			code:   http.StatusBadRequest,
			init:   func(s *FunctionalTestSuite) {},
		},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			test.init(s)
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.createShortURL(w, r)
			s.Require().Equal(test.code, w.Code)
			if test.code == http.StatusCreated {
				u, err := url.Parse(w.Body.String())
				s.Require().NoError(err)
				s.Require().Len(strings.Trim(u.Path, "/"), s.cnt.GetConfig().ShortURLLength)
			}
		})
	}
}

func (s *FunctionalTestSuite) TestApiShorten() {
	tests := []struct {
		method string
		body   string
		code   int
		init   func(s *FunctionalTestSuite)
	}{
		{
			method: http.MethodPost,
			body:   `{"url":"https://practicum.yandex.ru"}`,
			code:   http.StatusCreated,
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetMakeShortURLResult(&entity.URL{
					UUID:     uuid.UUID{},
					Short:    "********",
					Original: "2",
				}, nil)
			},
		},
		{
			method: http.MethodPost,
			body:   "{,}",
			code:   http.StatusBadRequest,
			init: func(s *FunctionalTestSuite) {
			},
		},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			test.init(s)
			r := httptest.NewRequest(test.method, "/api/shorten", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			s.app.shorten(w, r)
			s.Require().Equal(test.code, w.Code)

			if w.Body.Bytes() != nil {
				decoder := json.NewDecoder(w.Body)
				var resp response.ShortenResponse
				err := decoder.Decode(&resp)
				s.Require().NoError(err)

				if test.code == http.StatusCreated {
					u, err := url.Parse(resp.Result)
					s.Require().NoError(err)
					s.Require().Len(strings.Trim(u.Path, "/"), s.cnt.GetConfig().ShortURLLength)
					s.Require().Equal("application/json", w.Header().Get("content-type"))
					s.Require().Empty(w.Header().Get("Content-Encoding"))
				}
			}
		})
	}
}

func (s *FunctionalTestSuite) TestGetURL() {
	ctx := context.Background()
	tests := []struct {
		method   string
		URL      string
		code     int
		location string
		init     func(s *FunctionalTestSuite)
	}{
		{
			method:   http.MethodGet,
			URL:      "/test",
			code:     http.StatusTemporaryRedirect,
			location: "https://practicum.yandex.ru",
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetGetShortURLResult(&entity.URL{
					UUID:     uuid.UUID{},
					Short:    "test",
					Original: "https://practicum.yandex.ru",
				}, nil)
			},
		},
		{
			method:   http.MethodGet,
			URL:      "/",
			code:     http.StatusMethodNotAllowed,
			location: "",
			init:     func(s *FunctionalTestSuite) {},
		},
		{
			method:   http.MethodGet,
			URL:      "/undefined-short-url",
			code:     http.StatusNotFound,
			location: "",
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetGetShortURLResult(nil, nil)
			},
		},
	}
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()
	for _, test := range tests {
		s.Run(test.method, func() {
			test.init(s)
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
}

func (s *FunctionalTestSuite) TestCompress() {
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("send encoded request", func() {
		s.serviceURL.SetMakeShortURLResult(&entity.URL{
			UUID:     uuid.UUID{},
			Short:    "********",
			Original: "ya.ru",
		}, nil)
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
	})
}

func (s *FunctionalTestSuite) TestShortenBatch() { //nolint:funlen
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("shorten batch", func() {
		s.serviceURL.SetMakeShortURLBatchResult([]response.ShortenBatchResponse{
			{
				CorrelationID: "1",
				ShortURL:      "a",
			},
			{
				CorrelationID: "2",
				ShortURL:      "b",
			},
		}, nil)
		requestBody := `[
			{
				"correlation_id": "1",
				"original_url": "aaa"
			},
			{
				"correlation_id": "2",
				"original_url": "bbb"
			}
		]`
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		s.Require().NoError(err)
		err = zb.Close()
		s.Require().NoError(err)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten/batch", buf)

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
		s.Require().JSONEq(`[{"correlation_id":"1","short_url":"a"},{"correlation_id":"2","short_url":"b"}]`, string(b))
	})

	s.Run("shorten batch with correlation_id empty", func() {
		requestBody := `[
			{
				"correlation_id": "",
				"original_url": "aaa"
			}
		]`
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		s.Require().NoError(err)
		err = zb.Close()
		s.Require().NoError(err)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten/batch", buf)

		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.Run("shorten batch with original_url empty", func() {
		requestBody := `[
			{
				"correlation_id": "1",
				"original_url": ""
			}
		]`
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		s.Require().NoError(err)
		err = zb.Close()
		s.Require().NoError(err)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten/batch", buf)

		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.Run("shorten empty batch", func() {
		requestBody := `[]`
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		s.Require().NoError(err)
		err = zb.Close()
		s.Require().NoError(err)
		r := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten/batch", buf)

		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})
}
