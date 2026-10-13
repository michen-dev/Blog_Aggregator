package config

import (
	"encoding/json"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Fatal("expected a non-empty home directory path")
	}
	t.Logf("Config path: %s", path)
}

func TestRead(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Fatalf("Failed reading: %v", err)
	}
	data, _ := json.MarshalIndent(config, "", " ")
	t.Log(string(data))
}

func TestSetUser(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Fatalf("Failed reading: %v", err)
	}

	err = config.SetUser("michen_dev")
	if err != nil {
		t.Fatalf("Failed setting user: %v", err)
	}

	data, _ := json.MarshalIndent(config, "", " ")
	t.Log("Current config: ")
	t.Log(string(data))

	new_config, err := Read()
	if err != nil {
		t.Fatalf("Failed reading: %v", err)
	}
	data, _ = json.MarshalIndent(new_config, "", " ")
	t.Log("New config: ")
	t.Log(string(data))
}