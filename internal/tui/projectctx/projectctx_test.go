package projectctx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

func TestHeader(t *testing.T) {
	root := t.TempDir()
	if _, err := rootcfg.Init(root); err != nil {
		t.Fatal(err)
	}
	if err := rootcfg.Add(root, "alpha", rootcfg.Project{Path: "alpha"}); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"alpha/sub", "alpha.feature/sub", "other"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct{ name, dir, want string }{
		{"project folder", "alpha", "Project: alpha"},
		{"inside the project", "alpha/sub", "Project: alpha"},
		{"sibling worktree", "alpha.feature/sub", "Project: alpha"},
		{"root", ".", NoProject},
		{"unregistered folder", "other", NoProject},
	}
	for _, c := range cases {
		if got := Header(filepath.Join(root, c.dir)); got != c.want {
			t.Errorf("%s: header %q, want %q", c.name, got, c.want)
		}
	}
	if got := Header(t.TempDir()); got != NoProject {
		t.Errorf("outside any root: header %q, want %q", got, NoProject)
	}
}
