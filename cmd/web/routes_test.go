package main

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sanijo/rent-app/internal/config"
)

func TestRoutes(t *testing.T) {
	app := config.AppConfig{} // Create a sample AppConfig

	want := chi.NewRouter()

	got := routes(&app)

	// Check if the types of the objects are equal
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("routes() = %T, want %T", got, want)
	}
}
