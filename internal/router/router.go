package router

import (
	"database/sql"
	"github.com/AleksandrTitov/shortener/internal/config"
	"github.com/AleksandrTitov/shortener/internal/handler"
	"github.com/AleksandrTitov/shortener/internal/middleware"
	"github.com/AleksandrTitov/shortener/internal/model/id"
	"github.com/AleksandrTitov/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func NewRouter(repo repository.Repository, conf *config.Config, gen id.GeneratorID, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.GzipRead)
	router.Use(middleware.GzipWrite)
	router.Use(middleware.CookiesWrite)

	router.Get("/{urlID}", handler.GetOriginalURL(repo))
	router.Post("/", handler.GetSorterURL(repo, conf, gen))
	router.Post("/api/shorten", handler.GetSorterURLJson(repo, conf, gen))
	router.Post("/api/shorten/batch", handler.GetShorterURLJsonBatch(repo, conf, gen))
	router.Get("/ping", handler.Ping(db))
	return router
}
