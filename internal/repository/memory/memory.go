package memory

import (
	"fmt"
	"github.com/AleksandrTitov/shortener/internal/repository"
)

type Storage struct {
	Store map[string]string
}

func NewStorage() repository.Repository {
	return &Storage{
		Store: make(map[string]string),
	}
}

func (s *Storage) Set(id, url string) error {
	_, ok := s.Store[id]
	if !ok {
		s.Store[id] = url
		return nil
	}

	return fmt.Errorf("%w: %s", repository.ErrorAlreadyExist, id)
}

func (s *Storage) Get(id string) (string, bool) {
	url, ok := s.Store[id]
	return url, ok
}

func (s *Storage) GetAll() map[string]string {
	return s.Store
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
	for k, v := range s.GetAll() {
		if v == url {
			return k, nil
		}
	}
	return "", fmt.Errorf("%w: %s", repository.ErrorNotFound, url)
}
