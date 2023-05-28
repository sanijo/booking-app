package main

import (
	"net/http"
	"github.com/sanijo/booking-app/pkg/config"
	"github.com/sanijo/booking-app/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
    mux := chi.NewRouter()

    mux.Use(middleware.Recoverer)
    mux.Use(NoSurf)
    mux.Use(SessionLoad)

    mux.Get("/", handlers.Repo.Home)
    mux.Get("/about", handlers.Repo.About)

    // In static folder are all things that are not html template such as JS,
    // figures
    filseServer := http.FileServer(http.Dir("./static/"))
    mux.Handle("/static/*", http.StripPrefix("/static", filseServer))

    return mux
} 
