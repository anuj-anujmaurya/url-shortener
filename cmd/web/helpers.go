package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
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

func (app *application) base62encode(number int64) string {
	base62 := ""
	for number > 0 {
		remainder := number % 62
		base62 = string(base62Digits[remainder]) + base62
		number /= 62
	}
	return base62
}

// decoding function
func (app *application) base62decode(hash string) int64 {
	var res int64
	for i := len(hash) - 1; i >= 0; i-- {
		charValue := int64(strings.Index(base62Digits, string(hash[i])))
		res += charValue * int64(math.Pow(62, float64(len(hash)-1-i)))
	}
	return res
}

// function to validate the url
func (app *application) isValidUrl(inputURL string) bool {
	pattern := `^(http(s)?://)?([\w-]+\.)+[\w-]+(/[\w- ;,./?%&=]*)?$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(inputURL)
}
