package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	r := httptest.NewRequest("POST", "/url", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("Not Valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/url", nil)
	form := New(r.PostForm)
	form.Required("A", "B", "C")
	if form.Valid() {
		t.Error("Form Showing Valid When It Isn't")
	}
	postedData := url.Values{}
	postedData.Add("A", "A")
	postedData.Add("B", "A")
	postedData.Add("C", "A")

	r, _ = http.NewRequest("POST", "/url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("A", "B", "C")
	if !form.Valid() {
		t.Error("Form Showing Not Valid When It Is")
	}
}

func TestForm_Has(t *testing.T) {
	// r := httptest.NewRequest("POST", "/url", nil)
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("Whatever")
	if has {
		t.Error("Failure By Assuming That An Non-existent Field Exists")
	}

	postedData = url.Values{}
	postedData.Add("A", "ABC")
	form = New(postedData)

	has = form.Has("A")
	if !has {
		t.Error("Failure Aknowledging An Existent Field")
	}
}

func TestForm_MinLength(t *testing.T) {
	// r := httptest.NewRequest("POST", "/url", nil)
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("X", 10)
	if form.Valid() {
		t.Error("Form Shows MinLength for Non existent Field")
	}

	isError := form.Errors.Get("X")
	if isError == "" {
		t.Error("Should Have An Error But Didn't Get One")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")
	form = New(postedData)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("Form Shows MinLength of 100 Met when Data Is 10 Charactrs Long")
	}
	postedData = url.Values{}
	postedData.Add("another_field", "another_value")
	form = New(postedData)

	form.MinLength("another_field", 2)
	if !form.Valid() {
		t.Error("Form Shows MinLength of 2 Not Met when Data Is 13 Charactrs Long")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Shouldn't Have An Error But Did Get One")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// r := httptest.NewRequest("POST", "/url", nil)
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("X")
	if form.Valid() {
		t.Error("Form Shows Valid Email For Non Existent Field")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@here.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Form Shows Not Valid Email For Valid Email Field")
	}

	postedData = url.Values{}
	postedData.Add("email", "hjhdfssnsdl")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Form Shows Valid Email For Non-Valid Email Field")
	}
}
