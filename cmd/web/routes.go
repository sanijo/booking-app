package main

import (
	"net/http"
	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
    mux := chi.NewRouter()

    mux.Use(middleware.Recoverer)
    // cross-site request forgery protection
    mux.Use(NoSurf) 
    mux.Use(SessionLoad)

    mux.Get("/", handlers.Repo.Home)
    mux.Get("/model-3", handlers.Repo.Model3)
    mux.Get("/model-y", handlers.Repo.ModelY)

    mux.Get("/check-availability", handlers.Repo.CheckAvailability)
    mux.Post("/check-availability", handlers.Repo.PostAvailability)
    mux.Post("/check-availability-json", handlers.Repo.PostAvailabilityJSON)

    mux.Get("/rent", handlers.Repo.Rent)
    mux.Post("/rent", handlers.Repo.PostRent)
    mux.Get("/rent-summary", handlers.Repo.RentSummary)

    mux.Get("/about", handlers.Repo.About)
    mux.Get("/contact", handlers.Repo.Contact)

    // In static folder are all things that are not html template such as JS,
    // figures
    filesServer := http.FileServer(http.Dir("./static/"))
    mux.Handle("/static/*", http.StripPrefix("/static", filesServer))

    return mux
} 
