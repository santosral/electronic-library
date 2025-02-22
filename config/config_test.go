package config

import (
	"fmt"
	"os"
	"testing"
)

func createTempConfigFile(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		return "", fmt.Errorf("could not create temporary file: %v", err)
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("could not write to temp file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("could not close temp file: %v", err)
	}

	return tempFile.Name(), nil
}

func TestLoadConfig_Success(t *testing.T) {
	configContent := `{
		"server": {
			"host": "localhost",
			"port": 8080
		}
	}`
	tempFile, err := createTempConfigFile(configContent)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	config, err := LoadConfig(tempFile)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("expected server host 'localhost', got: %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got: %d", config.Server.Port)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent-file.json")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	invalidJSON := `Invalid JSON`

	tempFile, err := createTempConfigFile(invalidJSON)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	_, err = LoadConfig(tempFile)
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}
