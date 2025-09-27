package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBFileName      string `json:"db_file_name"`
	BusyTimeoutInMs int    `json:"busy_timeout_in_ms"`
	ServerAddr      string `json:"server_addr"`
	Debug           bool   `json:"debug"`
}

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
