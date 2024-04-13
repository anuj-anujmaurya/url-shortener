package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

// home
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

//
type ShortURLRequest struct {
	CreateShortURL bool   `json:"CreateShortURL"`
	LongURL        string `json:"LongURL"`
}

// create url shortener
func (app *application) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.handleFetchError(w, "method not allowed", 1000, http.StatusMethodNotAllowed)
		return
	}

	var requestBody ShortURLRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		app.handleFetchError(w, "bad request/empty url", 1001, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(requestBody.LongURL) == "" {
		app.handleFetchError(w, "enter long url", 1002, http.StatusUnprocessableEntity)
		return
	}

	if utf8.RuneCountInString(requestBody.LongURL) > 400 {
		app.handleFetchError(w, "url too long", 1003, http.StatusUnprocessableEntity)
		return
	}

	// now check if the url is valid
	if !app.isValidUrl(requestBody.LongURL) {
		app.handleFetchError(w, "invalid url", 1004, http.StatusUnprocessableEntity)
		return
	}

	// check for existence in db
	shortURL, err := app.URLModel.CheckIfExists(requestBody.LongURL)

	if err != nil {
		app.handleFetchError(w, "something went wrong", 1005, http.StatusInternalServerError)
		return
	}

	if shortURL != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]any{"success": true, "short_url": shortURL})
		return
	}

	// create the short url
	id, err := app.URLModel.Insert(requestBody.LongURL)

	if err != nil {
		app.handleFetchError(w, "something went wrong", 1006, http.StatusInternalServerError)
		return
	}

	// convert the id into base 62 encoded
	shortURL = app.base62encode(id)

	// update the entry
	err = app.URLModel.UpdateShortURL(id, shortURL)
	if err != nil {
		app.handleFetchError(w, "something went wrong", 1007, http.StatusInternalServerError)
		return
	}

	// return the short url
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "short_url": shortURL})
}

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	longURL, err := app.URLModel.FindURL(shortURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
