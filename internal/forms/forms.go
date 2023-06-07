package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
    url.Values
    Errors errors
}

//New initializes empty form struct
func New(data url.Values) *Form {
    return &Form{
        data,
        errors(map[string][]string{}),
    }
}

// Valid returns true if there is no errors
func (f *Form) Valid() bool {
    return len(f.Errors) == 0
}

// Required checks if required fields are not empty
func (f *Form) Required(fields ...string) {
    for _, field := range fields {
        value := f.Get(field)
        if strings.TrimSpace(value) == "" {
            f.Errors.Add(field, "This field cannot be empty")
        }
    }
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
    x := r.Form.Get(field)
    if x == "" {
        return false
    }
    return true
}

// MinLength checks for string min length
func (f *Form) MinLength(field string, length int, r*http.Request) bool {
    x := r.Form.Get(field)
    if len(x) < length {
        f.Errors.Add(
            field, 
            fmt.Sprintf("This field must be at least %d characters long", length))
        return false
    }
    return true
}

// IsEmail checks if the input is valid email
func (f *Form) IsEmail(field string) {
    if !govalidator.IsEmail(f.Get(field)) {
        f.Errors.Add(field, "Invalid email adress")
    }
}
