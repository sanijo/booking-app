package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sanijo/booking-app/internal/config"
	"github.com/sanijo/booking-app/internal/models"
	"github.com/sanijo/booking-app/internal/render"
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
    render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

// Model3 is model-3 page handler
func (m *Repository) Model3(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "model-3.page.html", &models.TemplateData{})
}

// ModelY is model-y page handler
func (m *Repository) ModelY(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "model-y.page.html", &models.TemplateData{})
}

// CheckAvailability is check-availability page handler
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "check-availability.page.html", &models.TemplateData{})
}

// PostkAvailability is check-availability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
    start := r.Form.Get("start")
    end := r.Form.Get("end")
    w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s",start,end)))
}

type jsonResponse struct {
    OK bool `json:"ok"`
    Message string `json:"message"`
}
// PostkAvailabilityJSON handles request for availability and sends JSON
// response
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
    resp := jsonResponse {
        OK: false,
        Message: "Available",
    }

    out, err := json.MarshalIndent(resp, "", "    ")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(out))
    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}

// Rent is rent page handler
func (m *Repository) Rent(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "rent.page.html", &models.TemplateData{})
}

// About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
    // send the data to the RenderTemplate
    render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{})
}

// Contact is contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}
