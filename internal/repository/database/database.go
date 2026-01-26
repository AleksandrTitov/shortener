package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/repository"
)

type Storage struct {
	context context.Context
	db      *sql.DB
}

func NewStorage(ctx context.Context, db *sql.DB) repository.Repository {
	return &Storage{
		context: ctx,
		db:      db,
	}
}

func (s *Storage) Get(id string) (string, bool) {
	var originalURL string

	row := s.db.QueryRowContext(s.context, "select original_url from public.shorter where url_id=$1", id)

	err := row.Scan(&originalURL)
	if err != nil {
		logger.Log.Errorf("Ошибка получения ID: %v", err)
		return "", false
	}

	return originalURL, true
}

func (s *Storage) Set(id, url, userID string) error {
	result, err := s.db.ExecContext(
		s.context,
		"INSERT INTO public.shorter (url_id, original_url, user_id) VALUES ($1, $2, $3)", id, url, userID,
	)
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

func (s *Storage) SetBatch(urls map[string]string, userID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for id, url := range urls {
		_, err = tx.ExecContext(
			s.context,
			"INSERT INTO public.shorter (url_id, original_url) VALUES ($1, $3)", id, url, userID,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAll() [][]string {
	return nil
}

func (s *Storage) Unic(id string) bool {
	return false
}

func (s *Storage) Delete(id string) bool {
	return false
}

func (s *Storage) GetByURL(url string) (string, error) {
	var urlID string

	row := s.db.QueryRowContext(s.context, "select url_id from public.shorter where original_url=$1", url)

	err := row.Scan(&urlID)
	if err != nil {
		return "", err
	}

	return urlID, nil
}

func (s *Storage) GetByUserID(userID string) ([]repository.UsersURL, error) {
	var UsersURLs []repository.UsersURL

	row, err := s.db.QueryContext(s.context, "select url_id, original_url from public.shorter where user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		usersURL := repository.UsersURL{}
		err = row.Scan(
			&usersURL.URLID,
			&usersURL.OriginalURL,
		)
		if err != nil {
			return nil, err
		}
		UsersURLs = append(UsersURLs, usersURL)
	}

	return UsersURLs, nil
}
