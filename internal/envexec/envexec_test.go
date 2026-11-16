package envexec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLookPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix modes")
	}
	wd := t.TempDir()
	abs := t.TempDir()
	for _, p := range []string{filepath.Join(wd, "rel", "tool"), filepath.Join(abs, "tool"), filepath.Join(wd, "tool"), filepath.Join(abs, "plain")} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o755)
		if filepath.Base(p) == "plain" {
			mode = 0o644
		}
		if err := os.WriteFile(p, []byte("#!/bin/sh\n"), mode); err != nil {
			t.Fatal(err)
		}
	}
	sep := string(os.PathListSeparator)
	for _, c := range []struct {
		name, file, path, want string
	}{
		{"first entry wins", "tool", abs + sep + "rel", filepath.Join(abs, "tool")},
		{"relative path joins wd", "rel/tool", "", filepath.Join(wd, "rel", "tool")},
		{"not executable", "plain", abs, ""},
		{"missing", "tool", t.TempDir(), ""},
		{"no PATH", "tool", "", ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			env := []string{"PATH=/elsewhere", "PATH=" + c.path}
			got, err := LookPath(c.file, env, wd)
			if c.want == "" {
				if !errors.Is(err, exec.ErrNotFound) && !errors.Is(err, os.ErrPermission) {
					t.Errorf("LookPath(%q) = %q, %v; want not found", c.file, got, err)
				}
				return
			}
			if err != nil || got != c.want {
				t.Errorf("LookPath(%q) = %q, %v; want %q", c.file, got, err, c.want)
			}
		})
	}
}

// TestLookPathErrDot checks that a program found only through a '.', empty
// or relative PATH entry is refused with exec.ErrDot, as exec.LookPath
// refuses it, and that Command does not run it; through an absolute entry
// the same program runs.
func TestLookPathErrDot(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	wd := t.TempDir()
	marker := filepath.Join(t.TempDir(), "ran")
	for _, p := range []string{filepath.Join(wd, "git"), filepath.Join(wd, "rel", "git")} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("#!"+sh+"\necho planted >> "+marker+"\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	sep := string(os.PathListSeparator)
	empty := t.TempDir()
	for _, c := range []struct{ name, path string }{
		{"dot entry", "." + sep + empty},
		{"empty entry", sep + empty},
		{"trailing empty entry", empty + sep},
		{"relative entry", "rel" + sep + empty},
	} {
		t.Run(c.name, func(t *testing.T) {
			env := []string{"PATH=" + c.path}
			if got, err := LookPath("git", env, wd); !errors.Is(err, exec.ErrDot) {
				t.Errorf("LookPath(git) = %q, %v; want exec.ErrDot", got, err)
			} else if e, ok := err.(*exec.Error); !ok || e.Name != "git" {
				t.Errorf("LookPath(git) error %#v, want an *exec.Error named git", err)
			}
			if err := Command(nil, env, wd, "git", "status").Run(); !errors.Is(err, exec.ErrDot) {
				t.Errorf("Command(git).Run = %v, want exec.ErrDot", err)
			}
			if _, err := os.Stat(marker); err == nil {
				t.Fatal("the planted ./git ran")
			}
		})
	}
	t.Run("absolute entry", func(t *testing.T) {
		env := []string{"PATH=" + empty + sep + wd}
		if got, err := LookPath("git", env, wd); err != nil || got != filepath.Join(wd, "git") {
			t.Errorf("LookPath(git) = %q, %v; want %q", got, err, filepath.Join(wd, "git"))
		}
		if err := Command(nil, env, wd, "git").Run(); err != nil {
			t.Fatalf("Command(git).Run = %v", err)
		}
		if b, err := os.ReadFile(marker); err != nil || string(b) != "planted\n" {
			t.Errorf("the absolute entry's git did not run (marker %q, %v)", b, err)
		}
	})
}

// TestCommand checks that Command runs the given PATH's program, keeps
// name as Args[0], and carries a miss as the command's error, never
// falling back to the process's PATH.
func TestCommand(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	wd := t.TempDir()
	if err := os.WriteFile(filepath.Join(wd, "hello"), []byte("#!"+sh+"\necho \"$0 hi\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := Command(nil, []string{"PATH=" + wd}, wd, "hello")
	if cmd.Args[0] != "hello" || cmd.Path != filepath.Join(wd, "hello") {
		t.Errorf("Args[0] %q, Path %q", cmd.Args[0], cmd.Path)
	}
	if out, err := cmd.Output(); err != nil || string(out) != filepath.Join(wd, "hello")+" hi\n" {
		t.Errorf("output %q, %v", out, err)
	}
	// sh is on the process's PATH, not the given one.
	if err := Command(nil, []string{"PATH=" + t.TempDir()}, wd, "sh", "-c", "exit 0").Run(); !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("Run = %v, want not found", err)
	}
}
