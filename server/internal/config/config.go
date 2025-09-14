package config

type Config struct {
	DatabaseDSN string `env:"DATABASE_DNS"`
	Debug       bool   `env:"DEBUG"`
}
