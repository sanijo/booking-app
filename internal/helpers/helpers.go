package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sanijo/rent-app/internal/config"
)


var app *config.AppConfig

// NewHelpers sets up the app config for the helpers package
func NewHelpers(a *config.AppConfig) {
    app = a
}

// ClientError logs client errors
func CilentError(w http.ResponseWriter, status int) {
    app.InfoLog.Println("Client error. Status:", status)
    http.Error(w, http.StatusText(status), status)
}

// ServerError logs server errors
func ServerError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.ErrorLog.Println(trace)
    http.Error(
        w, 
        http.StatusText(http.StatusInternalServerError),
        http.StatusInternalServerError)
}
