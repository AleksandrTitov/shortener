package handler

import (
	"context"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/go-chi/chi/v5"
	"io"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name   string
		target string
		id     string
		url    string
		body   string
		code   int
	}{
		{
			name:   "Получение оригинального URL",
			target: "dmiWnD",
			id:     "dmiWnD",
			url:    "http://test.aa",
			code:   http.StatusTemporaryRedirect,
		},
		{
			name:   "ID не найден",
			target: "dmiWnD",
			body:   "ID \"dmiWnD\" не найден\n",
			code:   http.StatusBadRequest,
		},
		{
			name:   "Длина запрашиваемого ID больше длины формата",
			target: "dmiWnDDD",
			body:   fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code:   http.StatusBadRequest,
		},
		{
			name:   "Длина запрашиваемого ID меньше длины формата",
			target: "dmi",
			body:   fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code:   http.StatusBadRequest,
		},
		{
			name: "Пустой ID",
			body: fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", test.target), nil)
			w := httptest.NewRecorder()

			// Создаем router context
			rctx := chi.NewRouteContext()
			// Добавляем URLParam urlID
			rctx.URLParams.Add("urlID", test.target)
			// Добавляем контекст к запросу
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Создаем MemoryStorage и записываем туда значения
			repo := memory.NewStorage()
			if test.id != "" {
				err := repo.Set(test.id, test.url)
				require.NoError(t, err)
			}

			// Выполняем запрос
			GetOriginalURL(repo).ServeHTTP(w, req)

			// Получаем результат запроса
			res := w.Result()

			// Проверяем HTTP Статус код
			assert.Equal(t, test.code, res.StatusCode)

			// Проверяем заголовок `Location`
			assert.Equal(t, test.url, res.Header.Get("Location"))

			// Получаем тело ответа и проверяем его
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.body, string(body))
		})
	}
}

func TestHTTPOk_GetSorterURL(t *testing.T) {
	const (
		name        = "Получение короткого URL"
		contentType = "text/plain"
		urlOrigin   = "http://test.aa"
		statusCode  = http.StatusCreated
	)

	t.Run(name, func(t *testing.T) {
		// Создаем Request
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(urlOrigin))
		// Устанавливаем "Content-Type" для Request
		req.Header.Set("Content-Type", contentType)
		// Создаем Recorder в который будет записываться ответ
		w := httptest.NewRecorder()

		// Создаем MemoryStorage
		repo := memory.NewStorage()

		// Создаем Config
		conf := config.Config{
			BaseHTTP: "https://shorter.123",
		}

		// Выполняем запрос
		GetSorterURL(repo, &conf).ServeHTTP(w, req)

		// Получаем ID соответствующий URL
		urlID, err := repo.GetByURL(urlOrigin)
		// Убеждаемся в успешном поиске URL ID
		require.NoError(t, err, "URL ID не найден")
		//Получаем короткий URL
		urlShort, err := url.JoinPath(conf.BaseHTTP, urlID)
		require.NoError(t, err)

		// Получаем результат запроса
		res := w.Result()

		// Проверяем HTTP Статус код
		assert.Equal(t, statusCode, res.StatusCode)

		// Получаем тело ответа и проверяем короткий URL в теле ответа
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, urlShort, string(body))
	})
}

func TestHTTPError_GetSorterURL(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		url         string
		code        int
	}{
		{
			name:        "Не верный Content-Type",
			contentType: "text/html",
			body:        "Разрешен только \"Content-Type: text/plain\"\n",
			url:         "http://test.aa",
			code:        http.StatusBadRequest,
		},
		{
			name:        "Не валидный URL",
			contentType: "text/plain",
			body:        "В данных запроса ожидаться валидный URL\n",
			url:         "http?test.aa",
			code:        http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем Request
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.url))
			// Устанавливаем "Content-Type" для Request
			req.Header.Set("Content-Type", test.contentType)
			// Создаем Recorder в который будет записываться ответ
			w := httptest.NewRecorder()

			// Создаем MemoryStorage
			repo := memory.NewStorage()

			// Создаем Config
			conf := config.Config{}

			// Выполняем запрос
			GetSorterURL(repo, &conf).ServeHTTP(w, req)

			// Получаем результат запроса
			res := w.Result()

			// Проверяем HTTP Статус код
			assert.Equal(t, test.code, res.StatusCode)

			// Получаем тело ответа и проверяем его
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.body, string(body))
		})
	}
}
