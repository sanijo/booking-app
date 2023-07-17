package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/sanijo/rent-app/internal/driver"
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

// data for the PostAvaialability handler 
var postAvailabilityTests = []struct {
    name string
    postedData url.Values
    expectedStatusCode int
    expectedLocation string
}{
    {
        name: "cannot parse form",
        postedData: nil,
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
    {
        name: "invalid start date",
        postedData: url.Values{
            "start": {"invalid"},
            "end": {"2021-05-20"},
        },
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
    {
        name: "invalid end date",
        postedData: url.Values{
            "start": {"2021-05-20"},
            "end": {"invalid"},
        },
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
    {
        name: "SearchAvailabilityForAllModels fails (start=2021-01-01)",
        postedData: url.Values{
            "start": {"2021-01-01"},
            "end": {"2021-01-02"},
        },
        expectedStatusCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
    {
        name: "length of models returned is 0",
        postedData: url.Values{
            "start": {"2021-05-20"},
            "end": {"2021-05-21"},
        },
        expectedStatusCode: http.StatusSeeOther,
        expectedLocation: "/check-availability",
    },
    {
        name: "models are available (start=2022-01-02)",
        postedData: url.Values{
            "start": {"2022-01-02"},
            "end": {"2022-01-03"},
        },
        expectedStatusCode: http.StatusOK,
        expectedLocation: "",
    },
}

// TestPostAvailability tests the PostAvailability handler /check-availability route
func TestPostAvailability(t *testing.T) {
    for _, e := range postAvailabilityTests {
        var r *http.Request
        if e.postedData != nil {
            r, _ = http.NewRequest("POST", "/check-availability", strings.NewReader(e.postedData.Encode()))
        } else {
            r, _ = http.NewRequest("POST", "/check-availability", nil)
        }
        ctx := getCtx(r)
        r = r.WithContext(ctx)
        r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        rr := httptest.NewRecorder()

        // create and call handler
        handler := http.HandlerFunc(Repo.PostAvailability)
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedStatusCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
        }
    }
}

var postAvailabilityJSONTests = []struct {
    name string
    postedData url.Values
    expectedOK bool
}{
    {
        name: "no available models (start=2021-01-01)",
        postedData: url.Values{
            "start": {"2021-01-01"},
            "end": {"2021-01-02"},
            "model_id": {"1"},
        },
        expectedOK: false,
    },
    {
        name: "models are available",
        postedData: url.Values{
            "start": {"2022-01-02"},
            "end": {"2022-01-03"},
            "model_id": {"1"},
        },
        expectedOK: true,
    },
    {
        name: "empty post body",
        postedData: nil,
        expectedOK: false,
    },
    {
        name: "database query returns error (start=2021-01-02)",
        postedData: url.Values{
            "start": {"2021-01-02"},
            "end": {"2021-01-03"},
            "model_id": {"1"},
        },
        expectedOK: false,
    },
}

// PostAvailabilityJSON tests the PostAvailabilityJSON handler /check-availability-json route
func TestPostAvailabilityJSON(t *testing.T) {
    for _, e := range postAvailabilityJSONTests {
        var r *http.Request
        if e.postedData != nil {
            r, _ = http.NewRequest("POST", "/check-availability-json", strings.NewReader(e.postedData.Encode()))
        } else {
            r, _ = http.NewRequest("POST", "/check-availability-json", nil)
        }
        ctx := getCtx(r)
        r = r.WithContext(ctx) // add context to request
        r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

        // make handler HandlerFunc
        handler := http.HandlerFunc(Repo.PostAvailabilityJSON)
        // create response recorder
        rr := httptest.NewRecorder()
        // make request to handler
        handler.ServeHTTP(rr, r)

        // test for json response 
        var response jsonResponse
        // get what is being handed back from the server and put it into the response struct
        err := json.Unmarshal([]byte(rr.Body.String()), &response)
        if err != nil {
            t.Errorf("error parsing json")
        }

        if response.OK != e.expectedOK {
            t.Errorf("for %s, expected %v but got %v", e.name, e.expectedOK, response.OK)
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

// postRentTests is data for the PostRent handler
var  postRentTests = []struct {
    name string
    inSession bool
    rent models.Rent
    postedData url.Values
    expectedResponseCode int
    expectedLocation string
    expectedHTML string
}{
    {
        name: "valid post data",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        postedData: url.Values{
            "start_date": {"2050-01-01"},
            "end_date": {"2050-01-02"}, 
            "first_name": {"John"},
            "last_name": {"Doe"},
            "email": {"john@doe.com"},
            "phone": {"+38599534256"},
            "model_id": {"1"},
        },
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "/rent-summary",
        expectedHTML: "",
    },
    {
        name: "invalid post data",
        inSession: true,
        rent: models.Rent{},
        postedData: nil,
        expectedResponseCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
        expectedHTML: "",
    },
    {
        name: "no rent in session",
        inSession: false,
        rent: models.Rent{},
        postedData: url.Values{
            "start_date": {"2050-01-01"},
            "end_date": {"2050-01-02"},
            "first_name": {"John"},
            "last_name": {"Doe"},
            "email": {"john@doe.com"},
            "phone": {"+38599534256"},
            "model_id": {"1"},
        },
        expectedResponseCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
        expectedHTML: "",
    },
    {
        name: "invalid form data (missing required fields)",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        postedData: url.Values{
            "start_date": {"2050-01-01"},
            "end_date": {"2050-01-02"},
            "first_name": {"John"},
            "last_name": {"Doe"},
            "phone": {"+38599534256"},
            "model_id": {"1"},
        },
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "",
        expectedHTML: "",
    },
    {
        name: "insert rent into database fails (ModelID == 3)",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 3,
            Model: models.Model{
                ID: 3,
                ModelName: "Model 3",
            },
        },
        postedData: url.Values{
            "start_date": {"2050-01-01"},
            "end_date": {"2050-01-02"},
            "first_name": {"John"},
            "last_name": {"Doe"},
            "email": {"john@doe.com"},
            "phone": {"+38599534256"},
            "model_id": {"3"},
        },
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "/",
    },
    {
        name: "insert rent restriction into database fails (ModelID == 4)",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 4,
            Model: models.Model{
                ID: 4,
                ModelName: "Model 3",
            },
        },
        postedData: url.Values{
            "start_date": {"2050-01-01"},
            "end_date": {"2050-01-02"},
            "first_name": {"John"},
            "last_name": {"Doe"},
            "email": {"john@doe.com"},
            "phone": {"+38599534256"},
            "model_id": {"4"},
        },
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "/",
    },
}

func TestPostRent(t *testing.T) {
    for _, e := range postRentTests {
        var r *http.Request
        if e.postedData != nil {
            r, _ = http.NewRequest("POST", "/rent", strings.NewReader(e.postedData.Encode()))
        } else {
            r, _ = http.NewRequest("POST", "/rent", nil)
        }
        // create context
        ctx := getCtx(r)
        // add context to request
        r = r.WithContext(ctx)
        // set content type
        r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        // create recorder 
        rr := httptest.NewRecorder()
        // create handler 
        handler := http.HandlerFunc(Repo.PostRent)
        // check if rent is in a session, and if so put it in a session 
        if e.inSession {
            session.Put(ctx, "rent", e.rent)
        }
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedResponseCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedResponseCode, rr.Code)
        }

        // test for Location
        if e.expectedLocation != "" {
            headers := rr.Result().Header
            if headers.Get("Location") != e.expectedLocation {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedLocation, headers.Get("Location"))
            }
        }

        // test for expected expected HTML 
        if e.expectedHTML != "" {
            if !strings.Contains(rr.Body.String(), e.expectedHTML) {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedHTML, rr.Body.String())
            }
        }
    }

}

