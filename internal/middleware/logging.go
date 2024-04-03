package middleware

import (
	"net/http"
	"time"
)

// WithLogging middleware для логирования.
func (mw *middleware) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData: &responseData{
				status: 0,
				size:   0,
			},
		}
		l := mw.logger.Info().Str("URI", r.RequestURI)
		next.ServeHTTP(&lw, r)
		l.Dur("duration", time.Since(start))
		l.Int("status", lw.responseData.status)
		l.Int("size", lw.responseData.size).Send()
	})
}
