package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/handlers"
	"github.com/sanijo/rent-app/internal/helpers"
	"github.com/sanijo/rent-app/internal/models"
	"github.com/sanijo/rent-app/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main app function
func main() {

    err := run()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Starting application on port", portNumber)

    server := &http.Server {
        Addr: portNumber,
        Handler: routes(&app),
    }

    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Server error:", err)
    }
        
} 

func run() error {
    // What to put in session
    gob.Register(models.Rent{})

    // Change to true if in production
    app.InProduction = false

    // Set up info and error loggers
    infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    app.InfoLog = infoLog
    errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    app.ErrorLog = errorLog

    // Set up session
    session = scs.New()
    session.Lifetime = 24 * time.Hour
    session.Cookie.Persist = true
    session.Cookie.SameSite = http.SameSiteLaxMode
    session.Cookie.Secure = app.InProduction // in production: true

    // Set pointer in config to session so that is available in program
    app.Session = session

    // Create template cache
    tc, err := render.CreateTemplateCache()
	if err != nil {
        log.Fatal("cannot create template cache")
        return err
	}
   
    // Setting TemplateCache in config so that is cached all the time while app
    // is running
    app.TemplateCache = tc
    app.UseCache = false

    // Create repo and handlers
    repo := handlers.NewRepo(&app)
    handlers.NewHandlers(repo)
    // Give access to app config variable inside helpers package
    helpers.NewHelpers(&app)
    // Give access to app config variable inside render package
    render.NewTemplates(&app)

    return nil
}
