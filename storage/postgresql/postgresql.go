package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"telegram-go/lib/e"
	"telegram-go/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, e.WrapError("can't open database", err)
	}
	if err := db.Ping(); err != nil {
		return nil, e.WrapError("can't ping database", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	query := `INSERT INTO pages (url, username) VALUES ($1, $2)`
	if _, err := s.db.ExecContext(ctx, query, page.URL, page.UserName); err != nil {
		return err
	}
	return nil
}

func (s *Storage) ChoseRandom(ctx context.Context, username string) (*storage.Page, error) {
	query := `SELECT url FROM pages WHERE username = $1 ORDER BY RANDOM() LIMIT 1`
	var url string

	err := s.db.QueryRowContext(ctx, query, username).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't chose random page: %w", err)
	}
	return &storage.Page{URL: url, UserName: username}, nil
}

// Remove Removes page from database
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	query := `DELETE FROM pages WHERE url = $1 AND username = $2`
	if _, err := s.db.ExecContext(ctx, query, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}

// Exists Checks if random page exists in database
func (s *Storage) Exists(ctx context.Context, p *storage.Page) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pages WHERE url = $1 AND username = $2)`
	var result bool
	if err := s.db.QueryRowContext(ctx, query, p.URL, p.UserName).Scan(&result); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}
	return result == true, nil
}

// InitDatabase Initializes database
func (s *Storage) InitDatabase(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, username TEXT)`
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}
