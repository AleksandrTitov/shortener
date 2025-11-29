package memory

import (
	"github.com/AleksandrTitov/shortener/internal/repository"
)

type InMemoryStorage struct {
	Store map[string]string
}

func NewInMemoryStorage() repository.Repository {
	return &InMemoryStorage{
		Store: make(map[string]string),
	}
}

func (ms *InMemoryStorage) Set(id, url string) error {
	ms.Store[id] = url
	return nil
}

func (ms *InMemoryStorage) Get(id string) (string, bool) {
	url, ok := ms.Store[id]
	return url, ok
}

func (ms *InMemoryStorage) GetAll() map[string]string {
	return ms.Store
}

func (ms *InMemoryStorage) Unic(id string) bool {
	_, ok := ms.Store[id]
	return !ok
}
