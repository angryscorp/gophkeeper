package pgx_test

import (
	"strings"
	"testing"

	"gophkeeper/pkg/pgx"
)

func TestCreatePGXPool_InvalidDSN(t *testing.T) {
	_, err := pgx.CreatePGXPool("not-a-valid-dsn")
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse dsn") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreatePGXPool_UnreachableDB(t *testing.T) {
	// DSN синтаксически валидный, но база явно не существует.
	_, err := pgx.CreatePGXPool("postgres://invalid:invalid@localhost:5432/doesnotexist")
	if err == nil {
		t.Fatal("expected error for unreachable DB, got nil")
	}
	if !strings.Contains(err.Error(), "failed to ping database") {
		t.Errorf("unexpected error: %v", err)
	}
}
