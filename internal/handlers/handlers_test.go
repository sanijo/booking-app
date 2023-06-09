package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
    key string
    value string
}

var theTests = []struct {
    name string
    url string
    method string
    params []postData
    expetedStatusCode int
}{
    {"home", "/", "GET", []postData{}, http.StatusOK},
    {"about", "/about", "GET", []postData{}, http.StatusOK},
    {"contact", "/contact", "GET", []postData{}, http.StatusOK},
    {"model-3", "/model-3", "GET", []postData{}, http.StatusOK},
    {"model-y", "/model-y", "GET", []postData{}, http.StatusOK},
    {"check-availability", "/check-availability", "GET", []postData{}, http.StatusOK},
    {"rent", "/rent", "GET", []postData{}, http.StatusOK},
    {"rent-summary", "/rent-summary", "GET", []postData{}, http.StatusOK},
    {"post-check-availability", "/check-availability", "POST", []postData{
        {key: "start", value: "2020-01-01"},
        {key: "end", value: "2020-01-02"},
    }, http.StatusOK},
    {"post-check-availability-json", "/check-availability-json", "POST", []postData{
        {key: "start", value: "2020-01-01"},
        {key: "end", value: "2020-01-02"},
    }, http.StatusOK},
    {"post-rent", "/rent", "POST", []postData{
        {key: "first_name", value: "John"},
        {key: "last_name", value: "Doe"},
        {key: "email", value: "name@lastname.com"},
        {key: "phone", value: "+38599534256"},
    }, http.StatusOK},
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
        } else { 
            values := url.Values{}
            for _, v := range e.params {
                values.Add(v.key, v.value)
            }

            response, err := ts.Client().PostForm(ts.URL+e.url, values)
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
