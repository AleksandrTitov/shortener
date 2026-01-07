package database

import (
	"context"
	"database/sql"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func Ping(config *config.Config) error {
	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		logger.Log.Error("Не удается подключиться к базе данных")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		logger.Log.Errorf("Не удалось выполнить проверку подключения к БД: %v", err)
		return err
	}
	logger.Log.Debugf("Проверка подключения к БД выполнена успешно")

	return nil
}
