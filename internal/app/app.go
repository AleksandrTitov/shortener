package app

import (
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/file"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/AleksandrTitov/shortener/internal/router"
	"github.com/AleksandrTitov/shortener/migrations"
	"net/http"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	server, err := a.CreateServer()
	if err != nil {
		return err
	}
	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) CreateServer() (*http.Server, error) {
	conf := config.NewConfig()
	err := logger.Initialize(conf.LogLevel, true)
	if err != nil {
		logger.Log.Warnf("Не удалось инициализировать логгер: %v", err.Error())
	}
	logger.Log.Infof("Адрес: %s, базовый http: %s", conf.Addr, conf.BaseHTTP)

	err = migrations.MigrateUP(conf.DatabaseDSN)
	if err != nil {
		logger.Log.Warnf("Не удалось запустить миграцию: %v", err.Error())
	}

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

	server := &http.Server{
		Addr:    conf.Addr,
		Handler: r,
	}

	return server, nil
}
