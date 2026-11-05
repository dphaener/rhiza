// Package tmpsweep removes the temp folders a killed test run left behind.
// A test that makes a fixed-prefix folder under os.TempDir() removes it when
// it exits normally, but a go test timeout kills the binary first, and each
// such run leaves its folder for good.
package tmpsweep

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MaxAge is how old a folder's modification time must be for Stale to
// remove it: a folder younger than this may belong to a run still going.
const MaxAge = time.Hour

// Stale removes the folders directly under root whose name starts with
// prefix and whose modification time is older than MaxAge. Files, folders
// with another prefix and fresh folders stay. Errors are ignored: this is
// housekeeping and must never fail a test run.
func Stale(root, prefix string) {
	stale(root, prefix, time.Now().Add(-MaxAge))
}

// StaleTemp is Stale over os.TempDir().
func StaleTemp(prefix string) {
	Stale(os.TempDir(), prefix)
}

func stale(root, prefix string, cutoff time.Time) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), prefix) {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.ModTime().Before(cutoff) {
			continue
		}
		os.RemoveAll(filepath.Join(root, e.Name()))
	}
}
