package middleware

import "net/http"

type (
	// Cтруктура для хранения сведений об ответе.
	responseData struct {
		status int
		size   int
	}

	// Кастомная реализация http.ResponseWriter для получения данных ответа.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер

	return size, err
}

// WriteHeader.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
