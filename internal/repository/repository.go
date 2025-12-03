package repository

import "errors"

var (
	ErrorAlreadyExist = errors.New("запись уже существует")
	ErrorNotFound     = errors.New("запись не найдена")
)

type Repository interface {
	Get(id string) (string, bool)
	Set(id, url string) error
	GetAll() map[string]string
	Unic(id string) bool
	Delete(id string) bool
	GetByURL(url string) (string, error)
}
