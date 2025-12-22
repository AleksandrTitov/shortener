package middleware

import (
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	mwResponseWriter struct {
		responseData *responseData
		http.ResponseWriter
	}
)

func (rw *mwResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.responseData.size = size

	return size, err
}

func (rw *mwResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.responseData.status = statusCode
}

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		now := time.Now()

		data := &responseData{
			size:   0,
			status: 0,
		}

		mrw := mwResponseWriter{
			responseData:   data,
			ResponseWriter: rw,
		}

		h.ServeHTTP(&mrw, r)

		logger.Log.WithFields(logrus.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": time.Since(now).String(),
			"size":     humanize.Bytes(uint64(mrw.responseData.size)),
			"status":   mrw.responseData.status,
		}).Info("http_request")
	})
}
