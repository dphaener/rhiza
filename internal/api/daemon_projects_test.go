package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

func TestDaemonHomePrefersWorkingProject(t *testing.T) {
	root := t.TempDir()
	for _, n := range []string{"a", "b"} {
		if err := os.MkdirAll(filepath.Join(root, n, ".rhiza"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, n, ".rhiza", "repo.json"), []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "bare"), 0o755); err != nil {
		t.Fatal(err)
	}
	reg := rootcfg.Config{Projects: map[string]rootcfg.Project{
		"a": {Path: "a"}, "b": {Path: "b"}, "bare": {Path: "bare"},
	}}
	for wd, want := range map[string]string{
		filepath.Join(root, "b"):    "b",
		filepath.Join(root, "bare"): "a", // not initialised: the first initialised one
		root:                        "a",
	} {
		name, _, ok := daemonHome(&Ctx{Wd: wd, Root: wd}, root, reg)
		if !ok || name != want {
			t.Errorf("daemonHome from %s = %q, %v; want %q", wd, name, ok, want)
		}
	}
	if projectInitialised(filepath.Join(root, "bare")) {
		t.Error("a bare checkout counts as initialised")
	}
}
