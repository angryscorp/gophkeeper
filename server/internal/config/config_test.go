package config_test

import (
	"os"
	"testing"

	"gophkeeper/server/internal/config"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("DATABASE_DNS", "postgres://user:pass@localhost:5432/db")
	t.Setenv("SERVER_ADDR", ":8443")
	t.Setenv("DEBUG", "true")
	t.Setenv("PRIVATE_KEY_PATH", "/tmp/private.pem")
	t.Setenv("PUBLIC_KEY_PATH", "/tmp/public.pem")
	t.Setenv("TLS_CERT_PATH", "/tmp/cert.pem")
	t.Setenv("TLS_KEY_PATH", "/tmp/key.pem")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DatabaseDSN != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("DatabaseDSN mismatch: got %q", cfg.DatabaseDSN)
	}
	if cfg.ServerAddr != ":8443" {
		t.Errorf("ServerAddr mismatch: got %q", cfg.ServerAddr)
	}
	if cfg.Debug != true {
		t.Errorf("Debug should be true")
	}
	if cfg.PrivateKeyPath != "/tmp/private.pem" {
		t.Errorf("PrivateKeyPath mismatch")
	}
	if cfg.PublicKeyPath != "/tmp/public.pem" {
		t.Errorf("PublicKeyPath mismatch")
	}
	if cfg.TLSCertPath != "/tmp/cert.pem" {
		t.Errorf("TLSCertPath mismatch")
	}
	if cfg.TLSCKeyPath != "/tmp/key.pem" {
		t.Errorf("TLSCKeyPath mismatch")
	}
}

func TestLoadFromEnv_MissingVars(t *testing.T) {
	os.Clearenv()

	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error when env vars missing, got nil")
	}
}
