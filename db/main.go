// Package db contains the database connection and the models for the database.
package db

import (
	"database/sql"
	"log"
	"sync"
)

var once sync.Once

// NewDB creates a new database connection
func NewDB() *DB {
	db := &DB{}
	once.Do(func() {
		conn, err := sql.Open("sqlite3", "/db/web.db")
		if err != nil {
			log.Fatal(err)
		}
		db.conn = conn
	})
	return db
}
