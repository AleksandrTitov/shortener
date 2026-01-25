package handler

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func getURLID(urlOrigin, userID string, repo repository.Repository, gen id.GeneratorID) (string, error) {
	maxAttempts := 15
	var urlID string
	var err error
	var pgErr *pgconn.PgError

	for i := 1; i <= maxAttempts; i++ {
		urlID, err = gen.GetID()
		if errors.Is(err, id.ErrGetID) {
			logger.Log.Errorf("Не удалось сгенерировать ID: %v", err.Error())
			break
		}

		logger.Log.Debugf("Записываем url oring: '%s', url id: '%s'", urlOrigin, urlID)
		err = repo.Set(urlID, urlOrigin, userID)
		logger.Log.Debugf("Значение записано url oring: '%s', url id: '%s'", urlOrigin, urlID)

		if errors.Is(err, repository.ErrorAlreadyExist) {
			logger.Log.Warnf("Не удалось записать id '%s' попытка %d(%d), %v", urlID, i, maxAttempts, err.Error())
			continue
		} else if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			logger.Log.Debugf("Запись origin url '%s' уже существует, получаем записаный url id", urlOrigin)
			urlID, err = repo.GetByURL(urlOrigin)
			if err != nil {
				logger.Log.Errorf("Ошибка получения url id: %v", err.Error())
				return "", repository.ErrorGet
			}
			return urlID, repository.ErrorOriginNotUnique

		} else if err != nil {
			logger.Log.Errorf("Не ожиданная ошибка при записи id '%s': %v", urlID, err)
			break
		}
		break
	}

	return urlID, err
}
