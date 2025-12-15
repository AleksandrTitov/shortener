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
	conf := config.NewConfig()
	err := logger.Initialize(conf.LogLevel, true)
	if err != nil {
		logger.Log.Warnf(err.Error())
	}
	logger.Log.Infof("адрес: %s, базовый http: %s", conf.Addr, conf.BaseHTTP)

	stor := memory.NewStorage()
	gen := id.NewGenerator()
	r := router.NewRouter(stor, conf, gen)

	err = http.ListenAndServe(conf.Addr, r)
	if err != nil {
		logger.Log.Errorf("Не удалось запустить сервер: %v", err)
	}
}
