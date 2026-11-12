package rootcfg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitFindProjectOf(t *testing.T) {
	top := filepath.Join(t.TempDir(), "top")
	if _, err := Init(top); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(top); !errors.Is(err, ErrExists) {
		t.Fatalf("second Init: %v", err)
	}
	cfgText := "[projects.app]\npath = \"app\"\n"
	if err := os.WriteFile(filepath.Join(top, File), []byte(cfgText), 0o644); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(top, "app.x", "sub")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	root, ok := Find(deep)
	if !ok || root != top {
		t.Fatalf("Find = %q, %v", root, ok)
	}
	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	for dir, want := range map[string]string{
		top:                                 "",
		filepath.Join(top, "app"):           "app",
		filepath.Join(top, "app", "a", "b"): "app",
		deep:                                "app",
		filepath.Join(top, "apple"):         "",
		filepath.Join(top, "app."):          "",
	} {
		if got, _ := cfg.ProjectOf(root, dir); got != want {
			t.Errorf("ProjectOf(%s) = %q, want %q", dir, got, want)
		}
	}
}

// An empty dir is no directory, not the process's cwd: Find reports no root
// even when the cwd sits under a folder holding rhiza.toml.
func TestFindEmptyDirFindsNoRoot(t *testing.T) {
	top := t.TempDir()
	if _, err := Init(top); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(top, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(sub)
	if root, ok := Find(""); ok {
		t.Fatalf("Find(\"\") = %q, true: want no root", root)
	}
	if root, ok := Find("."); !ok || root != top {
		t.Fatalf("Find(\".\") = %q, %v: want %q", root, ok, top)
	}
}

func TestLoadRejectsRemovedSelfUpdate(t *testing.T) {
	for _, toml := range []string{
		"self_update_project = \"rhiza\"\n",
		"[daemon]\n\"rhiza.daemon.selfUpdate\" = \"true\"\n",
	} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, File), []byte(toml), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := Load(root)
		if !errors.Is(err, ErrRemoved) || !strings.Contains(err.Error(), "rhiza.daemon.autoUpdate") {
			t.Fatalf("Load(%q) = %v", toml, err)
		}
	}
}
