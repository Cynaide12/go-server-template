package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url_shortener/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}


func (s *Storage) SaveUrl(urlToSave string, alias string) error {
	const fn = "storage.sqlite.SaveUrl"

	tx, err := s.db.Begin()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}
	
	var existingURL string
	err = tx.QueryRow("SELECT url FROM url WHERE alias=?", alias).Scan(&existingURL)

	if err == nil{
		tx.Rollback()
		return storage.ErrAliasExists
	}

	if err != sql.ErrNoRows {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}


	stmt, err := tx.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}


	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}


	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	defer stmt.Close()


	return nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const fn = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}

	return nil

}
