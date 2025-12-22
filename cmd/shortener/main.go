package main

import (
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/file"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/AleksandrTitov/shortener/internal/router"
	"net/http"
)

func main() {
	conf := config.NewConfig()
	err := logger.Initialize(conf.LogLevel, true)
	if err != nil {
		logger.Log.Warnf("Не удалось инициализировать логгер: %v", err.Error())
	}
	logger.Log.Infof("Адрес: %s, базовый http: %s", conf.Addr, conf.BaseHTTP)

	stor := memory.NewStorage()

	if conf.FileName != "" {
		logger.Log.Infof("Загружаем данные из файла %s", conf.FileName)
		items := file.NewShorterItems()
		data, err := items.LoadShorterItems(conf.FileName)
		if err != nil {
			logger.Log.Warnf("Ошибка чтения файла данных: %v", err)
		} else {
			n := 0
			for _, i := range *data {
				err = stor.Set(i.ShortURL, i.OriginalURL)
				if err != nil {
					logger.Log.Warnf("Ошибка записи в хранилище: %v", err)
				} else {
					n += 1
				}
			}
			logger.Log.Debugf("Загружено записей: %d", n)
		}
	}

	gen := id.NewGenerator()
	r := router.NewRouter(stor, conf, gen)

	err = http.ListenAndServe(conf.Addr, r)
	if err != nil {
		logger.Log.Errorf("Не удалось запустить сервер: %v", err)
	}
}
