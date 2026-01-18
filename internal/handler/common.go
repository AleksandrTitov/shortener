package handler

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"net/url"
)

func getURLID(url string, repo repository.Repository, gen id.GeneratorID) (string, error) {
	maxAttempts := 15
	var urlID string
	var err error

	for i := 1; i <= maxAttempts; i++ {
		urlID, err = gen.GetID()
		if errors.Is(err, id.ErrGetID) {
			logger.Log.Errorf("Не удалось сгенерировать ID: %v", err.Error())
			break
		}
		err = repo.Set(urlID, url)
		if errors.Is(err, repository.ErrorAlreadyExist) {
			logger.Log.Warnf("Не удалось записать id \"%s\" попытка %d(%d), %v", urlID, i, maxAttempts, err.Error())
			continue
		} else if err != nil {
			logger.Log.Errorf("Не ожиданная ошибка при записи id \"%s\": %v", urlID, err)
			break
		}
		break
	}

	return urlID, err
}

func originalURLUniqueViolation(err error, repo repository.Repository, conf *config.Config, rw http.ResponseWriter, urlOrigin string) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
		return false
	}

	rw.WriteHeader(http.StatusConflict)

	urlID, err := repo.GetByURL(urlOrigin)
	if err != nil {
		logger.Log.Errorf("Ошибка получения Original URL: %v", err.Error())
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}

	urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
	if err != nil {
		logger.Log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}

	_, err = rw.Write([]byte(urlShort))
	if err != nil {
		logger.Log.Errorf("Не удалось записать данные \"%v\"", err.Error())
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}

	return true
}
