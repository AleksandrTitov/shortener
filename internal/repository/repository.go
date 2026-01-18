package repository

import "errors"

var (
	ErrorAlreadyExist    = errors.New("запись уже существует")
	ErrorOriginNotUnique = errors.New("запись origin url не уникальна")
	ErrorNotFound        = errors.New("запись url id не найдена")
	ErrorGet             = errors.New("ошибка получения значения")
)

type Repository interface {
	Get(id string) (string, bool)
	Set(id, url string) error
	SetBatch(map[string]string) error
	GetAll() map[string]string
	Unic(id string) bool
	Delete(id string) bool
	GetByURL(url string) (string, error)
}
