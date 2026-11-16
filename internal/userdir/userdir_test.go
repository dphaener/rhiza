package userdir

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

// TestMatchesOS checks that each folder, read through os.Getenv, is the
// one the os package gives.
func TestMatchesOS(t *testing.T) {
	if want, werr := os.UserCacheDir(); true {
		got, err := Cache(nil)
		if got != want || (err == nil) != (werr == nil) {
			t.Errorf("Cache(nil) = %q, %v; want %q, %v", got, err, want, werr)
		}
	}
	if want, werr := os.UserHomeDir(); true {
		got, err := Home(nil)
		if got != want || (err == nil) != (werr == nil) {
			t.Errorf("Home(nil) = %q, %v; want %q, %v", got, err, want, werr)
		}
	}
	if got, want := Temp(nil), os.TempDir(); got != want {
		t.Errorf("Temp(nil) = %q, want %q", got, want)
	}
}

// TestFromGetenv checks that the folders follow the given environment, not
// the process's.
func TestFromGetenv(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the XDG rules are Linux's")
	}
	g := env(map[string]string{"HOME": "/h", "XDG_CACHE_HOME": "/c", "TMPDIR": "/t"})
	if d, err := Cache(g); err != nil || d != "/c" {
		t.Errorf("Cache = %q, %v; want /c", d, err)
	}
	if d, err := Cache(env(map[string]string{"HOME": "/h"})); err != nil || d != filepath.Join("/h", ".cache") {
		t.Errorf("Cache without XDG_CACHE_HOME = %q, %v; want /h/.cache", d, err)
	}
	if _, err := Cache(env(map[string]string{"XDG_CACHE_HOME": "rel"})); err == nil {
		t.Error("Cache with a relative XDG_CACHE_HOME: no error")
	}
	if _, err := Cache(env(nil)); err == nil {
		t.Error("Cache with no HOME: no error")
	}
	if d, err := Home(g); err != nil || d != "/h" {
		t.Errorf("Home = %q, %v; want /h", d, err)
	}
	if _, err := Home(env(nil)); err == nil {
		t.Error("Home with no HOME: no error")
	}
	if d := Temp(g); d != "/t" {
		t.Errorf("Temp = %q, want /t", d)
	}
	if d := Temp(env(nil)); d != "/tmp" {
		t.Errorf("Temp with no TMPDIR = %q, want /tmp", d)
	}
}

func TestGetenv(t *testing.T) {
	g := Getenv([]string{"A=1", "B=x=y", "A=2", "junk", "=C:=C:\\"})
	if g("A") != "2" || g("B") != "x=y" || g("junk") != "" || g("missing") != "" {
		t.Errorf("Getenv: A=%q B=%q junk=%q", g("A"), g("B"), g("junk"))
	}
	if g("=C:") != "C:\\" {
		t.Errorf("Getenv(=C:) = %q", g("=C:"))
	}
}
