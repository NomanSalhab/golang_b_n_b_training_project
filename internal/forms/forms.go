package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct and it embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid return true if there are no errors otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has Checks If Form Field Is In Post And Is Nopt Empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This Field Cannot Be Blank")
		return false
	}
	return true
}
