package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS links (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	shortcut      TEXT NOT NULL UNIQUE,
	destination   TEXT NOT NULL,
	created_at    TEXT NOT NULL,
	click_count   INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_links_created_at ON links(created_at);
`

func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	conn.SetMaxOpenConns(1)

	if _, err := conn.Exec(schema); err != nil {
		conn.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return conn, nil
}
