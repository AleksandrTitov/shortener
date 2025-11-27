package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const (
	lenID = 8
)

type URLShortener struct {
	store map[string]string
	// TODO: мьютекс
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		store: make(map[string]string),
	}
}

func getID(n int) string {
	const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	id := make([]byte, n)

	for i := range len(id) {
		id[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(id)
}

func validateURL(url string) error {
	//TODO: Имплементация валидации URL
	//return errors.New("Не является URL'ом")
	return nil
}

func (su URLShortener) unicID(id string) bool {
	_, ok := su.store[id]
	return ok
}

func (su URLShortener) sorterURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Разрешен только \"Content-Type: text/plain\"", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения данных запроса", http.StatusInternalServerError)
			return
		}
		url := string(body)

		err = validateURL(url)
		if err != nil {
			http.Error(w, "В данных запроса ожидаться валидный URL", http.StatusBadRequest)
			return
		}
		//TODO: проверка на уникальность
		id := getID(lenID)
		su.store[id] = url
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s\n", r.Host, id)))
		if err != nil {
			// TODO: что-то более внятное
			http.Error(w, "Ошибка ответа", http.StatusInternalServerError)
			return
		}
		fmt.Println(su.store)
		fmt.Println(len(su.store))

	case http.MethodGet:
		// RequestURI в начале имеет `/`, его убираем
		id := r.RequestURI[1:]
		if len(id) < lenID {
			http.Error(w, fmt.Sprintf("Длина ID должна быть равна %d символам", lenID), http.StatusBadRequest)
			return
		}
		url, ok := su.store[id]
		if !ok {
			http.Error(w, fmt.Sprintf("ID \"%s\" не найден", id), http.StatusBadRequest)
			return
		}
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		http.Error(w, "Разрешены только POST и GET методы!", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	su := NewURLShortener()

	mux := http.NewServeMux()
	mux.HandleFunc("/", su.sorterURL)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		// TODO: обработать ошибку
		fmt.Println(err.Error())
	}
}
