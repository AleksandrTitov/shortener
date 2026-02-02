package repository

import "errors"

var (
	ErrorAlreadyExist    = errors.New("запись уже существует")
	ErrorOriginNotUnique = errors.New("запись origin url не уникальна")
	ErrorNotFound        = errors.New("запись url id не найдена")
	ErrorUserNotFound    = errors.New("пользовательский id не найден")
	ErrorGet             = errors.New("ошибка получения значения")
)

type UsersURL struct {
	OriginalURL string `json:"original_url"`
	URLID       string `json:"short_url"`
}

type Repository interface {
	Get(id string) (string, bool, bool)
	GetByURL(url string) (string, error)
	GetAll() [][]string
	GetByUserID(userID string) ([]UsersURL, error)

	Set(id, url, userID string) error
	SetBatch(urls map[string]string, userID string) error

	Unic(id string) bool
	Delete(id string) bool
	DeleteIDs(userID string, id []string) error
}
