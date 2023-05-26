package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/sanijo/booking-app/pkg/config"
	"github.com/sanijo/booking-app/pkg/handlers/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
    app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
    return td
}

// V3
func RenderTemplate(
    w http.ResponseWriter, 
    tmpl string, 
    td *models.TemplateData) {

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
    td = AddDefaultData(td)

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

//// V1
//// RenderTemplate renders templates using html/template
//func RenderTemplateBase(w http.ResponseWriter, tmpl string) {
//	parsedTemplate, err := 
//        template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.html")
//	if err != nil {
//		fmt.Println("Error parsing template:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//
//	err = parsedTemplate.Execute(w, nil)
//	if err != nil {
//		fmt.Println("Error executing template:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//}

//// V2
//// Global map variable that stores template so that it doesnt need to be loaded
//// every time
//var tc = make(map[string]*template.Template)
//
//func RenderTemplate(w http.ResponseWriter, t string) {
//    var tmpl *template.Template
//
//    // Check to see if template allready exists in cache
//    _, inMap := tc[t]
//    if !inMap {
//        // Need to create the template (read from disk and parse it)
//        fmt.Println("Creating template and adding to cache")
//        err := createTemplateCache(t)
//	    if err != nil {
//	    	fmt.Println("Error parsing template:", err)
//	    	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	    	return
//	    }
//    } else {
//        // Template allready in the cache
//        fmt.Println("Using cached template")
//    }
//
//    tmpl = tc[t]
//
//    err := tmpl.Execute(w, nil)
//	if err != nil {
//		fmt.Println("Error executing template:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//}
//
//func createTemplateCache(t string) error  {
//    templates := [] string {
//        fmt.Sprintf("./templates/%s",t),
//        "./templates/base.layout.html",
//    }
//
//    // Parse the template
//    tmpl, err := template.ParseFiles(templates...)
//    if err != nil {
//        return err
//    }
//
//    // Add template to cache
//    tc[t] = tmpl
//
//    return nil
//}

