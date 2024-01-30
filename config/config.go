// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"encoding/json"
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
	Region    string `json:"region"`
	ProjectID string `json:"project_id"`
}

type NewConfigError struct {
	Location string
}

func (e *NewConfigError) Error() string {
	return fmt.Sprintf("Config file was just created at : %s\nPlease visit that location and fill in the config values", e.Location)
}

func LoadConfig(Location string) (*Config, error) {
	configLoc := Location
	c := &Config{Location: configLoc}
	var err error = nil
	if Location == "" {
		configLoc, err = configLocation()
		if err != nil {
			return nil, err
		}
		c.Location = configLoc
	}
	_, err = os.Stat(configLoc)
	if os.IsNotExist(err) {
		err := createConfig(configLoc)
		if err != nil {
			return c, err
		}
		return c, &NewConfigError{Location: configLoc}
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

	defaultConfig := &ConfigFile{GitHubKey: "INSERT KEY HERE", GeminiKey: "INSERT KEY HERE", VertexKey: "INSERT KEY HERE", Region: "us-central1", ProjectID: "Only required if using vertexAI"}
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

func (c *Config) GetGitHubKey() string {
	return c.Config.GitHubKey
}

func (c *Config) GetGeminiKey() string {
	return c.Config.GeminiKey
}

func (c *Config) GetVertexKey() string {
	return c.Config.VertexKey
}

func (c *Config) GetRegion() string {
	return c.Config.Region
}

func (c *Config) GetProjectID() string {
	return c.Config.ProjectID
}
