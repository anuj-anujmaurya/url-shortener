package models

import (
	"database/sql"
	"time"
)

// define url type to hold the individual row data
type Url struct {
	ID       int
	LongURL  string
	ShortURL string
	Created  time.Time
}

// wrap sql.DB connection pool
type URLModel struct {
	DB *sql.DB
}

// check if long url exist in db, return short
func (m *URLModel) CheckIfExists(longURL string) (shortURL string, err error) {
	stmt := `SELECT short_url FROM url_map WHERE long_url = ?`
	row := m.DB.QueryRow(stmt, longURL)

	if err != nil {
		return "", err
	}

	var short_url string

	err = row.Scan(&short_url)
	if err != nil {
		// no record found
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err

	}

	return short_url, nil
}

// insert the long url, return id
func (m *URLModel) Insert(longURL string) (insertID int64, err error) {
	stmt := `INSERT INTO url_map (long_url, created) VALUES (?, UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, longURL)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

// update with id,
func (m *URLModel) UpdateShortURL(insertID int64, shortURL string) error {
	stmt := `UPDATE url_map SET short_url = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, shortURL, insertID)

	if err != nil {
		return err
	}

	return nil
}
