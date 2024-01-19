package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Config   ConfigFile
	Location string
}

type ConfigFile struct {
	GitHubKey string `json:"github_key"`
	GeminiKey string `json:"gemini_key"`
	VertexKey string `json:"vertex_key"`
}

type NewConfigError struct {
	Location string
}

func (e *NewConfigError) Error() string {
	return fmt.Sprintf("Config file was just created at : %s\nPlease visit that location and fill in the config values", e.Location)
}

func LoadConfig(Location string) (*Config, error) {
	config := Location
	c := &Config{Location: config}
	var err error = nil
	if Location == "" {
		config, err = configLocation()
		if err != nil {
			return nil, err
		}
		c.Location = config
	}
	_, err = os.Stat(config)
	if os.IsNotExist(err) {
		err := createConfig(config)
		if err != nil {
			return c, err
		}
		return c, errors.New("Config file is new")
	}
	err = c.loadConfigValues()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) loadConfigValues() error {
	file, err := os.Open(c.Location)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	json.Unmarshal(b, &c.Config)
	return nil

}

func configLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, ".config", "expensive-water.json")
	return configPath, nil
}

func configExists(Location string) (bool, error) {

	config := Location
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func createConfig(Location string) error {
	config := Location
	file, err := os.Create(config)
	if err != nil {
		return err
	}
	defer file.Close()

	defaultConfig := &ConfigFile{GitHubKey: "INSERT KEY HERE", GeminiKey: "INSERT KEY HERE", VertexKey: "INSERT KEY HERE"}
	b, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) GetGitHubKey() {

}

func (c *Config) GetGeminiKey() {

}

func (c *Config) GetVertexKey() {

}
