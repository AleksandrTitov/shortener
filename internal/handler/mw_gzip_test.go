package handler

import (
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MiddlewareGzipWrite(t *testing.T) {
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
	MiddlewareGzipWrite(testHandler).ServeHTTP(w, req)

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
