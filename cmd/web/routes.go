package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// Define file server for serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	router.Get("/", app.home)
	router.Post("/create", app.shortenHandler)
	router.Get("/create/custom", app.createCustomPage)
	router.Post("/create/custom", app.createCustomPost)
	router.Get("/{shortURL:[0-9a-zA-Z]+}", app.redirectHandler)

	return router
}
