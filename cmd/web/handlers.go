package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode/utf8"
	"url_shortener/internal/models"
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
		app.handleFetchError(w, "Method Not Allowed", 1000, http.StatusMethodNotAllowed)
		return
	}

	var requestBody ShortURLRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		app.handleFetchError(w, "Bad request/empty url", 1001, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(requestBody.LongURL) == "" {
		app.handleFetchError(w, "This field can't be empty", 1002, http.StatusUnprocessableEntity)
		return
	} else if utf8.RuneCountInString(requestBody.LongURL) > 400 {
		app.handleFetchError(w, "URL too long", 1003, http.StatusUnprocessableEntity)
		return
	}

	urlModel := models.URLModel{DB: &sql.DB{}}

	// check for existence in db
	shortURL, err := urlModel.CheckIfExists(requestBody.LongURL)

	if err != nil {
		app.handleFetchError(w, "Something went wrong", 1004, http.StatusInternalServerError)
	} else {
		fmt.Println(shortURL)
		return
	}

	if shortURL != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"shortURL": shortURL})
		return
	}

	// create the short url

	// return the short url

	w.Write([]byte("will create short url"))
}

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	// find the long url
	var longURL string

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
