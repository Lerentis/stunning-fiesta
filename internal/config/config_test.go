package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureConfigAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-stunning-fiesta.yaml")

	// Ensure config file is created
	err := EnsureConfig(configPath)
	if err != nil {
		t.Fatalf("EnsureConfig failed: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Load config and check default values
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Endpoints.Template == "" ||
		cfg.Endpoints.Update == "" ||
		cfg.Endpoints.ClusterInfo == "" ||
		cfg.GitlabURL == "" {
		t.Errorf("Default config values are missing or empty")
	}
}