// rentSummaryTests is a struct that holds test data for the RentSummary handler
var rentSummaryTests = []struct {
    name string
    inSession bool
    rent models.Rent
    expectedResponseCode int
    expectedLocation string
}{
    {
        name: "rent in session",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        expectedResponseCode: http.StatusOK,
        expectedLocation: "",
    },
    {
        name: "rent not in session",
        inSession: false,
        rent: models.Rent{},
        expectedResponseCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
}

//TestRentSummary tests the RentSummary handler
func TestRentSummary(t *testing.T) {
    for _, e := range rentSummaryTests {
        var r *http.Request
        r, _ = http.NewRequest("GET", "/rent-summary", nil)
        // create context
        ctx := getCtx(r)
        // add context to request
        r = r.WithContext(ctx)
        // create recorder 
        rr := httptest.NewRecorder()
        // create handler 
        handler := http.HandlerFunc(Repo.RentSummary)
        // check if rent is in a session, and if so put it in a session 
        if e.inSession {
            session.Put(ctx, "rent", e.rent)
        }
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedResponseCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedResponseCode, rr.Code)
        }

        // test for Location
        if e.expectedLocation != "" {
            headers := rr.Result().Header
            if headers.Get("Location") != e.expectedLocation {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedLocation, headers.Get("Location"))
            }
        }
    }
}

