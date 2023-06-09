package render

import (
	"net/http"
//	"net/http/httptest"
	"testing"

	"github.com/sanijo/rent-app/internal/models"
)


func TestAddDefaultData(t *testing.T) {
    var td models.TemplateData

    r, err := getSession()
    if err != nil {
        t.Error(err)
    }

    session.Put(r.Context(), "flash", "123")

    result := AddDefaultData(&td, r)
    if result.Flash != "123" {
        t.Error("flash value of 123 not found in session")
    }
}

func TestRenderTemplate(t *testing.T) {
    pathToTemplates = "./../../templates"
    tc, err := CreateTemplateCache()
    if err != nil {
        t.Error(err)
    }

    r, err := getSession()
    if err != nil {
        t.Error(err)
    }

    app.TemplateCache = tc
//    ww := httptest.NewRecorder() // this can be used instead of dummyWriter
    ww := dummyWriter {}
    
	// Case 1: Existing template
	err = RenderTemplate(&ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Case 2: Non-existent template
	err = RenderTemplate(&ww, r, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func getSession() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/some-url", nil)
    if err != nil {
        return nil, err
    }
    
    // to avoid: panic: scs: no session data in context
    ctx := r.Context()
    ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
    r = r.WithContext(ctx)

    return r, nil
}

func TestNewTemplates(t *testing.T) {
    NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
    pathToTemplates = "./../../templates"

    _, err := CreateTemplateCache()
    if err != nil {
        t.Error(err)
    }
}
