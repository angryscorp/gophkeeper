package buildinfo_test

import (
	"testing"

	"gophkeeper/pkg/buildinfo"
)

func TestBuildInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		date     string
		expected string
	}{
		{
			name:    "both values provided",
			version: "v1.0.0",
			date:    "2025-10-07",
			expected: "Build version: v1.0.0\n" +
				"Build date: 2025-10-07\n",
		},
		{
			name:    "missing version",
			version: "",
			date:    "2025-10-07",
			expected: "Build version: N/A\n" +
				"Build date: 2025-10-07\n",
		},
		{
			name:    "missing both",
			version: "",
			date:    "",
			expected: "Build version: N/A\n" +
				"Build date: N/A\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := buildinfo.New(tt.version, tt.date)
			if got := bi.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
