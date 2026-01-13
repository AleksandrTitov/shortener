package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MigrateUP(db *sql.DB) error {
	// Создаем драйвер для миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать драйвер миграций: %w", err)
	}

	// Создаем объект миграции
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр миграции: %w", err)
	}

	// Выполняем миграции
	if err := m.Up(); err != nil {
		// Игнорируем ошибку "no change" - это нормально
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Info("Нет миграций для применения")
			return nil
		}
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	logger.Log.Info("Миграции успешно применены")
	return nil
}
