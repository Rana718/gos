package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Rana718/gos/internal/db"
)

type Config struct {
	Editor string `json:"editor"`
}

func Default() Config {
	return Config{
		Editor: "vim",
	}
}

func Load() (Config, error) {
	configPath := filepath.Join(db.GosDir(), "config.json")
	cfg := Default()

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		if err := Save(cfg); err != nil {
			return cfg, err
		}
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	configPath := filepath.Join(db.GosDir(), "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func OpenInEditor() {
	cfg, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(db.GosDir(), "config.json")

	cmd := exec.Command(cfg.Editor, configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config in %s: %v\n", cfg.Editor, err)
		os.Exit(1)
	}
}
