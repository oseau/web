// Package db contains the database connection and the models for the database.
package db

import (
	"database/sql"
	"log/slog"
)

// Count is a row in count table
type Count struct {
	Count int `json:"count"`
}

// GetCount return the latest count
func (d *DB) GetCount() Count {
	q := "SELECT count FROM count ORDER BY id DESC LIMIT 1"
	row := d.conn.QueryRow(q)
	count := Count{}
	if err := row.Scan(&count.Count); err != nil && err != sql.ErrNoRows {
		slog.Error("failed to get count", "error", err)
	}
	return count
}

// SetCount insert a new count record
func (d *DB) SetCount(count int) {
	q := "INSERT INTO count (count) VALUES (?)"
	if _, err := d.conn.Exec(q, count); err != nil {
		slog.Error("failed to insert count", "error", err)
	}
}
