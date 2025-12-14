package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Profile struct {
	CreatedAt string `json:"created_at"`
	Verified  bool   `json:"verified,omitempty"`
}

type Config struct {
	Profiles map[string]Profile `json:"profiles"`
	Current  string             `json:"current,omitempty"`
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
				Profiles: make(map[string]Profile),
			}, nil
		}

		return nil, err
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Profiles == nil {
		config.Profiles = make(map[string]Profile)
	}

	return &config, nil
}

func (c *Config) saveConfig() error {
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

func (c *Config) Save(name string, key string) error {
	if err := SaveKeyring(name, key); err != nil {
		return err
	}

	return c.saveConfig()
}

func (c *Config) Exists(name string) bool {
	_, exists := c.Profiles[name]
	return exists
}

func (c *Config) Add(name, key string) error {
	c.Profiles[name] = Profile{
		Verified:  true,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	c.Current = name

	if err := SaveKeyring(name, key); err != nil {
		return fmt.Errorf("failed to save key for profile '%s': %w", name, err)
	}

	if err := SaveKeyring("current", key); err != nil {
		return fmt.Errorf("failed to set 'current' key: %w", err)
	}

	return c.saveConfig()
}

func (c *Config) SetCurrent(name string) error {
	if !c.Exists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	key, err := GetKeyring(name)
	if err != nil {
		return fmt.Errorf("could not retrieve key for profile '%s': %w", name, err)
	}

	if err := SaveKeyring("current", key); err != nil {
		return fmt.Errorf("failed to set 'current' key: %w", err)
	}

	c.Current = name
	return c.saveConfig()
}

func (c *Config) ProfileExists(name string) bool {
	_, exists := c.Profiles[name]
	return exists
}
