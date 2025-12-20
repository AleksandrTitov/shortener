package router

import (
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/handler"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func NewRouter(repo repository.Repository, conf *config.Config, gen id.GeneratorID) *chi.Mux {
	router := chi.NewRouter()
	router.Use(handler.MiddlewareLogging)
	router.Use(handler.MiddlewareGzipRead)
	router.Use(handler.MiddlewareGzipWrite)

	router.Get("/{urlID}", handler.GetOriginalURL(repo))
	router.Post("/", handler.GetSorterURL(repo, conf, gen))
	router.Post("/api/shorten", handler.GetSorterURLJson(repo, conf, gen))

	return router
}
