//nolint:funlen
package application

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vagafonov/shortener/internal/config"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/cookie"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/internal/logger"
	"github.com/vagafonov/shortener/internal/response"
	"github.com/vagafonov/shortener/internal/service"
	"github.com/vagafonov/shortener/internal/storage"
	"github.com/vagafonov/shortener/pkg/encrypting"
	"github.com/vagafonov/shortener/pkg/entity"
)

const fileStoragePath = "short-url-db-test.json"

type FunctionalTestSuite struct {
	suite.Suite
	cnt                *container.Container
	app                *Application
	serviceURL         *service.URLServiceMock
	serviceHealthCheck *service.HealthCheckServiceMock
}

func TestFunctionalTestSuite(t *testing.T) {
	suite.Run(t, new(FunctionalTestSuite))
}

func (s *FunctionalTestSuite) SetupSuite() {
	fss, err := storage.NewFileSystemStorage(fileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig(
		"test",
		"http://test:8080",
		fileStoragePath,
		"test",
		false,
		[]byte("0123456789abcdef"),
		10,
		3,
		config.ModeTest,
		"192.168.0.1/24",
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

	s.app = NewApplication(
		s.cnt,
	)
}

func (s *FunctionalTestSuite) TearDownSuite() {
	os.Remove(fileStoragePath)
}

//nolint:dupl
func (s *FunctionalTestSuite) TestCreateURL() {
	srv := httptest.NewServer(s.app.Routes())
	tests := []struct {
		method   string
		body     string
		code     int
		init     func(s *FunctionalTestSuite)
		expected string
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
			expected: "http://test:8080/********",
		},
		{
			method:   http.MethodPost,
			body:     "",
			code:     http.StatusBadRequest,
			init:     func(s *FunctionalTestSuite) {},
			expected: "",
		},
		{
			method: http.MethodPost,
			body:   "http://test.local",
			code:   http.StatusConflict,
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetMakeShortURLResult(&entity.URL{
					UUID:     uuid.UUID{},
					Short:    "********",
					Original: "http://test.local",
				}, customerror.ErrURLAlreadyExists)
			},
			expected: "http://test:8080/********",
		},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			test.init(s)
			r := httptest.NewRequest(test.method, srv.URL+"/", strings.NewReader(test.body))
			ck := cookie.CreateCookieWithUserID(s.cnt.GetLogger(), s.cnt.GetConfig().CryptoKey)
			r.RequestURI = ""
			r.AddCookie(ck)
			resp, err := http.DefaultClient.Do(r)
			s.Require().NoError(err)
			defer resp.Body.Close()
			s.Require().Equal(test.code, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)
			s.Require().Equal(test.expected, string(b))
		})
	}
}

