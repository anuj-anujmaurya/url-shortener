package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// file server to server files from ui/static/ directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// href should look like '/static/css/main.css' for including css file
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/create", shortenHandler)
	mux.HandleFunc("/{shortURL}", redirectHandler)

	log.Println("Starting server on 8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
