package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RomanLevBy/UrlShortener/internal/config"
	"github.com/RomanLevBy/UrlShortener/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(conf *config.Config) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			conf.Postgres.Host,
			conf.Postgres.User,
			conf.Postgres.DBName,
			conf.Postgres.Password,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS urls(
            id SERIAL PRIMARY KEY,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL
        );
    `)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	stmt, err = db.Prepare(`
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
	const fn = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls (url, alias) VALUES ($1, $2);")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgres.SaveURLGetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = $1")
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
	const fn = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s, %w", fn, err)
	}

	return nil
}
