//go:build unix && !linux

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReapAdoptOther: off Linux adoption lists nothing.
func TestReapAdoptOther(t *testing.T) {
	common := t.TempDir()
	adoptChildren(common)
	if _, err := os.Stat(filepath.Join(common, "rhiza", "drivers")); err == nil {
		t.Error("adoption wrote a driver list off Linux")
	}
}
