package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func MigrateFromJSON(database *sql.DB, dir string) error {
	jsonPath := filepath.Join(dir, "paths.json")
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		return nil
	}

	var count int
	database.QueryRow("SELECT COUNT(*) FROM paths").Scan(&count)
	if count > 0 {
		return nil
	}

	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open paths.json: %w", err)
	}
	defer file.Close()

	paths := make(map[string]string)
	if err := json.NewDecoder(file).Decode(&paths); err != nil {
		return fmt.Errorf("failed to decode paths.json: %w", err)
	}

	tx, err := database.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO paths (name, path) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for name, path := range paths {
		if _, err := stmt.Exec(name, path); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return os.Rename(jsonPath, jsonPath+".bak")
}
