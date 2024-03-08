package tests

import (
	"net/http"
	"net/http/httptest"
)

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
