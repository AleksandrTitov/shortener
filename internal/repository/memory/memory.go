package memory

import (
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/repository"
)

type Storage struct {
	Store map[string]StoreData
}

type StoreData struct {
	OriginalURL string
	UserID      string
}

func NewStorage() repository.Repository {
	return &Storage{
		Store: make(map[string]StoreData),
	}
}

func (s *Storage) Set(id, url string, userID string) error {
	_, ok := s.Store[id]
	if !ok {
		s.Store[id] = StoreData{
			OriginalURL: url,
			UserID:      userID,
		}
		return nil
	}

	return fmt.Errorf("%w: %s", repository.ErrorAlreadyExist, id)
}

func (s *Storage) SetBatch(urls map[string]string, userID string) error {
	for id, url := range urls {
		_, ok := s.Store[id]
		if ok {
			return fmt.Errorf("%w: %s", repository.ErrorAlreadyExist, id)
		}
		s.Store[id] = StoreData{
			OriginalURL: url,
			UserID:      userID,
		}
	}

	return nil
}

func (s *Storage) Get(id string) (string, bool, bool) {
	data, ok := s.Store[id]
	return data.OriginalURL, ok, false
}

func (s *Storage) GetAll() [][]string {
	var data [][]string
	for k, v := range s.Store {
		item := make([]string, 3)
		item[0] = k
		item[1] = v.OriginalURL
		item[2] = v.UserID

		data = append(data, item)
	}

	return data
}

func (s *Storage) Unic(id string) bool {
	_, ok := s.Store[id]
	return !ok
}

func (s *Storage) Delete(id string) bool {
	_, ok := s.Store[id]
	if ok {
		delete(s.Store, id)
	}
	return ok
}

func (s *Storage) GetByURL(url string) (string, error) {
	for k, v := range s.Store {
		if v.OriginalURL == url {
			return k, nil
		}
	}
	return "", fmt.Errorf("%w: %s", repository.ErrorNotFound, url)
}

func (s *Storage) GetByUserID(userID string) ([]repository.UsersURL, error) {
	var UsersURLs []repository.UsersURL

	for k, v := range s.Store {
		if v.UserID == userID {
			usersURL := repository.UsersURL{
				OriginalURL: v.OriginalURL,
				URLID:       k,
			}
			UsersURLs = append(UsersURLs, usersURL)
		}
	}
	return UsersURLs, nil
}

func (s *Storage) DeleteIDs(userID string, id []string) error {
	return nil
}
