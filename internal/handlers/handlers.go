package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/driver"
	"github.com/sanijo/rent-app/internal/forms"
	"github.com/sanijo/rent-app/internal/helpers"
	"github.com/sanijo/rent-app/internal/models"
	"github.com/sanijo/rent-app/internal/render"
	"github.com/sanijo/rent-app/internal/repository"
	"github.com/sanijo/rent-app/internal/repository/dbrepo"
)

// Repository is the repository type
type Repository struct {
    App *config.AppConfig
    DB repository.DatabaseRepo
}

// Repo repository used by the handlers
var Repo *Repository

// NewRepo sets a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
    return &Repository {
        App: a,
        DB: dbrepo.NewPostgresRepo(db.SQL, a),
    }
}

// NewHandlers just sets the Repo variable for the handlers
func NewHandlers(r *Repository) {
    Repo = r
}

// Home is homepage handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// Model3 is model-3 page handler
func (m *Repository) Model3(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "model-3.page.html", &models.TemplateData{})
}

// ModelY is model-y page handler
func (m *Repository) ModelY(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "model-y.page.html", &models.TemplateData{})
}

// CheckAvailability is check-availability page handler
func (m *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "check-availability.page.html", &models.TemplateData{})
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
// PostAvailabilityJSON handles request for availability and sends JSON
// response
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
    resp := jsonResponse {
        OK: false,
        Message: "Available",
    }

    out, err := json.MarshalIndent(resp, "", "    ")
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    fmt.Println(string(out))
    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}

// Rent is rent page handler
func (m *Repository) Rent(w http.ResponseWriter, r *http.Request) {
    var emptyRent models.Rent
    data := make(map[string]interface{})
    data["rent"] = emptyRent

    render.Template(w, r, "rent.page.html", &models.TemplateData{
        Form: forms.New(nil),
        Data: data,
    })
}

// PostRent handles the posting of a rent form
func (m *Repository) PostRent(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    rent := models.Rent{
        FirstName: r.Form.Get("first_name"),
        LastName: r.Form.Get("last_name"),
        Email: r.Form.Get("email"),
        Phone: r.Form.Get("phone"),
    }

    form := forms.New(r.PostForm)
    form.Required("first_name", "last_name", "email")
    form.MinLength("first_name", 2)
    form.IsEmail("email")

    if !form.Valid() {
        data := make(map[string]interface{})
        data["rent"] = rent

        render.Template(w, r, "rent.page.html", &models.TemplateData{
            Form: form,
            Data: data,
        })
        return
    }

    // store rent value into session (type enabled in main)
    m.App.Session.Put(r.Context(), "rent", rent)
    http.Redirect(w, r, "/rent-summary", http.StatusSeeOther)
}

// About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Contact is contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
    render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// RentSummary is rent-summary page handler
func (m *Repository) RentSummary(w http.ResponseWriter, r *http.Request) {
    rent, ok := m.App.Session.Get(r.Context(), "rent").(models.Rent)
    if !ok {
        m.App.ErrorLog.Println("Cannot get item from session")
        m.App.Session.Put(r.Context(), "error", "Can't get rent from session")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    m.App.Session.Remove(r.Context(), "rent")

    data := make(map[string]interface{})
    data["rent"] = rent

    render.Template(w, r, "rent-summary.page.html", &models.TemplateData{
        Data: data,
    })
}
