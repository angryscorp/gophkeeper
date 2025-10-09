package device_test

import (
	"runtime"
	"strings"
	"testing"

	"gophkeeper/pkg/device"
)

func TestGenerateDeviceName(t *testing.T) {
	got := device.GenerateDeviceName()

	if got == "" {
		t.Fatal("expected non-empty device name")
	}
	if got != strings.ToLower(got) {
		t.Errorf("expected lowercase result, got %q", got)
	}
	if strings.Contains(got, " ") {
		t.Errorf("expected no spaces in result, got %q", got)
	}
	if !strings.Contains(got, runtime.GOOS) {
		t.Errorf("expected device name to contain GOOS %q, got %q", runtime.GOOS, got)
	}
	if !strings.Contains(got, runtime.GOARCH) {
		t.Errorf("expected device name to contain GOARCH %q, got %q", runtime.GOARCH, got)
	}
}
