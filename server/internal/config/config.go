package config

type Config struct {
	DatabaseDSN string `env:"DATABASE_DNS"`
	ServerAddr  string `env:"SERVER_ADDR"`
	Debug       bool   `env:"DEBUG"`
}
