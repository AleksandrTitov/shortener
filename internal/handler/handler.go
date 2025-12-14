package handler

import (
	"errors"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
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

func MiddlewareLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log := logger.NewLogger()

		data := &responseData{
			size:   0,
			status: 0,
		}

		mrw := mwResponseWriter{
			responseData:   data,
			ResponseWriter: rw,
		}

		h.ServeHTTP(&mrw, r)

		log.WithFields(logrus.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": fmt.Sprintf("%s", time.Since(now)),
			"size":     humanize.Bytes(uint64(mrw.responseData.size)),
			"status":   mrw.responseData.status,
		}).Info("Ok")
	}
}

func GetSorterURL(repo repository.Repository, conf *config.Config, gen id.GeneratorID) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(rw, "Разрешен только \"Content-Type: text/plain\"", http.StatusBadRequest)
			return
		}

		log := logger.NewLogger()
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		urlOrigin := string(body)

		_, err = url.ParseRequestURI(urlOrigin)
		if err != nil {
			http.Error(rw, "В данных запроса ожидаться валидный URL", http.StatusBadRequest)
			return
		}

		maxAttempts := 15
		var urlID string

		for i := 1; i <= maxAttempts; i++ {
			urlID, err = gen.GetID()
			if errors.Is(err, id.ErrGetID) {
				log.Errorf("ERROR: Не удалось сгенерировать ID: %v", err.Error())
				break
			}
			err = repo.Set(urlID, urlOrigin)
			if errors.Is(err, repository.ErrorAlreadyExist) {
				log.Warnf("Не удалось записать id \"%s\" попытка %d(%d), %v", urlID, i, maxAttempts, err.Error())
				continue
			} else if err != nil {
				log.Errorf("Не ожиданная ошибка при записи id \"%s\": %v", urlID, err)
				break
			}
			break
		}
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)

		urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
		if err != nil {
			log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, err = rw.Write([]byte(urlShort))
		if err != nil {
			log.Errorf("Не удалось записать данные \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func GetOriginalURL(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		urlID := r.PathValue("urlID")

		if len(urlID) < id.LenID || len(urlID) > id.LenID {
			http.Error(rw, fmt.Sprintf("Длина ID должна быть равна %d символам", id.LenID), http.StatusBadRequest)
			return
		}
		urlOrigin, ok := repo.Get(urlID)
		if !ok {
			http.Error(rw, fmt.Sprintf("ID \"%s\" не найден", urlID), http.StatusBadRequest)
			return
		}
		rw.Header().Add("Location", urlOrigin)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
