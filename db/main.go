// Package db contains the database connection and the models for the database.
package db

import (
	"database/sql"
	"log/slog"
	"os"
	"sync"
)

var (
	once sync.Once
	db   *DB
)

// DB is the database
type DB struct {
	conn *sql.DB
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.conn.Close()
}

// NewDB creates a new database connection
func NewDB() *DB {
	once.Do(func() {
		conn, err := sql.Open("sqlite3", "/data/web.db")
		if err != nil {
			slog.Error("failed to open database", "error", err)
			os.Exit(1)
		}
		db = &DB{conn: conn}
	})
	return db
}
