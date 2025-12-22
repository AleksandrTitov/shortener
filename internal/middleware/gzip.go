package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipWrite(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(rw, r)
			return
		}

		gz, err := gzip.NewWriterLevel(rw, gzip.BestSpeed)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		rw.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{
			rw,
			gz,
		}, r)
	})
}

func GzipRead(h http.Handler) http.Handler {
	allowContentTypes := map[string]bool{
		"application/json":                true,
		"application/json; charset=utf-8": true,
		"text/plain":                      true,
	}
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.ServeHTTP(rw, r)
			return
		}
		contentType := r.Header.Get("Content-Type")
		_, ok := allowContentTypes[contentType]
		if !ok {
			h.ServeHTTP(rw, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(rw, "некорректные данные в формате gzip", http.StatusBadRequest)
			return
		}
		defer gz.Close()
		r.Body = gz
		h.ServeHTTP(rw, r)
	})
}
