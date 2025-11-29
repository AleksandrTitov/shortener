package repository

type Repository interface {
	Get(id string) (string, bool)
	Set(id, url string) error
	GetAll() map[string]string
	Unic(id string) bool
}
