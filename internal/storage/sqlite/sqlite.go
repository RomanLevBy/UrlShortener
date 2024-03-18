package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RomanLevBy/UrlShortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS urls(
        id INTEGER PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias)
    `)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const fn = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls (url, alias) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		//todo refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s, %w", fn, storage.ErrURLExists)
		}

		return fmt.Errorf("%s, %w", fn, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.sqlite.SaveURLGetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s, %w", fn, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s, %w", fn, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.sqlite.DeleteURLL"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	return nil
}
