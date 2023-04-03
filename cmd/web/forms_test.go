package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_application_formHas(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has filed when it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "b")

	form = NewForm(postedData)

	has = form.Has("a")

	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func Test_application_formRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "something")
	postedData.Add("b", "something")
	postedData.Add("c", "something")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("shows post does not have required fields when it does")
	}
}

func Test_application_formCheck(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")

	if form.Valid() {
		t.Error("Valid() returns true and it should be false after calling Check()")
	}
}

func Test_application_formErrorGet(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")

	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from Get but do not")
	}

	s = form.Errors.Get("whatever")
	if len(s) != 0 {
		t.Error("should not have an error but got one")
	}
}
