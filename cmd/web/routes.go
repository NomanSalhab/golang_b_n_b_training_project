package main

import (
	"net/http"

	"github.com/NomanSalhab/golang_b_n_b_training_project/pkg/config"
	"github.com/NomanSalhab/golang_b_n_b_training_project/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer) // * For When The Application Panics

	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
