package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GzipWrite(t *testing.T) {
	inData := `{"url":"http://test.aa"}`

	// Создаем Handler
	testHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write([]byte(inData))
		assert.NoError(t, err)
	})

	// Создаем Request
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	// Создаем Recorder
	w := httptest.NewRecorder()

	// Запускаем Handler с использованием MW функции сжимающей ответ
	GzipWrite(testHandler).ServeHTTP(w, req)

	// Получаем ответ и распаковываем его
	gz, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gz.Close()

	outData, err := io.ReadAll(gz)
	assert.NoError(t, err)

	// Сравниваем распакованный ответ и данные из Request
	assert.Equal(t, inData, string(outData))

	// Проверяем заголовок Content-Encoding
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}

func Test_GzipRead(t *testing.T) {
	// Тестовые данные
	inData := `{"url":"http://test.aa"}`

	// Создаем буфер в которую будут записываться сжатые данные
	var b bytes.Buffer

	// Создаем gzip writer для записи данных в буфер
	gzw := gzip.NewWriter(&b)

	// Записываем тестовые данные в gzip writer, они будут храниться в буфере
	_, err := gzw.Write([]byte(inData))
	assert.NoError(t, err)

	// Закрываем gzip writer
	err = gzw.Close()
	assert.NoError(t, err)

	// Создаем Handler
	testHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		// Проверяем, что данные распаковались
		assert.Equal(t, inData, string(body))
		t.Logf("Данные запроса: %s, полученые данные %s", inData, string(body))
	})

	// Создаем Request
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(b.Bytes()))
	// Устанавливаем заголовки
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	// Создаем Recorder
	w := httptest.NewRecorder()

	// Запускаем Handler с использованием MW функции распаковывающей запрос
	GzipRead(testHandler).ServeHTTP(w, req)
}
