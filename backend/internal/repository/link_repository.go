package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"go-links/internal/models"
)

const timeLayout = time.RFC3339Nano

var NoError = errors.New("link not found")

var ErrShortcut = errors.New("shortcut already exists")

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(shortcut, destination string) (*models.Link, error) {
	res, err := r.db.Exec(
		`INSERT INTO links (shortcut, destination, created_at) VALUES (?, ?, ?)`,
		shortcut, destination, time.Now().UTC().Format(timeLayout),
	)
	if err != nil {
		if isErr(err) {
			return nil, ErrShortcut
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *LinkRepository) GetByID(id int64) (*models.Link, error) {
	row := r.db.QueryRow(
		`SELECT id, shortcut, destination, created_at, click_count FROM links WHERE id = ?`,
		id,
	)
	return scanLink(row)
}

func (r *LinkRepository) GetByShortcut(shortcut string) (*models.Link, error) {
	row := r.db.QueryRow(
		`SELECT id, shortcut, destination, created_at, click_count FROM links WHERE shortcut = ?`,
		shortcut,
	)
	return scanLink(row)
}

func (r *LinkRepository) List() ([]*models.Link, error) {
	rows, err := r.db.Query(
		`SELECT id, shortcut, destination, created_at, click_count FROM links ORDER BY created_at DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []*models.Link{}
	for rows.Next() {
		var l models.Link
		var createdAt string
		if err := rows.Scan(&l.ID, &l.Shortcut, &l.Destination, &createdAt, &l.ClickCount); err != nil {
			return nil, err
		}
		l.CreatedAt, err = time.Parse(timeLayout, createdAt)
		if err != nil {
			return nil, err
		}
		links = append(links, &l)
	}
	return links, rows.Err()
}

func (r *LinkRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM links WHERE id = ?`, id)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if x == 0 {
		return NoError
	}
	return nil
}

func (r *LinkRepository) IncrementClickCount(shortcut string) (*models.Link, error) {
	res, err := r.db.Exec(
		`UPDATE links SET click_count = click_count + 1 WHERE shortcut = ?`,
		shortcut,
	)
	if err != nil {
		return nil, err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if x == 0 {
		return nil, NoError
	}
	return r.GetByShortcut(shortcut)
}

func scanLink(row *sql.Row) (*models.Link, error) {
	var l models.Link
	var createdAt string
	err := row.Scan(&l.ID, &l.Shortcut, &l.Destination, &createdAt, &l.ClickCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NoError
	}
	if err != nil {
		return nil, err
	}
	l.CreatedAt, err = time.Parse(timeLayout, createdAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func isErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
