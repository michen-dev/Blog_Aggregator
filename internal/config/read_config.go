package config

import (
	"encoding/json"
	"fmt"
	"os"
)


func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get home directory: %v", err) 
	}
	if homeDir == "" {
		return homeDir, fmt.Errorf("Expected home directory path non-empty")
	}
	configFilePath := homeDir + config_json_file_path
	return configFilePath, nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Failed retrieving config file path: %v", err)
	}

	fileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading file %v", err)
	}

	var config Config

	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing JSON: %v", err)
	}

	return config, nil
}	