package handler

import (
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"io"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name   string
		method string
		target string
		id     string
		url    string
		body   string
		code   int
	}{
		{
			name:   "Получение оригинального URL",
			method: http.MethodGet,
			target: "dmiWnD",
			id:     "dmiWnD",
			url:    "http://test.aa",
			code:   http.StatusTemporaryRedirect,
		},
		{
			name:   "ID не найден",
			method: http.MethodGet,
			target: "dmiWnD",
			body:   "ID \"dmiWnD\" не найден\n",
			code:   http.StatusBadRequest,
		},
		{
			name:   "Не верный метод",
			method: http.MethodPut,
			body:   "Разрешен только метод GET\n",
			code:   http.StatusMethodNotAllowed,
		},
		{
			name:   "Длина запрашиваемого ID больше длины формата",
			method: http.MethodGet,
			target: "dmiWnDDD",
			body:   fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code:   http.StatusBadRequest,
		},
		{
			name:   "Длина запрашиваемого ID меньше длины формата",
			method: http.MethodGet,
			target: "dmi",
			body:   fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code:   http.StatusBadRequest,
		},
		{
			name:   "Пустой ID",
			method: http.MethodGet,
			body:   fmt.Sprintf("Длина ID должна быть равна %d символам\n", id.LenID),
			code:   http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, fmt.Sprintf("/%s", test.target), nil)
			w := httptest.NewRecorder()

			// Создаем MemoryStorage и записываем туда значения
			repo := memory.NewInMemoryStorage()
			err := repo.Set(test.id, test.url)
			require.NoError(t, err)

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
		method      = http.MethodPost
		contentType = "text/plain"
		url         = "http://test.aa"
		statusCode  = http.StatusCreated
	)

	t.Run(name, func(t *testing.T) {
		req := httptest.NewRequest(method, "/", strings.NewReader(url))
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		repo := memory.NewInMemoryStorage()
		GetSorterURL(repo).ServeHTTP(w, req)

		// Получаем URL ID соответствующий test.url
		var urlID string
		for k, v := range repo.GetAll() {
			if v == url {
				urlID = k
			}
		}
		// Убеждаемся в успешном поиске URL ID
		require.NotEmpty(t, urlID, "URL ID не найден")

		res := w.Result()

		// Проверяем HTTP Статус код
		assert.Equal(t, statusCode, res.StatusCode)

		// Получаем тело ответа и проверяем короткий URL в теле ответа
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("http://%s/%s", req.Host, urlID), string(body))
	})
}

func TestHTTPError_GetSorterURL(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		id          string
		url         string
		code        int
	}{
		{
			name:        "Не верный Content-Type",
			method:      http.MethodPost,
			contentType: "text/html",
			body:        "Разрешен только \"Content-Type: text/plain\"\n",
			url:         "http://test.aa",
			code:        http.StatusBadRequest,
		},
		{
			name:        "Не верный метод",
			method:      http.MethodGet,
			contentType: "text/plain",
			body:        fmt.Sprintf("Разрешен только метод %s\n", http.MethodPost),
			url:         "http://test.aa",
			code:        http.StatusMethodNotAllowed,
		},
		{
			name:        "Не валидный URL",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        "В данных запроса ожидаться валидный URL\n",
			url:         "http?test.aa",
			code:        http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем Request
			req := httptest.NewRequest(test.method, "/", strings.NewReader(test.url))
			// Устанавливаем "Content-Type" для Request
			req.Header.Set("Content-Type", test.contentType)
			// Создаем Recorder в который будет записываться ответ
			w := httptest.NewRecorder()

			// Создаем MemoryStorage
			repo := memory.NewInMemoryStorage()

			// Выполняем запрос
			GetSorterURL(repo).ServeHTTP(w, req)

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
