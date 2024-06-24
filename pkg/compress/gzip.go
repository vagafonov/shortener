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

// Constructor for CompressGzipWriter.
func NewCompressGzipWriter(w http.ResponseWriter) *compressGzipWriter {
	return &compressGzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header return header.
func (c *compressGzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write write to gzip writer.
func (c *compressGzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader set header Content-Encoding and write status.
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

// Constructor for CompressGzipReader.
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

// Read read from reader.
func (c compressGzipReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

// Close ReadCloser and Reader.
func (c *compressGzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}
