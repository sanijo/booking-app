package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sanijo/rent-app/internal/config"
	"github.com/sanijo/rent-app/internal/models"
)


var session *scs.SessionManager
var testApp config.AppConfig


func TestMain(m *testing.M) {
    // What to put in session
    gob.Register(models.Rent{})

    // Change to true if in production
    testApp.InProduction = false

    // Set up info and error loggers
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    testApp.InfoLog = infoLog
    errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    testApp.ErrorLog = errorLog

    session = scs.New()
    session.Lifetime = 24 * time.Hour
    session.Cookie.Persist = true
    session.Cookie.SameSite = http.SameSiteLaxMode
    session.Cookie.Secure = false

    // Set pointer in config to session so that is available in program
    testApp.Session = session

    app = &testApp

    os.Exit(m.Run())
}

// Type that satisfies structure of the http.ResponseWriter
type dummyWriter struct {}

func (tw *dummyWriter) Header() http.Header {
    return http.Header {}
}

func (tw *dummyWriter) WriteHeader(i int) {}

func (tw *dummyWriter) Write(b []byte) (int, error) {
    //var length = 1 // CAUTION: this won't work
    length := len(b)
    return length, nil
}
