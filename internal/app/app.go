package app

import (
	"context"
	"database/sql"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/file"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/AleksandrTitov/shortener/internal/repository/database"
	"github.com/AleksandrTitov/shortener/internal/repository/memory"
	"github.com/AleksandrTitov/shortener/internal/router"
	"github.com/AleksandrTitov/shortener/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

type App struct {
	DB *sql.DB
}

func NewApp() *App {
	return &App{
		DB: nil,
	}
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

	stor, err := a.CreateStorage(conf)
	if err != nil {
		logger.Log.Errorf("Ошибка создания хранилища: %v", err)
		return nil, err
	}

	gen := id.NewGenerator()
	r := router.NewRouter(stor, conf, gen, a.DB)

	server := &http.Server{
		Addr:    conf.Addr,
		Handler: r,
	}

	return server, nil
}

func (a *App) CreateStorage(conf *config.Config) (repository.Repository, error) {
	switch {
	case conf.DatabaseDSN != "":
		logger.Log.Info("Используемое хранилище: база данных")
		return a.createDatabaseStorage(conf.DatabaseDSN)
	case conf.FileName != "":
		logger.Log.Info("Используемое хранилище: файловая система")
		return a.createFileStorage(conf.FileName)
	default:
		logger.Log.Info("Используемое хранилище: оперативная память")
		return a.createMemoryStorage()
	}
}

func (a *App) createDatabaseStorage(dsn string) (repository.Repository, error) {
	db, err := sql.Open("pgx", dsn)
	a.DB = db
	if err != nil {
		logger.Log.Errorf("Не удалось установить соединение с базой данных: %v", err)
		return nil, err
	}

	err = migrations.MigrateUP(db)
	if err != nil {
		logger.Log.Warnf("Не удалось запустить миграцию: %v", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		logger.Log.Errorf("Не удалось подключиться к базе данных: %v", err)
		return nil, err
	}

	stor := database.NewStorage(ctx, db)

	return stor, nil
}

func (a *App) createFileStorage(filename string) (repository.Repository, error) {
	logger.Log.Infof("Загружаем данные из файла %s", filename)
	items := file.NewShorterItems()
	stor := memory.NewStorage()
	data, err := items.LoadShorterItems(filename)
	if err != nil {
		logger.Log.Warnf("Ошибка чтения файла данных: %v", err)
	} else {
		n := 0
		for _, i := range *data {
			err = stor.Set(i.ShortURL, i.OriginalURL, i.UserID)
			if err != nil {
				logger.Log.Warnf("Ошибка записи в хранилище: %v", err)
			} else {
				n += 1
			}
		}
		logger.Log.Debugf("Загружено записей: %d", n)
	}

	return stor, nil
}

func (a *App) createMemoryStorage() (repository.Repository, error) {
	stor := memory.NewStorage()
	return stor, nil
}
