package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents application configuration loaded from a JSON file.
//
// Fields:
//   - DBFileName: path to the SQLite/SQLCipher database file.
//   - BusyTimeoutInMs: SQLite busy timeout in milliseconds.
//   - ServerAddr: address of the gRPC server (host:port).
//   - ServerName: TLS server name for certificate verification.
//   - Debug: enables debug mode (e.g. plaintext gRPC, verbose logs).
type Config struct {
	DBFileName      string `json:"db_file_name"`
	BusyTimeoutInMs int    `json:"busy_timeout_in_ms"`
	ServerAddr      string `json:"server_addr"`
	ServerName      string `json:"server_name"`
	Debug           bool   `json:"debug"`
}

// LoadFromFile loads configuration from a JSON file located in the same
// directory as the provided executable path.
func LoadFromFile(execPath, fileName string) (Config, error) {
	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, fileName)

	cfg := Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file (%s): %w", configPath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	if cfg.DBFileName != "" {
		cfg.DBFileName = filepath.Join(execDir, cfg.DBFileName)
	}

	return cfg, nil
}
