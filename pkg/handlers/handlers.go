package handlers

import (
	"net/http"
	"github.com/sanijo/booking-app/pkg/config"
	"github.com/sanijo/booking-app/pkg/handlers/models"
	"github.com/sanijo/booking-app/pkg/render"
)

// Repository is the repository type
type Repository struct {
    App *config.AppConfig
}

// Repo repository used by the handlers
var Repo *Repository

// NewRepo sets a new repository
func NewRepo(a *config.AppConfig) *Repository {
    return &Repository {
        App: a,
    }
}

// NewHandlers just sets the Repo variable for the handlers
func NewHandlers(r *Repository) {
    Repo = r
}

// Home is homepage handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// Model3 is model-3 page handler
func (m *Repository) Model3(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "model-3.page.html", &models.TemplateData{})
}

// ModelY is model-y page handler
func (m *Repository) ModelY(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "model-y.page.html", &models.TemplateData{})
}

// CheckAvailability is check-availability page handler
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "check_availability.page.html", &models.TemplateData{})
}

// Rent is rent page handler
func (m *Repository) Rent(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "rent.page.html", &models.TemplateData{})
}

// About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
    // send the data to the RenderTemplate
    render.RenderTemplate(w, "about.page.html", &models.TemplateData{})
}

// Contact is contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, "contact.page.html", &models.TemplateData{})
}
