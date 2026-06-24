package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
	Connection_string string `json:"connection_string"`
}

func getConfigFilepath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(dir, ".gatorconfig.json")
	return filepath, nil
}

func write(cfg Config) error {
	filepath, err := getConfigFilepath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	filepath, err := getConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) SetUser(user string) error {
	c.Current_user_name = user
	if err := write(c); err != nil {
		return err
	}

	return nil
}
