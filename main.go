package main

import (
	"log"
	"net/http"
)

// home
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("creating url shortener"))
}

// create url shortener
func createShortUrl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("will create short url"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	mux.HandleFunc("/create/", createShortUrl)

	log.Println("Starting server on 8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
