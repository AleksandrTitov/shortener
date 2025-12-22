package handler

import (
	"errors"
	"github.com/AleksandrTitov/shortener/internal/logger"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
)

func getURLID(url string, repo repository.Repository, gen id.GeneratorID) (string, error) {
	maxAttempts := 15
	var urlID string
	var err error

	for i := 1; i <= maxAttempts; i++ {
		urlID, err = gen.GetID()
		if errors.Is(err, id.ErrGetID) {
			logger.Log.Errorf("Не удалось сгенерировать ID: %v", err.Error())
			break
		}
		err = repo.Set(urlID, url)
		if errors.Is(err, repository.ErrorAlreadyExist) {
			logger.Log.Warnf("Не удалось записать id \"%s\" попытка %d(%d), %v", urlID, i, maxAttempts, err.Error())
			continue
		} else if err != nil {
			logger.Log.Errorf("Не ожиданная ошибка при записи id \"%s\": %v", urlID, err)
			break
		}
		break
	}

	return urlID, err
}
