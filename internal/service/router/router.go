package router

import (
	"github.com/AleksandrTitov/shortener/internal/handler"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(repo repository.Repository) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/{urlID}", handler.GetOriginalURL(repo))
	router.Post("/", handler.GetSorterURL(repo))

	return router
}
