package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

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
