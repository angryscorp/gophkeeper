package help

import (
	"strings"
	"testing"
)

func TestHelp_OutputContainsBuildInfo(t *testing.T) {
	h := New("v1.0.0", "10.10.2025")
	out := h.Help()

	if !strings.Contains(out, "Build version: v1.0.0") {
		t.Errorf("expected version in output, got %s", out)
	}
	if !strings.Contains(out, "Build date: 10.10.2025") {
		t.Errorf("expected date in output, got %s", out)
	}
}
