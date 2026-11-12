package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckWorkOutsideRepo(t *testing.T) {
	mk := func(t *testing.T, paths ...string) string {
		t.Helper()
		root := t.TempDir()
		for _, p := range paths {
			if err := os.MkdirAll(filepath.Join(root, p), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	t.Run("work tree", func(t *testing.T) {
		root := mk(t, ".git", "a/work")
		dir := filepath.Join(root, "a/work")
		err := checkWorkOutsideRepo(dir)
		if err == nil || !strings.Contains(err.Error(), "GOTMPDIR") || !strings.Contains(err.Error(), "work") {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("git dir", func(t *testing.T) {
		root := mk(t, ".git/rhiza/tmp")
		os.RemoveAll(filepath.Join(root, ".git", "x"))
		err := checkWorkOutsideRepo(filepath.Join(root, ".git/rhiza/tmp"))
		if err == nil || !strings.Contains(err.Error(), "GOTMPDIR") || !strings.Contains(err.Error(), ".git") {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("bare repo", func(t *testing.T) {
		root := mk(t, "r.git/objects", "r.git/refs", "r.git/tmp")
		os.WriteFile(filepath.Join(root, "r.git/HEAD"), []byte("ref: refs/heads/main\n"), 0o644)
		if err := checkWorkOutsideRepo(filepath.Join(root, "r.git/tmp")); err == nil {
			t.Fatal("want error")
		}
	})
	t.Run("outside", func(t *testing.T) {
		if err := checkWorkOutsideRepo(mk(t, "w")); err != nil {
			t.Fatal(err)
		}
	})
}
