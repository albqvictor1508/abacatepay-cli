package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
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
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Profiles: make(map[string]ProfileData),
			}, nil
		}

		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]ProfileData)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := getPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func (c *Config) ProfileExists(name string) bool {
	_, exists := c.Profiles[name]
	return exists
}

func (c *Config) AddProfile(name, apiKey string) error {
	c.Profiles[name] = ProfileData{
		CreatedAt: time.Now().Format(time.RFC3339),
		Verified:  true,
	}

	return c.Save()
}

func SaveCurrent(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	cfg.Current = name
	return cfg.Save()
}
