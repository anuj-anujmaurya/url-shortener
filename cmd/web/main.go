package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"url_shortener/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// application struct to hold the app's dependencies for web app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	URLModel *models.URLModel
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

	username := "user"
	password := "*****"
	dbName := "url_shortener"
	dbHost := "localhost"
	dbPort := "3306"

	// Create a DSN (Data Source Name) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errorLog.Println("Error opening database connection:", err)
		return
	}

	defer db.Close()

	// initialize new instance of app struct (containing the dependencies)
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		URLModel: &models.URLModel{DB: db},
	}

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on 8080")

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