// testChooseModel is a struct that holds tests for the ChooseModel handler
var testChooseModel = []struct {
    name string
    inSession bool
    rent models.Rent
    url string
    expectedResponseCode int
    expectedLocation string
}{
    {
        name: "rent in session",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        url: "/choose-model/1",
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "/rent",
    },
    {
        name: "rent not in session",
        inSession: false,
        rent: models.Rent{},
        url: "/choose-model/1",
        expectedResponseCode: http.StatusSeeOther,
        expectedLocation: "/",
    },
    {
        name: "invalid model id (URL parameter)",
        inSession: true,
        rent: models.Rent{
            FirstName: "John",
            LastName: "Doe",
            Email: "john@doe.com",
            Phone: "+38599534256",
            ModelID: 1,
            Model: models.Model{
                ID: 1,
                ModelName: "Model 3",
            },
        },
        url: "/choose-model/invalid",
        expectedResponseCode: http.StatusTemporaryRedirect,
        expectedLocation: "/",
    },
}

// TestChooseModel tests the ChooseModel handler
func TestChooseModel(t *testing.T) {
    for _, e := range testChooseModel {
        var r *http.Request
        r, _ = http.NewRequest("GET", e.url, nil)
        // create context
        ctx := getCtx(r)
        // add context to request
        r = r.WithContext(ctx)
        // create recorder 
        rr := httptest.NewRecorder()
        // create handler 
        handler := http.HandlerFunc(Repo.ChooseModel)
        // check if rent is in a session, and if so put it in a session 
        if e.inSession {
            session.Put(ctx, "rent", e.rent)
        }
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedResponseCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedResponseCode, rr.Code)
        }

        // test for Location
        if e.expectedLocation != "" {
            headers := rr.Result().Header
            if headers.Get("Location") != e.expectedLocation {
                t.Errorf("for %s, expected %s but got %s", e.name, e.expectedLocation, headers.Get("Location"))
            }
        }
    }
}

// testRentVehicle is a struct that holds tests for the RentVehicle handler
var testRentVehicle = []struct {
    name string
    url string
    expectedResponseCode int
}{
    {
        name: "valid url parameters",
        url: "/rent-vehicle?s=2023-01-01&e=2023-01-02&id=1",
        expectedResponseCode: http.StatusSeeOther,
    },
    {
        name: "invalid url parameters",
        url: "/rent-vehicle?s=invalid&e=invalid&id=invalid",
        expectedResponseCode: http.StatusTemporaryRedirect,
    },
    {
        name: "invalid start date",
        url: "/rent-vehicle?s=invalid&e=2023-01-02&id=1",
        expectedResponseCode: http.StatusSeeOther,
    },
    {
        name: "invalid end date",
        url: "/rent-vehicle?s=2023-01-01&e=invalid&id=1",
        expectedResponseCode: http.StatusSeeOther,
    },
    {
        name: "invalid model id",
        url: "/rent-vehicle?s=2023-01-01&e=2023-01-02&id=4",
        expectedResponseCode: http.StatusSeeOther,
    },
}


// TestRentVehicle tests the RentVehicle handler
func TestRentVehicle(t *testing.T) {
    for _, e := range testRentVehicle {
        var r *http.Request
        r, _ = http.NewRequest("GET", e.url, nil)
        // create context
        ctx := getCtx(r)
        // add context to request
        r = r.WithContext(ctx)
        // create recorder 
        rr := httptest.NewRecorder()
        // create handler 
        handler := http.HandlerFunc(Repo.RentVehicle)
        handler.ServeHTTP(rr, r)

        // test for status code
        if rr.Code != e.expectedResponseCode {
            t.Errorf("for %s, expected %d but got %d", e.name, e.expectedResponseCode, rr.Code)
        }
    }
}


// TestNewRepo tests the NewRepo function
func TestNewRepo(t *testing.T) {
    var db driver.DB 
    testRepo := NewRepo(&app, &db)

    if reflect.TypeOf(testRepo) != reflect.TypeOf(&Repository{}) {
        t.Error("NewRepo did not return a pointer to testrepo")
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
































