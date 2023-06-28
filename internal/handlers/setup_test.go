package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/models"
	"github.com/sanijo/rent-app/internal/render"
)


var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
    // What to put in session
    gob.Register(models.Rent{})

    // Change to true if in production
    app.InProduction = false

    session = scs.New()
    session.Lifetime = 24 * time.Hour
    session.Cookie.Persist = true
    session.Cookie.SameSite = http.SameSiteLaxMode
    session.Cookie.Secure = app.InProduction // in production: true

    // Set up info and error loggers
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    app.InfoLog = infoLog
    errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    app.ErrorLog = errorLog

    // Set pointer in config to session so that is available in program
    app.Session = session

    tc, err := CreateTestTemplateCache()
	if err != nil {
        log.Fatal("cannot create template cache")
	}
   
    // Setting TemplateCache in config so that is cached all the time while app
    // is running
    app.TemplateCache = tc
    app.UseCache = true

    repo := NewTestRepo(&app)
    NewHandlers(repo)
    
    // Give access to app config variable inside render package
    render.NewRenderer(&app)

    mux := chi.NewRouter()

    mux.Use(middleware.Recoverer)
    // cross-site request forgery protection (not necessary for testing)
    //mux.Use(NoSurf) 
    mux.Use(SessionLoad)

    mux.Get("/", Repo.Home)
    mux.Get("/model-3", Repo.Model3)
    mux.Get("/model-y", Repo.ModelY)

    mux.Get("/check-availability", Repo.CheckAvailability)
    mux.Post("/check-availability", Repo.PostAvailability)
    mux.Post("/check-availability-json", Repo.PostAvailabilityJSON)

    mux.Get("/rent", Repo.Rent)
    mux.Post("/rent", Repo.PostRent)
    mux.Get("/rent-summary", Repo.RentSummary)

    mux.Get("/about", Repo.About)
    mux.Get("/contact", Repo.Contact)

    // In static folder are all things that are not html template such as JS,
    // figures
    filesServer := http.FileServer(http.Dir("./static/"))
    mux.Handle("/static/*", http.StripPrefix("/static", filesServer))

    return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
    csrfHandler := nosurf.New(next)

    csrfHandler.SetBaseCookie(http.Cookie{
        HttpOnly: true,
        Path: "/",
        Secure: app.InProduction,
        SameSite: http.SameSiteLaxMode,
    })

    return csrfHandler
}

// SessionLoad loads and saves session on every request
func SessionLoad(next http.Handler) http.Handler {
    return session.LoadAndSave(next)
}

//CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
    // Create empty cache map of pointers to templates
    cache := make(map[string]*template.Template)
   
    // Get all files ending with *.page.html from ./templates/
    pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
    if err != nil {
        return cache, err
    }

    // Iterate through each page template file ending with *.page.html
    for _, page := range pages {
        name := filepath.Base(page)
        // Parse the page template file
        ts, err := template.New(name).Funcs(functions).ParseFiles(page)
        if err != nil {
            return cache, err
        }

        // Get all files ending with *.layout.html
        layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
        if err != nil {
            return cache, err
        }
        
        // If layout files exist, parse and add them to the template set
        if len(layouts) > 0 {
            ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
            if err != nil {
                return cache, err
            }
        }

        // Add the parsed template to the cache map
        cache[name] = ts
    }

    return cache, nil
}
