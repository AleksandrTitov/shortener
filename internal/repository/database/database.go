package database

import "github.com/AleksandrTitov/shortener/internal/repository"

type Storage struct {
	ShortURL    string
	OriginalURL string
	ID          string
}

func NewStorage() repository.Repository {
	return &Storage{}
}

func (s *Storage) Get(id string) (string, bool) {
	return "", false
}

func (s *Storage) Set(id, url string) error {
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
