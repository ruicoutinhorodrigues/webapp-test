package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html/"))))
	mux.Route("/web", func(mux chi.Router) {
		mux.Post("/auth", app.authenticate)
		// refresh-token
		mux.Get("/refresh-token", app.refreshUsingCookie)
		// logout
		mux.Get("/logout", app.deleteRefreshCookie)
	})

	// authentication routes auth handler, refresh
	mux.Post("/auth", app.authenticate)
	mux.Post("/refresh-token", app.refresh)

	// protected routes
	mux.Route("/users", func(mux chi.Router) {
		// use auth middleware
		mux.Use(app.authRequired)

		mux.Get("/", app.allUsers)
		mux.Get("/{userID}", app.getUser)
		mux.Delete("/{userID}", app.deleteUser)
		mux.Put("/", app.insertUser)
		mux.Patch("/", app.updateUser)
	})

	// test handler
	// mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	var payload = struct {
	// 		Message string `json:"message"`
	// 	}{
	// 		Message: "hello, world",
	// 	}

	// 	app.writeJSON(w, http.StatusOK, payload)
	// })

	return mux
}
