package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)

    isValid := form.Valid()
    if !isValid {
        t.Error("got invalid when valid is expected")
    }
}

func TestForm_Required(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)

    form.Required("a", "b", "c")
    if form.Valid() {
        t.Error("forms shows valid when required fields are missing")
    }

    postedData = url.Values{}
    postedData.Add("a", "a")
    postedData.Add("b", "b")
    postedData.Add("c", "c")

    form = New(postedData)

    form.Required("a", "b", "c")
    if !form.Valid() {
        t.Error("forms shows invalid when required fields are present")
    }
}

func TestForm_Has(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)

    has := form.Has("a")
    if has {
        t.Error("form shows valid when form field is empty")
    }

    postedData = url.Values{}
    postedData.Add("a", "a")
    form = New(postedData)
    has = form.Has("a")
    
    if !has {
        t.Error("form shows invalid when form field is not empty")
    }
}

func TestForm_MinLength(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)
    
    form.MinLength("field", 10)
    if form.Valid() {
        t.Error("form shows min length for non-existent field")
    }

    postedData = url.Values{}
    postedData.Add("field", "value")
    form = New(postedData)

    form.MinLength("field", 4)
    if !form.Valid() {
        t.Error("form shows invalid when field lenght is larger than min length")
    }

    isError := form.Errors.Get("field")
    if isError != "" {
        t.Error("retrieved error message for field that passed validation")
    }

    postedData = url.Values{}
    postedData.Add("field", "123")
    form = New(postedData)

    form.MinLength("field", 4)
    if form.Valid() {
        t.Error("form shows valid when field lenght is smaller than min length")
    }

    isError = form.Errors.Get("field")
    if isError == "" {
        t.Error("did not retrieve error message for field that failed validation")
    }

}

func TestForm_IsEmail(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)

    form.IsEmail("email")
    if form.Valid() {
        t.Error("form shows valid email for non-existent field")
    }

    postedData = url.Values{}
    postedData.Add("email", "name@gmail.com")
    form = New(postedData)

    form.IsEmail("email")
    if !form.Valid() {
        t.Error("got invalid for valid email address")
    }

    postedData = url.Values{}
    postedData.Add("email", "wrong")
    form = New(postedData)

    form.IsEmail("email")
    if form.Valid() {
        t.Error("got valid for invalid email address")
    }
}

func TestForm_IsURL(t *testing.T) {
    postedData := url.Values{}
    form := New(postedData)

    form.IsURL("url")
    if form.Valid() {
        t.Error("form shows valid url for non-existent field")
    }

    postedData = url.Values{}
    postedData.Add("url", "https://www.google.com")
    form = New(postedData)

    form.IsURL("url")
    if !form.Valid() {
        t.Error("got invalid for valid url")
    }

    postedData = url.Values{}
    postedData.Add("url", "wrong")
    form = New(postedData)

    form.IsURL("url")
    if form.Valid() {
        t.Error("got valid for invalid url")
    }
}
