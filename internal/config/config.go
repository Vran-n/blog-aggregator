package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func getConfigPath() string {
	const path = "/.gatorconfig.json"

	home_dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return home_dir + path
}

func Read() (Config, error) {
	cp := getConfigPath()
	config_data, err := os.ReadFile(cp)
	if err != nil {
		return Config{}, fmt.Errorf("Cannot read file of path \"%s\": %w", cp, err)
	}

	var config Config
	if err := json.Unmarshal(config_data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(data []byte) error {
	cp := getConfigPath()
	err := os.WriteFile(cp, data, 0666)
	if err != nil {
		return fmt.Errorf("Cannot write file of path \"%s\": %w", cp, err)
	}

	return nil
}

func SetUser(config *Config, name string) error {
	config.CurrentUsername = name

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := write(data); err != nil {
		return err
	}

	return nil
}
