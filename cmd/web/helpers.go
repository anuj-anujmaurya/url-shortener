package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// handle fetch error (this function isn't working as expected)
func (app *application) handleFetchError(w http.ResponseWriter, message string, code, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]any{"error": message, "code": code}
	json.NewEncoder(w).Encode(response)
}

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (app *application) convertToBase62(number int64) string {
	base62 := ""
	for number > 0 {
		remainder := number % 62
		base62 = string(base62Digits[remainder]) + base62
		number /= 62
	}
	return base62
}

// function to validate the url
func (app *application) isValidUrl(inputURL string) bool {
	pattern := `^(http(s)?://)?([\w-]+\.)+[\w-]+(/[\w- ;,./?%&=]*)?$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(inputURL)
}
