package db

import "fmt"

type Alias struct {
	ID        int
	AliasName string
	Command   string
	CreatedAt string
}

func SetAlias(name, command string) error {
	_, err := GetDB().Exec(
		"INSERT INTO aliases (alias_name, command) VALUES (?, ?) ON CONFLICT(alias_name) DO UPDATE SET command = ?",
		name, command, command,
	)
	return err
}

func GetAlias(name string) (string, error) {
	var command string
	err := GetDB().QueryRow("SELECT command FROM aliases WHERE alias_name = ?", name).Scan(&command)
	if err != nil {
		return "", fmt.Errorf("no alias found for '@%s'", name)
	}
	return command, nil
}

func ListAliases() ([]Alias, error) {
	rows, err := GetDB().Query("SELECT id, alias_name, command, created_at FROM aliases ORDER BY alias_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aliases []Alias
	for rows.Next() {
		var a Alias
		if err := rows.Scan(&a.ID, &a.AliasName, &a.Command, &a.CreatedAt); err != nil {
			return nil, err
		}
		aliases = append(aliases, a)
	}
	return aliases, rows.Err()
}

func RemoveAlias(name string) error {
	result, err := GetDB().Exec("DELETE FROM aliases WHERE alias_name = ?", name)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no alias found for '@%s'", name)
	}
	return nil
}
