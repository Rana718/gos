package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	instance *sql.DB
	once     sync.Once
	gosDir   string
)

func GosDir() string {
	if gosDir != "" {
		return gosDir
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	gosDir = filepath.Join(homeDir, ".gos")
	if err := os.MkdirAll(gosDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .gos directory: %v\n", err)
		os.Exit(1)
	}
	return gosDir
}

func GetDB() *sql.DB {
	once.Do(func() {
		dbPath := filepath.Join(GosDir(), "gos.db")
		var err error
		instance, err = sql.Open("sqlite", dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		instance.SetMaxOpenConns(1)

		if err := createTables(instance); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating tables: %v\n", err)
			os.Exit(1)
		}

		if err := MigrateFromJSON(instance, GosDir()); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: JSON migration failed: %v\n", err)
		}
	})
	return instance
}

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS paths (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS aliases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alias_name TEXT NOT NULL UNIQUE,
			command TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute: %w", err)
		}
	}
	return nil
}
