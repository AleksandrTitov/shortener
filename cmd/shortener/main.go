package main

import (
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/AleksandrTitov/shortener/internal/service/router"
	"log"
	"net/http"
)

func main() {
	s := memory.NewStorage()
	r := router.NewRouter(s)
	cfg := config.NewConfig()

	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
