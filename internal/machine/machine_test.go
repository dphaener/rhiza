package machine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, body string) string {
	t.Helper()
	p := filepath.Join(dir, "rhiza", "machine.toml")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestAbsentIsDaily(t *testing.T) {
	home := t.TempDir()
	c, err := Load(func(string) string { return "" }, home)
	if err != nil || c.Trust != Daily {
		t.Fatalf("got %+v, %v", c, err)
	}
}

func TestTrustValues(t *testing.T) {
	for _, want := range []string{Dedicated, Daily} {
		dir := t.TempDir()
		write(t, dir, "trust = \""+want+"\"\n")
		c, err := Load(func(k string) string {
			if k == "XDG_CONFIG_HOME" {
				return dir
			}
			return ""
		}, "/nonexistent")
		if err != nil || c.Trust != want {
			t.Fatalf("%s: got %+v, %v", want, c, err)
		}
	}
}

func TestHomeDefault(t *testing.T) {
	home := t.TempDir()
	write(t, filepath.Join(home, ".config"), "trust = \"dedicated\"\n")
	c, err := Load(func(string) string { return "" }, home)
	if err != nil || c.Trust != Dedicated {
		t.Fatalf("got %+v, %v", c, err)
	}
}

func TestBadFilesNameThePath(t *testing.T) {
	for _, body := range []string{"not = [toml", "trust = \"root\"\n"} {
		dir := t.TempDir()
		p := write(t, dir, body)
		_, err := Read(p)
		if err == nil || !strings.Contains(err.Error(), p) {
			t.Fatalf("%q: err %v does not name %s", body, err, p)
		}
	}
}

func TestPathUnknown(t *testing.T) {
	if p := Path(nil, ""); p != "" {
		t.Fatal(p)
	}
	if c, err := Read(""); err != nil || c.Trust != Daily {
		t.Fatal(c, err)
	}
}
