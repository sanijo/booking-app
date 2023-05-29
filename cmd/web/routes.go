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
    mux.Get("/model-3", handlers.Repo.Model3)
    mux.Get("/model-y", handlers.Repo.ModelY)
    mux.Get("/check-availability", handlers.Repo.CheckAvailability)
    mux.Get("/rent", handlers.Repo.Rent)
    mux.Get("/about", handlers.Repo.About)
    mux.Get("/contact", handlers.Repo.Contact)

    // In static folder are all things that are not html template such as JS,
    // figures
    filseServer := http.FileServer(http.Dir("./static/"))
    mux.Handle("/static/*", http.StripPrefix("/static", filseServer))

    return mux
} 
