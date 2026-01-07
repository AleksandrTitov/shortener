package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/database"
	"github.com/AleksandrTitov/shortener/internal/file"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	requestJSON struct {
		URL string `json:"url"`
	}

	responseJSON struct {
		Result string `json:"result"`
	}
)

func GetSorterURL(repo repository.Repository, conf *config.Config, gen id.GeneratorID) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(rw, "Разрешен только \"Content-Type: text/plain\"", http.StatusBadRequest)
			return
		}

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
		urlID, err := getURLID(urlOrigin, repo, gen)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)

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

		urlID, err := getURLID(urlOrigin.URL, repo, gen)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

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
		urlOrigin, ok := repo.Get(urlID)
		if !ok {
			http.Error(rw, fmt.Sprintf("ID \"%s\" не найден", urlID), http.StatusBadRequest)
			return
		}
		rw.Header().Add("Location", urlOrigin)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func Ping(conf *config.Config) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := database.Ping(conf)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		rw.WriteHeader(http.StatusOK)
	}
}
