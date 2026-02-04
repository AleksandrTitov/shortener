package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/lib/pq"
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

func (s *Storage) Get(id string) (string, error) {
	var originalURL string
	var gone bool

	row := s.db.QueryRowContext(s.context, "select original_url, deleted_flag from public.shorter where url_id=$1", id)

	err := row.Scan(&originalURL, &gone)
	if err != nil {
		logger.Log.Errorf("Ошибка получения ID: %v", err)
		return "", err
	}

	if gone {
		return "", repository.ErrorGone
	}
	return originalURL, nil
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
			"INSERT INTO public.shorter (url_id, original_url, user_id) VALUES ($1, $2, $3)", id, url, userID,
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

	rows, err := s.db.QueryContext(s.context, "select url_id, original_url from public.shorter where user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	for rows.Next() {
		usersURL := repository.UsersURL{}
		err = rows.Scan(
			&usersURL.URLID,
			&usersURL.OriginalURL,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных пользователя: %w", err)
		}
		UsersURLs = append(UsersURLs, usersURL)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов запроса пользователей: %w", err)
	}

	return UsersURLs, nil
}

func (s *Storage) DeleteIDs(userID string, urlsID []string) error {
	logger.Log.Debugf("Получено на удаление %d записей пользователя %s", len(urlsID), userID)
	result, err := s.db.ExecContext(
		s.context,
		"update public.shorter set deleted_flag = true where url_id=any($1) and user_id=$2", pq.Array(urlsID), userID,
	)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользовательских url id: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Log.Warnf("Не удалось получить количество обновленных строк: %v", err)
		return nil
	}

	logger.Log.Debugf("Удалено %d записей пользователя %s", rows, userID)

	return nil
}
