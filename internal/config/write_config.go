package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Failed retrieving config file path: %v", err)
	}
	
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %v", err)
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed writing to file: %v", err)
	}
	
	return nil
}

func (c *Config) SetUser(current_user_name string) error {
	c.Current_user_name = current_user_name
	err := write(*c)
	if err != nil {
		return fmt.Errorf("Unexpected error: %v", err)
	}
	return nil
}