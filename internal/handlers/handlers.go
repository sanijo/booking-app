package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

// NewTestRepo sets a new test repository
func NewTestRepo(a *config.AppConfig) *Repository {
    return &Repository {
        App: a,
        DB: dbrepo.NewTestingRepo(a),
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

// PostAvailability is check-availability page handler. After user submits
// the form, this handler is called.
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
    // parse the form
    err := r.ParseForm()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Can't parse form")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    // get the form values
    start := r.Form.Get("start")
    end := r.Form.Get("end")

    // convert the date to time.Time type
    // 2020-01-01 -- 01/02 03:04:05PM '06 -0700
    layout := "2006-01-02"

    startDate, err := time.Parse(layout, start)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Can't parse start date")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    endDate, err := time.Parse(layout, end)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Can't parse end date")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    // get availability
    availableCarModels, err := m.DB.SearchAvailabilityForAllModels(startDate, endDate)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    // if slice is empty means no availability
    if len(availableCarModels) == 0 {
        m.App.Session.Put(r.Context(), "error", "No available vehicles for specified dates")
        http.Redirect(w, r, "/check-availability", http.StatusSeeOther)
        return
    }

    // create a map to store data to be sent to the template
    data := make(map[string]interface{}) 
    data["models"] = availableCarModels

    // create a rent struct to store data in session to be available in next page
    // gob.Register(models.Rent{}) allready in main.go therefore no need to register here
    // Idea: when user clicks on a model, start and end date are pulled from the
    // session and id of the chosen model is assigned to the rent struct and
    // stored in the session. Then rent page is rendered.
    rent := models.Rent{
        StartDate: startDate,
        EndDate: endDate,
    }

    m.App.Session.Put(r.Context(), "rent", rent)

    // send data to the template
    render.Template(w, r, "choose-model.page.html", &models.TemplateData{
        Data: data,
    })
}

type jsonResponse struct {
    OK bool `json:"ok"`
    Message string `json:"message"`
    ModelID string `json:"model_id"`
    StartDate string `json:"start_date"`
    EndDate string `json:"end_date"`
}

// PostAvailabilityJSON handles request for availability and sends JSON
// response
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
    sd := r.Form.Get("start")
    ed := r.Form.Get("end")

    layout := "2006-01-02"

    // convert the date to time.Time type
    startDate, _ := time.Parse(layout, sd)
    endDate, _ := time.Parse(layout, ed)
    
    modelID, _ := strconv.Atoi(r.Form.Get("model_id"))

    available, _ := m.DB.SearchAvailabilityByDatesAndModelID(startDate, endDate, modelID)

    resp := jsonResponse {
        OK: available,
        Message: "",
        ModelID: strconv.Itoa(modelID),
        StartDate: sd,
        EndDate: ed,
    }

    out, err := json.MarshalIndent(resp, "", "    ")
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}

