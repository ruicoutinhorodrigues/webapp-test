package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"local/rui.rodrigues/webapp-test/pkg/data"

	"github.com/go-chi/chi/v5"
)

func Test_application_authenticate(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"not json", "not json body", http.StatusUnauthorized},
		{"empty json", "{}", http.StatusUnauthorized},
		{"empty email", `{"email":"","password":"secret"}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@example.com","password":""}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"admin@domain.com","password":"secret"}`, http.StatusUnauthorized},
	}

	for _, e := range tests {

		reader := strings.NewReader(e.requestBody)

		req, _ := http.NewRequest("POST", "/auth", reader)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_application_refresh(t *testing.T) {
	tests := []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid", "", http.StatusOK, true},
		{"expired token", expiredToken, http.StatusBadRequest, false},
		{"valid and not yet expired", "", http.StatusTooEarly, false},
	}

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	oldRefreshTime := refreshTokenExpiry

	for _, e := range tests {

		var tkn string

		if e.token == "" {

			if e.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}

			tokens, _ := app.generateTokenPair(&testUser)

			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		postedData := url.Values{
			"refresh_token": {tkn},
		}

		req, _ := http.NewRequest("POST", "/refresh-token", strings.NewReader(postedData.Encode()))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.refresh)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status of %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		refreshTokenExpiry = oldRefreshTime
	}
}

func Test_application_userHandlers(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{"allUsers", "GET", "", "", app.allUsers, http.StatusOK},
		{"deleteUser", "DELETE", "", "1", app.deleteUser, http.StatusNoContent},
		{"deleteUser bad URL param", "DELETE", "", "Y", app.deleteUser, http.StatusBadRequest},
		{"getUser valid", "GET", "", "1", app.getUser, http.StatusOK},
		{"getUser invalid", "GET", "", "2", app.getUser, http.StatusBadRequest},
		{"getUser bad URL param", "GET", "", "Y", app.getUser, http.StatusBadRequest},
		{
			"updateUser valid",
			"PATCH",
			`{"id":1,"first_name":"Administrator","last_name":"User","email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusNoContent,
		},
		{
			"updateUser invalid",
			"PATCH",
			`{"id":2,"first_name":"Administrator","last_name":"User","email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
		{
			"updateUser invalid json",
			"PATCH",
			`{"id":1,first_name:"Administrator","last_name":"User","email":"admin@example.com"}`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
		{
			"insertUser valid",
			"PUT",
			`{"id":1,"first_name":"JAck","last_name":"Smith","email":"smith@example.com"}`,
			"",
			app.insertUser,
			http.StatusNoContent,
		},
		{
			"insertUser invalid",
			"PUT",
			`{"foo":"bar","id":1,"first_name":"JAck","last_name":"Smith","email":"smith@example.com"}`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
		{
			"insertUser invalid json",
			"PUT",
			`{"id":1,first_name:"JAck","last_name":"Smith","email":"smith@example.com"}`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
	}

	for _, e := range tests {

		var req *http.Request

		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}

		if e.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(e.handler)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_application_refreshUsingCookie(t *testing.T) {
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPair(&testUser)

	testCookie := &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}

	badCookie := &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    "some_bad_string",
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	}

	tests := []struct {
		name           string
		addCookie      bool
		cookie         *http.Cookie
		expectedStatus int
	}{
		{"valid cookie", true, testCookie, http.StatusOK},
		{"invalid cookie", true, badCookie, http.StatusBadRequest},
		{"no cookie", false, nil, http.StatusUnauthorized},
	}

	for _, e := range tests {

		rr := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)

		if e.addCookie {
			req.AddCookie(e.cookie)
		}

		handler := http.HandlerFunc(app.refreshUsingCookie)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status code returned; expected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}

func Test_application_deleteRefreshCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.deleteRefreshCookie)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("wrong status; expected %d but got %d", http.StatusAccepted, rr.Code)
	}

	foundCookie := false

	for _, c := range rr.Result().Cookies() {
		if c.Name == "__Host-refresh_token" {

			foundCookie = true

			if c.Expires.After(time.Now()) {
				t.Errorf("cookie expiration in the future and should not be: %v", c.Expires.UTC())
			}
		}
	}

	if !foundCookie {
		t.Error("__Host-refresh_token cookie not found")
	}
}
