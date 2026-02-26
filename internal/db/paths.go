package db

import "fmt"

type Path struct {
	ID        int
	Name      string
	Path      string
	CreatedAt string
}

func AddPath(name, path string) error {
	_, err := GetDB().Exec(
		"INSERT INTO paths (name, path) VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET path = ?",
		name, path, path,
	)
	return err
}

func GetPath(name string) (string, error) {
	var path string
	err := GetDB().QueryRow("SELECT path FROM paths WHERE name = ?", name).Scan(&path)
	if err != nil {
		return "", fmt.Errorf("no path found for '%s'", name)
	}
	return path, nil
}

func ListPaths() ([]Path, error) {
	rows, err := GetDB().Query("SELECT id, name, path, created_at FROM paths ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []Path
	for rows.Next() {
		var p Path
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	return paths, rows.Err()
}

func RemovePath(name string) error {
	result, err := GetDB().Exec("DELETE FROM paths WHERE name = ?", name)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no path found for '%s'", name)
	}
	return nil
}

func LoadPathsMap() map[string]string {
	paths := make(map[string]string)
	all, err := ListPaths()
	if err != nil {
		return paths
	}
	for _, p := range all {
		paths[p.Name] = p.Path
	}
	return paths
}