// Rent is rent page handler
func (m *Repository) Rent(w http.ResponseWriter, r *http.Request) {
    // get rent struct from session. NOTE: Get returns an
    // interface, so we need to type assert it to models.Rent
    rent, ok := m.App.Session.Get(r.Context(), "rent").(models.Rent)
    if !ok {
        m.App.Session.Put(r.Context(), "error", "Can't get rent from session")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    // get model from database by model id
    model, err := m.DB.GetModelByID(rent.ModelID)
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Can't get model from database")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    
    // store model name into rent struct Model field
    rent.Model.ModelName = model.ModelName

    // store rent struct with model name into session
    m.App.Session.Put(r.Context(), "rent", rent)

    // convert the date to string type to be able to use it in rent template
    sd := rent.StartDate.Format("2006-01-02")
    ed := rent.EndDate.Format("2006-01-02")

    data := make(map[string]interface{})
    data["rent"] = rent

    // create string map (see TemplateData struct in models/models.go)
    // to store data to be sent to the template
    stringMap := make(map[string]string)    
    stringMap["start_date"] = sd
    stringMap["end_date"] = ed

    render.Template(w, r, "rent.page.html", &models.TemplateData{
        StringMap: stringMap,
        Form: forms.New(nil),
        Data: data,
    })
}

// PostRent handles the posting of a rent form
func (m *Repository) PostRent(w http.ResponseWriter, r *http.Request) {
    // get rent struct from session. NOTE: allready contains start date, end
    // date, model id and model name
    rent, ok := m.App.Session.Get(r.Context(), "rent").(models.Rent)
    if !ok {
        helpers.ServerError(w, errors.New("can't get rent from session"))
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    err := r.ParseForm()
    if err != nil {
        m.App.Session.Put(r.Context(), "error", "Can't parse form")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    // update rent struct with data from the form
    rent.FirstName = r.Form.Get("first_name")
    rent.LastName = r.Form.Get("last_name")
    rent.Email = r.Form.Get("email")
    rent.Phone = r.Form.Get("phone")

    // create a form struct to validate the data
    form := forms.New(r.PostForm)
    form.Required("first_name", "last_name", "email")
    form.MinLength("first_name", 2)
    form.IsEmail("email")

    // if there are any errors, redisplay the form
    if !form.Valid() {
        // convert the date to string type to be able to use it in rent-summary template
        sd := rent.StartDate.Format("2006-01-02")
        ed := rent.EndDate.Format("2006-01-02")

        // create string map (see TemplateData struct in models/models.go)
        // to store data to be sent to the template
        stringMap := make(map[string]string)
        stringMap["start_date"] = sd
        stringMap["end_date"] = ed

        data := make(map[string]interface{})
        data["rent"] = rent

        render.Template(w, r, "rent.page.html", &models.TemplateData{
            StringMap: stringMap,
            Form: form,
            Data: data,
        })
        return
    }

    // insert rent into database
    rentID, err := m.DB.InsertRent(rent)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    // create rent restriction struct
    rentRestriction := models.RentRestriction{
        StartDate: rent.StartDate,
        EndDate: rent.EndDate,
        ModelID: rent.ModelID,
        RentID: rentID,
        RestrictionID: 1,
    }

    // insert restriction into database
    err = m.DB.InsertRentRestriction(rentRestriction)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    // put rent value back into session (type enabled in main)
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

    // convert the date to string type to be able to use it in rent-summary template
    sd := rent.StartDate.Format("2006-01-02")
    ed := rent.EndDate.Format("2006-01-02")

    // create string map (see TemplateData struct in models/models.go)
    // to store data to be sent to the template
    stringMap := make(map[string]string)
    stringMap["start_date"] = sd
    stringMap["end_date"] = ed

    render.Template(w, r, "rent-summary.page.html", &models.TemplateData{
        StringMap: stringMap,
        Data: data,
    })
}

// ChooseModel is choose-model page handler
func (m *Repository) ChooseModel(w http.ResponseWriter, r *http.Request) {
    // get model id from url
    modelID, err := strconv.Atoi(chi.URLParam(r, "id")) 
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    // get rent from session and update modelID value. NOTE: Get returns an
    // interface, so we need to type assert it to models.Rent
    rent, ok := m.App.Session.Get(r.Context(), "rent").(models.Rent)
    if !ok {
        m.App.Session.Put(r.Context(), "error", "Can't get rent from session")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // update modelID value and save back into session
    rent.ModelID = modelID
    m.App.Session.Put(r.Context(), "rent", rent)

    // redirect to rent page
    http.Redirect(w, r, "/rent", http.StatusSeeOther)
}

// RentVehicle is rent-vehicle page handler
func (m *Repository) RentVehicle(w http.ResponseWriter, r *http.Request) {
    // grab id, s, and e values from url
    modelID, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    sd := r.URL.Query().Get("s")
    ed := r.URL.Query().Get("e")

    // convert the date to time.Time type to be able to use it in rent-vehicle template
    startDate, err := time.Parse("2006-01-02", sd)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }
    endDate, err := time.Parse("2006-01-02", ed)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }

    // get model from database
    model, err := m.DB.GetModelByID(modelID)
    if err != nil {
        helpers.ServerError(w, err)
        return
    }
    
    // create rent struct to store data to be sent to put into session
    var rent models.Rent

    rent.ModelID = modelID
    rent.Model.ModelName = model.ModelName
    rent.StartDate = startDate
    rent.EndDate = endDate

    // put rent value into session (type enabled in main)
    m.App.Session.Put(r.Context(), "rent", rent)

    // redirect to rent page
    http.Redirect(w, r, "/rent", http.StatusSeeOther)
}
