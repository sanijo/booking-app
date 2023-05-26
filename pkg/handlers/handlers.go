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
    
    // Store remote_ip in the session
    remote_ip := r.RemoteAddr
    m.App.Session.Put(r.Context(), "remote_ip", remote_ip)
    
    render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

    // perform som business logic here
    stringMap := make(map[string]string)
    stringMap["test"] = "Hello again"

    remote_ip := m.App.Session.GetString(r.Context(), "remote_ip")
    stringMap["remote_ip"] = remote_ip

    // send the data to the RenderTemplate
    render.RenderTemplate(w, "about.page.html", &models.TemplateData{
        StringMap: stringMap,
    })
}

