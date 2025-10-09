package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile_Success_JoinDBPath(t *testing.T) {
	t.Parallel()

	// Arrange
	execDir := t.TempDir()
	execPath := filepath.Join(execDir, "app-bin")
	cfgName := "config.json"

	json := `{
		"db_file_name": "vault.db",
		"busy_timeout_in_ms": 2500,
		"server_addr": "127.0.0.1:50051",
		"server_name": "localhost",
		"debug": true
	}`

	if err := os.WriteFile(filepath.Join(execDir, cfgName), []byte(json), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Act
	cfg, err := LoadFromFile(execPath, cfgName)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}

	// Assert
	wantDB := filepath.Join(execDir, "vault.db")
	if cfg.DBFileName != wantDB {
		t.Errorf("DBFileName: got %q, want %q", cfg.DBFileName, wantDB)
	}
	if cfg.BusyTimeoutInMs != 2500 {
		t.Errorf("BusyTimeoutInMs: got %d, want %d", cfg.BusyTimeoutInMs, 2500)
	}
	if cfg.ServerAddr != "127.0.0.1:50051" {
		t.Errorf("ServerAddr: got %q", cfg.ServerAddr)
	}
	if cfg.ServerName != "localhost" {
		t.Errorf("ServerName: got %q", cfg.ServerName)
	}
	if !cfg.Debug {
		t.Errorf("Debug: got false, want true")
	}
}

func TestLoadFromFile_ReadError(t *testing.T) {
	t.Parallel()

	execDir := t.TempDir()
	execPath := filepath.Join(execDir, "app-bin")

	_, err := LoadFromFile(execPath, "nope.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromFile_BadJSON(t *testing.T) {
	t.Parallel()

	execDir := t.TempDir()
	execPath := filepath.Join(execDir, "app-bin")
	cfgName := "bad.json"

	bad := `{"db_file_name": "vault.db", "debug": tru`
	if err := os.WriteFile(filepath.Join(execDir, cfgName), []byte(bad), 0o644); err != nil {
		t.Fatalf("write bad config: %v", err)
	}

	_, err := LoadFromFile(execPath, cfgName)
	if err == nil {
		t.Fatal("expected JSON parse error, got nil")
	}
}

func TestLoadFromFile_EmptyDBName_NoJoin(t *testing.T) {
	t.Parallel()

	execDir := t.TempDir()
	execPath := filepath.Join(execDir, "app-bin")
	cfgName := "cfg.json"

	json := `{
		"db_file_name": "",
		"busy_timeout_in_ms": 1000,
		"server_addr": "s:1",
		"server_name": "",
		"debug": false
	}`
	if err := os.WriteFile(filepath.Join(execDir, cfgName), []byte(json), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFromFile(execPath, cfgName)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.DBFileName != "" {
		t.Errorf("DBFileName: got %q, want empty", cfg.DBFileName)
	}
}
