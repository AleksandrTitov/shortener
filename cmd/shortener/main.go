package main

import (
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/AleksandrTitov/shortener/internal/service/router"
	"net/http"
)

func main() {
	log := logger.NewLogger()

	conf := config.NewConfig()
	stor := memory.NewStorage()
	gen := id.NewGenerator()
	r := router.NewRouter(stor, conf, gen)

	err := http.ListenAndServe(conf.Addr, r)
	if err != nil {
		log.Errorf("Не удалось запустить сервер: %v", err)
	}
}
