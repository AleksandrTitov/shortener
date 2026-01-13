package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/repository"
)

type Storage struct {
	Context context.Context
	DB      *sql.DB
}

func NewStorage() repository.Repository {
	return &Storage{}
}

func (s *Storage) Get(id string) (string, bool) {
	var originalURL string

	row := s.DB.QueryRowContext(s.Context, "select original_url from public.shorter where url_id=$1", id)

	err := row.Scan(&originalURL)
	if err != nil {
		logger.Log.Errorf("Ошибка получения ID: %v", err)
		return "", false
	}

	return originalURL, true
}

func (s *Storage) Set(id, url string) error {
	result, err := s.DB.ExecContext(s.Context, "INSERT INTO public.shorter (url_id, original_url) VALUES ($1, $2)", id, url)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("количество измененных строк должно ровняться 1, текущее: %d", rows)
	}

	return nil
}

func (s *Storage) GetAll() map[string]string {
	return nil
}

func (s *Storage) Unic(id string) bool {
	return false
}

func (s *Storage) Delete(id string) bool {
	return false
}

func (s *Storage) GetByURL(url string) (string, error) {
	return "", nil
}
