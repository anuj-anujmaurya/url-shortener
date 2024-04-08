package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode/utf8"
)

// home
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is home handler")

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

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func convertToBase62(number int64) string {
	base62 := ""
	for number > 0 {
		remainder := number % 62
		base62 = string(base62Digits[remainder]) + base62
		number /= 62
	}
	return base62
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

	// check for existence in db
	shortURL, err := app.URLModel.CheckIfExists(requestBody.LongURL)

	if err != nil {
		app.handleFetchError(w, "Something went wrong", 1004, http.StatusInternalServerError)
	}

	if shortURL != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]any{"success": true, "shortURL": shortURL})
		return
	}

	// create the short url
	id, err := app.URLModel.Insert(requestBody.LongURL)

	if err != nil {
		app.handleFetchError(w, "Something went wrong", 1005, http.StatusInternalServerError)
		return
	}

	// convert the id into base 62 encoded
	shortURL = convertToBase62(id)

	// update the entry
	err = app.URLModel.UpdateShortURL(id, shortURL)
	if err != nil {
		app.handleFetchError(w, "Something went wrong", 1006, http.StatusInternalServerError)
	}

	// return the short url
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "short_url": shortURL})
}

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is redirect handler")
	shortURL := r.URL.Path
	fmt.Println(shortURL)
	// http.Redirect(w, r, shortURL, http.StatusMovedPermanently)
}
