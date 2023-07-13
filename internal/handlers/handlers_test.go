package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

// data for the Rent handler, /rent route 
var rentTests = []struct {
    name string
    rent models.Rent
    expectedStatusCode int
    expectedLocation string 
    expectedHTML string
}{
    {
        name: "rent in session",
        rent: models.Rent{
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        expectedStatusCode: http.StatusOK,
        expectedHTML: `action="/rent"`,
    },
    {
        name: "no rent in session",
        rent: models.Rent{},
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
        expectedHTML: "",
    },
    {
        name: "non existent model",
        rent: models.Rent{
            ModelID: 3,
            Model: models.Model{
                ID: 3,
                ModelName: "Model 3",
            },
        },
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
        expectedHTML: "",
    },
}

// TestRent tests the Rent handler  
func TestRent(t *testing.T) {
    for _, e := range rentTests {
        // create request
        r, _ := http.NewRequest("GET", "/rent", nil)
        // create context
        ctx := getCtx(r)
        // add context to request
        r = r.WithContext(ctx)

        // create recorder
        rr := httptest.NewRecorder()
        if e.rent.ModelID > 0 {
            session.Put(ctx, "rent", e.rent)
        }

        handler := http.HandlerFunc(Repo.Rent)
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedStatusCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
        }

        // test for location    
        if e.expectedLocation != "" {
            headers := rr.Result().Header
            if headers.Get("Location") != e.expectedLocation {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedLocation, headers.Get("Location"))
            }
        }

        // test for expected HTML
        if e.expectedHTML != "" {
            if !strings.Contains(rr.Body.String(), e.expectedHTML) {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedHTML, rr.Body.String())
            }
        }
    }
}

func TestRepository_PostRent(t *testing.T) {
    // dummy rent struct for testing
    rent := models.Rent{
        FirstName: "John",
        LastName: "Doe",
        Email: "john@doe.com",
        Phone: "+38599534256",
        ModelID: 1,
        Model: models.Model{
            ID: 1,
            ModelName: "Model 3",
        },
    }

    // post data necessary to create request
    reqBody := "start_date=2050-01-01"
    reqBody += "&end_date=2050-01-02"
    reqBody += "&first_name=John"
    reqBody += "&last_name=Doe"
    reqBody += "&email=john@doe.com"
    reqBody += "&phone=+38599534256"
    reqBody += "&model_id=1"

    // test case when there is post body data
    r, _ := http.NewRequest("POST", "/rent", strings.NewReader(reqBody))
    ctx := getCtx(r)
    r = r.WithContext(ctx)
    // set the header for the request (not necessary for this test but it is 
    // good practice). It is information to the server about the request type.
    // In this case it says that it is form post request.
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Repo.PostRent)
    session.Put(ctx, "rent", rent)
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusSeeOther {
        t.Errorf("PostRent handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
    }

    // test case when there is no rent in session
    r, _ = http.NewRequest("POST", "/rent", strings.NewReader(reqBody))
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("Rent handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
    }

    // test case when there is no post body data
    r, _ = http.NewRequest("POST", "/rent", nil)
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    // set the header for the request (not necessary for this test but it is 
    // good practice). It is information to the server about the request type.
    // In this case it says that it is form post request.
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr = httptest.NewRecorder()
    handler = http.HandlerFunc(Repo.PostRent)
    session.Put(ctx, "rent", rent)
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("PostRent handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
    }

    // test case when the form is not valid
    reqBody = "start_date=invalid"
    reqBody += "&end_date=invalid"
    reqBody += "&first_name=John"
    reqBody += "&last_name=Doe"
    reqBody += "&phone=+38599534256"
    reqBody += "&model_id=1"

    r, _ = http.NewRequest("POST", "/rent", strings.NewReader(reqBody))
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr = httptest.NewRecorder()
    handler = http.HandlerFunc(Repo.PostRent)
    session.Put(ctx, "rent", rent)
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusSeeOther {
        t.Errorf("PostRent handler returned wrong response code for invalid form: got %d, wanted %d", rr.Code, http.StatusSeeOther)
    }

    // test case for inserting rent into database using test-db
    // dummy rent struct for testing
    rent = models.Rent{
        FirstName: "John",
        LastName: "Doe",
        Email: "john@doe.com",
        Phone: "+38599534256",
        ModelID: 3,
        Model: models.Model{
            ID: 3,
            ModelName: "Model 3",
        },
    }

    // post data necessary to create request
    reqBody = "start_date=2050-01-01"
    reqBody += "&end_date=2050-01-02"
    reqBody += "&first_name=John"
    reqBody += "&last_name=Doe"
    reqBody += "&email=john@doe.com"
    reqBody += "&phone=+38599534256"
    reqBody += "&model_id=3"

    r, _ = http.NewRequest("POST", "/rent", strings.NewReader(reqBody))
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr = httptest.NewRecorder()
    handler = http.HandlerFunc(Repo.PostRent)
    // uddate session with new rent struct
    session.Put(ctx, "rent", rent)
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusSeeOther {
        t.Errorf("PostRent handler returned wrong response code for inserting rent into database: got %d, wanted %d", rr.Code, http.StatusSeeOther)
    }

    // test case for inserting rent restriction into database using test-db
    // dummy rent struct for testing
    rent = models.Rent{
        FirstName: "John",
        LastName: "Doe",
        Email: "john@doe.com",
        Phone: "+38599534256",
        ModelID: 4,
        Model: models.Model{
            ID: 4,
            ModelName: "Model 3",
        },
    }

    // post data necessary to create request
    reqBody = "start_date=2050-01-01"
    reqBody += "&end_date=2050-01-02"
    reqBody += "&first_name=John"
    reqBody += "&last_name=Doe"
    reqBody += "&email=john@doe.com"
    reqBody += "&phone=+38599534256"
    reqBody += "&model_id=4"

    r, _ = http.NewRequest("POST", "/rent", strings.NewReader(reqBody))
    ctx = getCtx(r)
    r = r.WithContext(ctx)
    r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    rr = httptest.NewRecorder()
    handler = http.HandlerFunc(Repo.PostRent)
    // uddate session with new rent struct
    session.Put(ctx, "rent", rent)
    rent = session.Get(ctx, "rent").(models.Rent)
    handler.ServeHTTP(rr, r)
    if rr.Code != http.StatusSeeOther {
        t.Errorf("PostRent handler returned wrong response code for inserting rent restriction into database: got %d, wanted %d", rr.Code, http.StatusSeeOther)
    }
}

// getCtx is a helper function that returns a request with a context
func getCtx(r *http.Request) context.Context {
    ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
    if err != nil {
        log.Println(err)
    }

    return ctx
}
































