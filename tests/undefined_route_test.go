package tests

import (
	"net/http"
	"net/http/httptest"
)

func (s *FunctionalTestSuite) TestUndefinedURL() {
	s.Run("undefined method PUT", func() {
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		w := httptest.NewRecorder()
		s.app.Route(w, r)
		s.Require().Equal(http.StatusMethodNotAllowed, w.Code)
	})
}
