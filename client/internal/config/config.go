package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBFileName      string `json:"db_file_name"`
	BusyTimeoutInMs int    `json:"busy_timeout_in_ms"`
	ServerAddr      string `json:"server_addr"`
	Debug           bool   `json:"debug"`
}

func LoadFromFile(filePath string) (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file (%s): %w", filePath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	return cfg, nil
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf("file:%s?_pragma=busy_timeout=%d", c.DBFileName, c.BusyTimeoutInMs)
}
