package main

import (
	"github.com/AleksandrTitov/shortener/internal/handler"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	ms := memory.NewInMemoryStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.GetSorterURL(ms))
	mux.HandleFunc("/{urlID}", handler.GetOriginalURL(ms))

	err := http.ListenAndServe("localhost:8089", mux)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
