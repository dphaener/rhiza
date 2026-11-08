package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultsWhenNoFile(t *testing.T) {
	r, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range Kinds {
		if r.ForRequest(k) != "plan" || r.ForReady(k) != "ship-change" {
			t.Errorf("%s: got %q/%q", k, r.ForRequest(k), r.ForReady(k))
		}
	}
}

func TestReadyBugOnly(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, File), []byte("[ready]\nbug = \"bugfix\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	d := Defaults()
	for _, k := range Kinds {
		want := d.ForReady(k)
		if k == KindBug {
			want = "bugfix"
		}
		if r.ForReady(k) != want || r.ForRequest(k) != "plan" {
			t.Errorf("%s: got %q/%q", k, r.ForRequest(k), r.ForReady(k))
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct{ in, want string }{
		{"[other]\nbug = \"x\"\n", "[other]"},
		{"[request]\nchore = \"x\"\n", "[request] unknown key \"chore\""},
		{"[ready]\nmilestone = \"x\"\n", "[ready] unknown key \"milestone\": want one of feature, bug"},
		{"[ready]\nbug = 3\n", "[ready] bug"},
		{"[ready]\nbug = \"a/b\"\n", "[ready] bug"},
		{"[request]\nfeature = \"x.toml\"\n", "[request] feature"},
		{"[request]\nfeature = \"\"\n", "[request] feature"},
	}
	for _, c := range cases {
		_, err := Parse([]byte(c.in))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: got %v, want %q", c.in, err, c.want)
		}
	}
}
