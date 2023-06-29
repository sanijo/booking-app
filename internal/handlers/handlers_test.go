package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanijo/rent-app/internal/models"
)

type postData struct {
    key string
    value string
}

var theTests = []struct {
    name string
    url string
    method string
    expetedStatusCode int
}{
    {"home", "/", "GET", http.StatusOK},
    {"about", "/about", "GET", http.StatusOK},
    {"contact", "/contact", "GET", http.StatusOK},
    {"model-3", "/model-3", "GET", http.StatusOK},
    {"model-y", "/model-y", "GET", http.StatusOK},
    {"check-availability", "/check-availability", "GET", http.StatusOK},
    {"rent-summary", "/rent-summary", "GET", http.StatusOK},

// Those below use session and database, so they are not suitable for testing
//    {"post-check-availability", "/check-availability", "POST", []postData{
//        {key: "start", value: "2020-01-01"},
//        {key: "end", value: "2020-01-02"},
//    }, http.StatusOK},
//    {"post-check-availability-json", "/check-availability-json", "POST", []postData{
//        {key: "start", value: "2020-01-01"},
//        {key: "end", value: "2020-01-02"},
//    }, http.StatusOK},
//    {"post-rent", "/rent", "POST", []postData{
//        {key: "first_name", value: "John"},
//        {key: "last_name", value: "Doe"},
//        {key: "email", value: "name@lastname.com"},
//        {key: "phone", value: "+38599534256"},
//    }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
    routes := getRoutes()
    // create test server
    ts := httptest.NewTLSServer(routes)
    defer ts.Close()

    for _, e := range theTests {
        if e.method == "GET" {
            response, err := ts.Client().Get(ts.URL + e.url)
            if err != nil {
                t.Log(err)
                t.Fatal(err)
            }
            if response.StatusCode != e.expetedStatusCode {
                t.Errorf("for %s, expected %d but %d received",
                    e.name, e.expetedStatusCode, response.StatusCode)
            }
        }
    }
}

func TestRepository_Rent(t *testing.T) {
    // test case when there is rent in session
    rent := models.Rent{
        ModelID: 1,
        Model: models.Model{
            ID: 1,
            ModelName: "Model 3",
        },
    }

    r, _ := http.NewRequest("GET", "/rent", nil)
    ctx := getCtx(r)
    r = r.WithContext(ctx)

    rr := httptest.NewRecorder()
    session.Put(ctx, "rent", rent)
    handler := http.HandlerFunc(Repo.Rent)
    handler.ServeHTTP(rr, r)

    if rr.Code != http.StatusOK {
        t.Errorf("Rent handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
    }

    // test case when there is no rent in session
    r, _ = http.NewRequest("GET", "/rent", nil)
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, r)

    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("Rent handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
    }

    // test for case where there is no model 
    r, _ = http.NewRequest("GET", "/rent", nil)
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    rr = httptest.NewRecorder()
    rent.ModelID = 3
    session.Put(ctx, "rent", rent)
    handler.ServeHTTP(rr, r)

    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("Rent handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
    }
}

func TestRepository_PostRent(t *testing.T) {
}

func getCtx(r *http.Request) context.Context {
    ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
    if err != nil {
        log.Println(err)
    }

    return ctx
}
































