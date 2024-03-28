package main

import (
	"log"
	"net/http"
	"os"
)

// application struct to hold the app's dependencies for web app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Open separate files for info and error logs
	infoLogFile, err := os.OpenFile("./logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer infoLogFile.Close()

	errorLogFile, err := os.OpenFile("./logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer errorLogFile.Close()

	// Create loggers
	infoLog := log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(errorLogFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize new instance of app struct (containing the dependencies)
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	// file server to server files from ui/static/ directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.shortenHandler)
	mux.HandleFunc("/{shortURL}", app.redirectHandler)

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Println("Starting server on 8080")

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
