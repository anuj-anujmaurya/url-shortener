package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode/utf8"
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

// create url shortener
func (app *application) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.handleFetchError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	fmt.Println("hello")
	if err != nil {
		app.handleFetchError(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	long_url := r.PostForm.Get("long_url")

	if strings.TrimSpace(long_url) == "" {
		app.handleFetchError(w, "This field can't be empty", http.StatusUnprocessableEntity)
		return
	} else if utf8.RuneCountInString(long_url) > 400 {
		app.handleFetchError(w, "URL too long", http.StatusUnprocessableEntity)
		return
	}
	fmt.Println("hello again")


	// check for existence in db

	// create the short url

	// return the short url

	w.Write([]byte("will create short url"))
}

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	// find the long url
	var longURL string

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
