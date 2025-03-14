package middleware

import (
	"compress/gzip"
	"net/http"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Переопределяет метод Header интерфейса ResponseWriter
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Переопределяет метод Write интерфейса ResponseWriter
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Переопределяет метод WriteHeader интерфейса ResponseWriter
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}
