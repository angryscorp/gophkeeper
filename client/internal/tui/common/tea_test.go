package common

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

func TestInputWithPlaceholder_SetsPlaceholderAndWidth(t *testing.T) {
	inp := InputWithPlaceholder("hello")
	if inp.Placeholder != "hello" {
		t.Fatalf("placeholder != %q, got %q", "hello", inp.Placeholder)
	}
	if inp.Width != len("hello") {
		t.Fatalf("width != %d, got %d", len("hello"), inp.Width)
	}
}

func TestExtendInputWidthIfNeeded_GrowsOnlyWhenNeeded(t *testing.T) {
	var m textinput.Model
	m.Placeholder = "xx"
	m.Width = 2

	// smaller than width -> unchanged
	m.SetValue("a")
	ExtendInputWidthIfNeeded(&m)
	if m.Width != 2 {
		t.Fatalf("width changed unexpectedly, got %d", m.Width)
	}

	// larger than width -> grows
	m.SetValue("abcdef")
	ExtendInputWidthIfNeeded(&m)
	if m.Width != len("abcdef") {
		t.Fatalf("width not extended, want %d, got %d", len("abcdef"), m.Width)
	}
}
