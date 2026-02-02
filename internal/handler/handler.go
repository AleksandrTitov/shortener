package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/file"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/middleware"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	requestJSON struct {
		URL string `json:"url"`
	}

	responseJSON struct {
		Result string `json:"result"`
	}

	requestBatchJSON struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	responseBatchJSON struct {
		CorrelationID string `json:"correlation_id"`
		SortURL       string `json:"short_url"`
	}
)

func GetSorterURL(repo repository.Repository, conf *config.Config, gen id.GeneratorID) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(rw, "Разрешен только \"Content-Type: text/plain\"", http.StatusBadRequest)
			return
		}

		userID := r.Header.Get(middleware.UserIDHeader)

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		urlOrigin := strings.TrimSpace(string(body))
		_, err = url.ParseRequestURI(urlOrigin)
		if err != nil {
			http.Error(rw, "В данных запроса ожидаться валидный URL", http.StatusBadRequest)
			return
		}
		urlID, err := getURLID(urlOrigin, userID, repo, gen)

		switch {
		case err == nil:
			rw.WriteHeader(http.StatusCreated)
		case errors.Is(err, repository.ErrorOriginNotUnique):
			rw.WriteHeader(http.StatusConflict)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
		if err != nil {
			logger.Log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if conf.FileName != "" {
			items := file.NewShorterItems()
			err = items.SaveShorterItems(conf.FileName, repo)
			if err != nil {
				logger.Log.Errorf("Не удалось записать данные в файл %s: %v", conf.FileName, err)
			}
		}

		_, err = rw.Write([]byte(urlShort))
		if err != nil {
			logger.Log.Errorf("Не удалось записать данные \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func GetSorterURLJson(repo repository.Repository, conf *config.Config, gen id.GeneratorID) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(rw, "Разрешен только \"Content-Type: application/json\"", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var urlOrigin requestJSON
		err = json.Unmarshal(body, &urlOrigin)
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = url.ParseRequestURI(urlOrigin.URL)
		if err != nil {
			http.Error(rw, "В данных запроса ожидаться валидный URL", http.StatusBadRequest)
			return
		}

		userID := r.Header.Get(middleware.UserIDHeader)

		urlID, err := getURLID(urlOrigin.URL, userID, repo, gen)

		rw.Header().Set("Content-Type", "application/json")

		switch {
		case err == nil:
			rw.WriteHeader(http.StatusCreated)
		case errors.Is(err, repository.ErrorOriginNotUnique):
			rw.WriteHeader(http.StatusConflict)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
		if err != nil {
			logger.Log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if conf.FileName != "" {
			items := file.NewShorterItems()
			err = items.SaveShorterItems(conf.FileName, repo)
			if err != nil {
				logger.Log.Errorf("Не удалось записать данные в файл %s: %v", conf.FileName, err)
			}
		}

		urlShortJSON := responseJSON{
			Result: urlShort,
		}
		resp, _ := json.Marshal(urlShortJSON)
		_, err = rw.Write(resp)
		if err != nil {
			logger.Log.Errorf("Не удалось записать данные \"%v\"", err.Error())
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
		urlOrigin, ok, gone := repo.Get(urlID)
		if !ok {
			http.Error(rw, fmt.Sprintf("ID \"%s\" не найден", urlID), http.StatusBadRequest)
			return
		}
		if gone {
			http.Error(rw, fmt.Sprintf("ID \"%s\" удален", urlID), http.StatusGone)
			return
		}
		rw.Header().Add("Location", urlOrigin)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func Ping(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if db == nil {
			logger.Log.Errorf("Подключение к БД отсутствует")
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		timeout := time.Second * 3
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Log.Errorf("Проверка подключения к БД завершилась по таймауту(%s): %v", timeout, err)
			} else {
				logger.Log.Errorf("Проверка подключения к БД: %v", err)
			}

			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
		logger.Log.Debugf("Проверка подключения к БД выполнена успешно")
	}
}

func GetShorterURLJsonBatch(repo repository.Repository, conf *config.Config, gen id.GeneratorID) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var requestBatch []requestBatchJSON
		var responseBatch []responseBatchJSON

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(rw, "Разрешен только \"Content-Type: application/json\"", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &requestBatch)
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if len(requestBatch) == 0 {
			logger.Log.Errorf("Пустой батч")
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		urls := map[string]string{}

		for _, i := range requestBatch {
			urlID, err := gen.GetID()
			if i.OriginalURL == "" {
				logger.Log.Errorf("Original URL пустой")
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if i.CorrelationID == "" {
				logger.Log.Errorf("Correlation ID пустой")
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if errors.Is(err, id.ErrGetID) {
				logger.Log.Errorf("Не удалось сгенерировать ID: %v", err.Error())
				break
			}
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
			if err != nil {
				logger.Log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			responseBatch = append(responseBatch, responseBatchJSON{
				CorrelationID: i.CorrelationID,
				SortURL:       urlShort,
			})
			urls[urlID] = i.OriginalURL
		}

		userID := r.Header.Get(middleware.UserIDHeader)

		err = repo.SetBatch(urls, userID)
		if err != nil {
			logger.Log.Errorf("Не удалось сохранть батч \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

		resp, _ := json.Marshal(responseBatch)
		_, err = rw.Write(resp)
		if err != nil {
			logger.Log.Errorf("Не удалось записать данные \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func GetUsersURLJson(repo repository.Repository, conf *config.Config) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get(middleware.UserIDHeader)
		urls, err := repo.GetByUserID(userID)
		if err != nil {
			logger.Log.Errorf("Ошибка получения Users URLs пользователя \"%s\": %v", userID, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")

		switch {
		case len(urls) == 0:
			rw.WriteHeader(http.StatusNoContent)
			return
		default:
			rw.WriteHeader(http.StatusOK)
		}
		for i := range urls {
			urls[i].URLID, err = url.JoinPath(conf.BaseHTTP, urls[i].URLID)
			if err != nil {
				logger.Log.Errorf("Не удалось создать короткий URL \"%v\"", err.Error())
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
		logger.Log.Debugf("Users URLs пользователя %s получены: %v", userID, urls)

		resp, err := json.Marshal(urls)
		if err != nil {
			logger.Log.Errorf("Ошибка десерелизации Users URLs пользователя \"%s\": %v", userID, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = rw.Write(resp)
		if err != nil {
			logger.Log.Errorf("Не удалось записать данные \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteURLs(repo repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(rw, "Разрешен только \"Content-Type: application/json\"", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			logger.Log.Errorf("Ошибка чтения запроса \"%v\"", err.Error())
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		userID := r.Header.Get(middleware.UserIDHeader)

		var urlsToDelete []string
		err = json.Unmarshal(body, &urlsToDelete)
		if err != nil {
			logger.Log.Errorf("Ошибка десерелизации Users URLs пользователя \"%s\": %v", userID, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Log.Debugf("Сприсок URL на удаление: %v, от пользователя: %s", urlsToDelete, userID)
		rw.WriteHeader(http.StatusAccepted)

		go func() {
			if len(urlsToDelete) == 0 {
				return
			}

			err = repo.DeleteIDs(userID, urlsToDelete)
			if err != nil {
				logger.Log.Errorf("Ошибка при асинхронном удалении URL для user %s: %v", userID, err)
			} else {
				logger.Log.Debugf("Асинхронно удалено %d URL для user %s", len(urlsToDelete), userID)
			}
		}()
	}
}
