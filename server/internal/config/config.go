package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	DatabaseDSN    string `env:"DATABASE_DNS"`
	ServerAddr     string `env:"SERVER_ADDR"`
	Debug          bool   `env:"DEBUG"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
	TLSCertPath    string `env:"TLS_CERT_PATH"`
	TLSCKeyPath    string `env:"TLS_KEY_PATH"`
}

// LoadFromEnv parses environment variables into a Config struct.
func LoadFromEnv() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse env config: %w", err)
	}
	return cfg, nil
}
