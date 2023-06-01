package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/sanijo/booking-app/internal/config"
	"github.com/sanijo/booking-app/internal/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
    app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
    td.CSRFToken = nosurf.Token(r)
    return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

    tc := make(map[string]*template.Template)
    var err error

    if app.UseCache {
        // Get the template cache from the app config
        tc = app.TemplateCache 
    } else {
        tc, err = CreateTemplateCache()
	    if err != nil {
            fmt.Println("Cannot create template cache:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	    }
    }

    // Get requested template from cache
    t, available := tc[tmpl]
    if !available {
		fmt.Println("Template unavailable:", tmpl)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
    }

    // Render the template
    buffer := new(bytes.Buffer)
    
    // Template default data
    td = AddDefaultData(td, r)

    err = t.Execute(buffer, td)
    if err != nil {
        fmt.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    _, err = buffer.WriteTo(w)
    if err != nil {
        fmt.Println("Error writing template to response:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
    // Create empty cache map of pointers to templates
    cache := make(map[string]*template.Template)
   
    // Get all files ending with *.page.html from ./templates/
    pages, err := filepath.Glob("./templates/*.page.html")
    if err != nil {
        return cache, err
    }

    // Iterate through each page template file ending with *.page.html
    for _, page := range pages {
        name := filepath.Base(page)
        // Parse the page template file
        ts, err := template.New(name).ParseFiles(page)
        if err != nil {
            return cache, err
        }

        // Get all files ending with *.layout.html
        layouts, err := filepath.Glob("./templates/*.layout.html")
        if err != nil {
            return cache, err
        }
        
        // If layout files exist, parse and add them to the template set
        if len(layouts) > 0 {
            ts, err = ts.ParseGlob("./templates/*.layout.html")
            if err != nil {
                return cache, err
            }
        }

        // Add the parsed template to the cache map
        cache[name] = ts
    }

    return cache, nil
}