//nolint:funlen, dupl
func (s *FunctionalTestSuite) TestApiShorten() {
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()
	tests := []struct {
		method   string
		body     string
		code     int
		init     func(s *FunctionalTestSuite)
		expected string
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
			expected: `{"result":"http://test:8080/********"}`,
		},
		{
			method: http.MethodPost,
			body:   "{,}",
			code:   http.StatusBadRequest,
			init: func(s *FunctionalTestSuite) {
			},
			expected: `{"result":"http://test:8080/********"}`,
		},
		{
			method: http.MethodPost,
			body:   `{"url":"https://practicum.yandex.ru"}`,
			code:   http.StatusConflict,
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetMakeShortURLResult(&entity.URL{
					UUID:     uuid.UUID{},
					Short:    "********",
					Original: "http://test.local",
				}, customerror.ErrURLAlreadyExists)
			},
			expected: `{"result":"http://test:8080/********"}`,
		},
	}

	for _, test := range tests {
		s.Run(test.method, func() {
			test.init(s)
			r := httptest.NewRequest(test.method, srv.URL+"/api/shorten", strings.NewReader(test.body))
			ck := cookie.CreateCookieWithUserID(s.cnt.GetLogger(), s.cnt.GetConfig().CryptoKey)
			r.RequestURI = ""
			r.AddCookie(ck)
			resp, err := http.DefaultClient.Do(r)
			s.Require().NoError(err)
			defer resp.Body.Close()
			s.Require().Equal(test.code, resp.StatusCode)
			b, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)

			if string(b) != "" {
				s.Require().Equal(`application/json`, resp.Header.Get("Content-Type"))
				s.Require().JSONEq(test.expected, string(b))
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
		{
			method:   http.MethodGet,
			URL:      "/deleted-short-url",
			code:     http.StatusGone,
			location: "",
			init: func(s *FunctionalTestSuite) {
				s.serviceURL.SetGetShortURLResult(nil, customerror.ErrURLDeleted)
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

		ck := cookie.CreateCookieWithUserID(s.cnt.GetLogger(), s.cnt.GetConfig().CryptoKey)
		r.AddCookie(ck)
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

func (s *FunctionalTestSuite) TestCheckUserIDInCookie() {
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("exist and valid", func() {
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/ping", strings.NewReader(""))
		r.RequestURI = ""
		userID, err := uuid.NewUUID()
		s.Require().NoError(err)
		uuidString := userID.String()
		encrypted, err := encrypting.Encrypt(uuidString, s.cnt.GetConfig().CryptoKey)
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
	})

	s.Run("exist and invalid", func() {
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/ping", strings.NewReader(""))
		r.RequestURI = ""
		userID, err := uuid.NewUUID()
		s.Require().NoError(err)
		uuidString := userID.String()
		encrypted, err := encrypting.Encrypt(uuidString, []byte("****************"))
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
	})

	s.Run("does not exist", func() {
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/ping", strings.NewReader(""))
		r.RequestURI = ""
		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
	})
}

func (s *FunctionalTestSuite) TestApiUserURLs() { //nolint:funlen
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("get user URLs", func() {
		userID := uuid.New()
		s.serviceURL.SetGetUserURLsResult([]*entity.URL{
			{
				UUID:     uuid.UUID{},
				Short:    "********",
				Original: "2",
				UserID:   userID,
			},
		}, nil)
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/api/user/urls", strings.NewReader(""))
		r.RequestURI = ""
		encrypted, err := encrypting.Encrypt(userID.String(), s.cnt.GetConfig().CryptoKey)
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		s.Require().Equal(`application/json`, resp.Header.Get("Content-Type"))
		s.Require().JSONEq(`[{"short_url":"********","original_url":"2"}]`, string(b))
	})

	// Если кука не содержит ID пользователя, хендлер должен возвращать HTTP-статус 401 Unauthorized.
	s.Run("request with empty userID in cookie", func() {
		userID := uuid.New()
		s.serviceURL.SetGetUserURLsResult([]*entity.URL{
			{
				UUID:     uuid.UUID{},
				Short:    "********",
				Original: "2",
				UserID:   userID,
			},
		}, nil)
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/api/user/urls", strings.NewReader(""))
		r.RequestURI = ""
		cookie := &http.Cookie{Name: "userID", Value: ""}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	// При отсутствии сокращённых пользователем URL хендлер должен отдавать HTTP-статус 204 No Content.
	s.Run("request with empty userID in cookie", func() {
		userID := uuid.New()
		s.serviceURL.SetGetUserURLsResult([]*entity.URL{}, nil)
		r := httptest.NewRequest(http.MethodGet, srv.URL+"/api/user/urls", strings.NewReader(""))
		r.RequestURI = ""
		encrypted, err := encrypting.Encrypt(userID.String(), s.cnt.GetConfig().CryptoKey)
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
	})
}

func (s *FunctionalTestSuite) TestDeleteUserURLs() {
	srv := httptest.NewServer(s.app.Routes())
	defer srv.Close()

	s.Run("delete user URLs", func() {
		userID := uuid.New()
		s.serviceURL.SetDeleteUserURLsResult(nil)
		r := httptest.NewRequest(
			http.MethodDelete,
			srv.URL+"/api/user/urls",
			strings.NewReader(`["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`),
		)
		r.RequestURI = ""
		encrypted, err := encrypting.Encrypt(userID.String(), s.cnt.GetConfig().CryptoKey)
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusAccepted, resp.StatusCode)
	})

	s.Run("delete empty user URLs", func() {
		userID := uuid.New()
		s.serviceURL.SetDeleteUserURLsResult(nil)
		r := httptest.NewRequest(
			http.MethodDelete,
			srv.URL+"/api/user/urls",
			strings.NewReader(`[]`),
		)
		r.RequestURI = ""
		encrypted, err := encrypting.Encrypt(userID.String(), s.cnt.GetConfig().CryptoKey)
		s.Require().NoError(err)
		cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encrypted)}
		r.AddCookie(cookie)

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *FunctionalTestSuite) TestInternalStats() {
	ctx := context.Background()
	ts := httptest.NewServer(s.app.Routes())
	defer ts.Close()
	s.Run("successfully", func() {
		s.serviceURL.SetGetStatResult(&entity.Stat{
			Urls:  1,
			Users: 1,
		}, nil)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/internal/stats", nil)
		s.Require().NoError(err)
		r.Header.Set("X-Real-IP", "192.168.0.1")

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		s.Require().Equal(`application/json`, resp.Header.Get("Content-Type"))
		s.Require().JSONEq(`{"urls":1,"users":1}`, string(b))
	})

	s.Run("not in subnet", func() {
		r, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/internal/stats", nil)
		s.Require().NoError(err)
		r.Header.Set("X-Real-IP", "0.0.0.0")

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusForbidden, resp.StatusCode)
	})

	s.Run("empty config value for trusted_subnet", func() {
		s.cnt.GetConfig().TrustedSubnet = ""

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/internal/stats", nil)
		s.Require().NoError(err)
		r.Header.Set("X-Real-IP", "0.0.0.0")

		resp, err := http.DefaultClient.Do(r)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusForbidden, resp.StatusCode)
	})
}
