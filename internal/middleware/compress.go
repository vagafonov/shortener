package middleware

import (
	"net/http"
	"strings"

	"github.com/vagafonov/shortener/pkg/compress"
)

// middleware для сжатия.
func (mw *middleware) WithCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w
		if strings.Contains(r.RequestURI, "/debug/") {
			next.ServeHTTP(originalWriter, r)

			return
		}

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		gzipContents := map[string]struct{}{
			"application/json": {},
			"text/html":        {},
			"text/plain":       {},
		}

		_, foundGzipFormat := gzipContents[r.Header.Get("Content-Type")]

		if supportsGzip || foundGzipFormat {
			compressWriter := compress.NewCompressGzipWriter(w)
			originalWriter = compressWriter
			defer compressWriter.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			compressReader, err := compress.NewCompressGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			// меняем тело запроса на новое
			r.Body = compressReader

			defer compressReader.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(originalWriter, r)
	})
}
