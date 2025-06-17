// Package main defines the entry point of the backend.
package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Exit(run())
}
