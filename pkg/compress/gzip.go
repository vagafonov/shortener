package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

const statusToGzip = 500

type compressGzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressGzipWriter(w http.ResponseWriter) *compressGzipWriter {
	return &compressGzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressGzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressGzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressGzipWriter) WriteHeader(statusCode int) {
	if statusCode < statusToGzip {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressGzipWriter) Close() error {
	return c.zw.Close()
}

type compressGzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressGzipReader(r io.ReadCloser) (*compressGzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressGzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressGzipReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressGzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}
