package config

import (
	"os"
	"path/filepath"
)

type ProfileData struct {
	CreatedAt string `json:"created_at"`
	Verified  bool   `json:"verified,omitempty"`
}

type Config struct {
	Profiles map[string]ProfileData `json:"profiles"`
	Current  string                 `json:"current,omitempty"`
}

func getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".abacate", "abacate.json"), nil
}

func Load() (*Config, error) {
	return nil
}
