// Package db contains the database connection and the models for the database.
package db

import (
	"database/sql"
	"log"
)

// DB is the database
type DB struct {
	conn *sql.DB
}

// Close closes the database connection
func (d *DB) Close() {
	d.conn.Close()
}

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
		log.Fatal(err)
	}
	return count
}

// SetCount return the latest count
func (d *DB) SetCount(count int) {
	q := "INSERT INTO count (count) VALUES (?)"
	_, err := d.conn.Exec(q, count)
	if err != nil {
		log.Fatal(err)
	}
}
