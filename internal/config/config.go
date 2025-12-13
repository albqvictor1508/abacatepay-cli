package config

import (
	"encoding/json"
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

func (c *Config) Save(name string, key string) error {
	path, err := getPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	const permUserReadWriteExec = 0o755

	if err := os.MkdirAll(dir, permUserReadWriteExec); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := SaveKeyring(name, key); err != nil {
		return err
	}

	const permFileReadable = 0o644

	return os.WriteFile(path, data, permFileReadable)
}

func (c *Config) Exists(name string) bool {
	// Trimmar a string no futuro? name = strings.TrimSpace(name)

	_, exists := c.Profiles[name]

	return exists
}

func (c *Config) Add(name, key string) error {
	c.Profiles[name] = Profile{
		Verified:  true,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return c.Save(name, key)
}

func SaveCurrent(name string, key string) error {
	config, err := Load()
	if err != nil {
		return err
	}

	config.Current = name

	return config.Save(name, key)
}

func (c *Config) ProfileExists(name string) bool {
	_, exists := c.Profiles[name]
	return exists
}
