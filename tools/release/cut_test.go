package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// cutRepo is a fixture repo whose last tag is tag, made at 2026-11-05, with
// the fixture log committed after it: m2 finished and a bug fix landed.
func cutRepo(t *testing.T, tag string) (dir, head string) {
	t.Helper()
	dir = t.TempDir()
	git := func(date string, args ...string) string {
		t.Helper()
		c := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_COMMITTER_DATE="+date, "GIT_AUTHOR_DATE="+date)
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("2026-11-05T00:00:00Z", "init", "-q")
	git("2026-11-05T00:00:00Z", "commit", "-q", "--allow-empty", "-m", "base")
	git("2026-11-05T00:00:00Z", "tag", tag)
	logs, _ := fixture(t)
	if err := os.CopyFS(filepath.Join(dir, ".rhiza", "log"), os.DirFS(logs)); err != nil {
		t.Fatal(err)
	}
	git("2026-11-10T00:00:00Z", "add", ".")
	git("2026-11-10T00:00:00Z", "commit", "-q", "-m", "log")
	return dir, git("2026-11-10T00:00:00Z", "rev-parse", "HEAD")
}

// fakeGit records every git call except the read-only ones it answers from
// tags (the local and remote tag names).
type fakeGit struct {
	local, remote []string
	calls         [][]string
	notes         string
}

func (f *fakeGit) run(args ...string) (string, error) {
	f.calls = append(f.calls, args)
	has := func(ts []string, name string) string {
		for _, x := range ts {
			if x == name {
				return x
			}
		}
		return ""
	}
	switch args[0] {
	case "tag":
		if args[1] == "--list" {
			return has(f.local, args[2]), nil
		}
		if args[1] == "-a" {
			b, err := os.ReadFile(args[len(args)-1])
			f.notes = string(b)
			return "", err
		}
	case "ls-remote":
		return has(f.remote, strings.TrimPrefix(args[len(args)-1], "refs/tags/")), nil
	}
	return "", nil
}

func (f *fakeGit) pushes() (n int) {
	for _, c := range f.calls {
		if c[0] == "push" {
			n++
		}
	}
	return
}

func TestCut(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	runs := filepath.Join(t.TempDir(), "runs.jsonl")
	envFor := func(head, ci string, at time.Time) CheckEnv {
		os.WriteFile(runs, []byte(`{"commit":"`+head+`","at":"`+at.Format(time.RFC3339)+`"}`+"\n"), 0o644)
		return CheckEnv{Now: func() time.Time { return now }, Runs: runs,
			CI:       func(string) (string, error) { return ci, nil },
			Breaking: func(string) (string, []Break, error) { return "", nil, nil }}
	}

	t.Run("pushes once when every gate passes", func(t *testing.T) {
		dir, head := cutRepo(t, "v0.1.0")
		g := &fakeGit{}
		var b strings.Builder
		out, err := Cut(dir, &b, envFor(head, "success", now.Add(-30*time.Hour)), g.run)
		if err != nil || out != Tagged || g.pushes() != 1 {
			t.Fatalf("%v %v pushes=%d\n%s", out, err, g.pushes(), b.String())
		}
		if !strings.HasPrefix(g.notes, "v0.2.0\n\n") || !strings.Contains(g.notes, "Milestone m2") {
			t.Fatalf("tag message %q", g.notes)
		}
		var tagCall []string
		for _, c := range g.calls {
			if c[0] == "tag" && c[1] == "-a" {
				tagCall = c
			}
		}
		if tagCall[2] != "v0.2.0" || tagCall[3] != head {
			t.Fatalf("tag call %v", tagCall)
		}
		last := g.calls[len(g.calls)-1]
		if last[0] != "push" || last[len(last)-1] != "refs/tags/v0.2.0" {
			t.Fatalf("last call %v", last)
		}
	})

	t.Run("never twice for one version", func(t *testing.T) {
		dir, head := cutRepo(t, "v0.1.0")
		for _, g := range []*fakeGit{{local: []string{"v0.2.0"}}, {remote: []string{"v0.2.0"}}} {
			out, err := Cut(dir, &strings.Builder{}, envFor(head, "success", now.Add(-30*time.Hour)), g.run)
			if err != nil || out != Existing || g.pushes() != 0 {
				t.Fatalf("%v %v pushes=%d", out, err, g.pushes())
			}
		}
	})

	t.Run("never when a gate fails", func(t *testing.T) {
		dir, head := cutRepo(t, "v0.1.0")
		for name, env := range map[string]CheckEnv{
			"ci red":     envFor(head, "failure", now.Add(-30*time.Hour)),
			"soak":       envFor(head, "success", now.Add(-2*time.Hour)),
			"not in run": envFor("other", "success", now.Add(-30*time.Hour)),
		} {
			g := &fakeGit{}
			out, err := Cut(dir, &strings.Builder{}, env, g.run)
			if err != nil || out != Blocked || g.pushes() != 0 {
				t.Fatalf("%s: %v %v pushes=%d", name, out, err, g.pushes())
			}
			for _, c := range g.calls {
				if c[0] == "tag" && c[1] == "-a" {
					t.Fatalf("%s: tagged", name)
				}
			}
		}
	})

	t.Run("never for a major", func(t *testing.T) {
		dir, head := cutRepo(t, "v1.0.0")
		writeEntry(t, filepath.Join(dir, ".rhiza", "log"), "CS5", "zz", note("NB", "breaking: drops --old")+","+discusses("LB", "NB", "CS5"))
		g := &fakeGit{}
		var b strings.Builder
		out, err := Cut(dir, &b, envFor(head, "success", now.Add(-30*time.Hour)), g.run)
		if err != nil || out != SignOff || g.pushes() != 0 || !strings.Contains(b.String(), "v2.0.0") {
			t.Fatalf("%v %v pushes=%d\n%s", out, err, g.pushes(), b.String())
		}
	})
}

func TestRecordDaemonRun(t *testing.T) {
	p := filepath.Join(t.TempDir(), "runs.jsonl")
	t0 := time.Date(2026, 11, 12, 8, 0, 0, 0, time.UTC)
	rec := func(commit string, at time.Time) bool {
		t.Helper()
		added, err := RecordDaemonRun(p, func() (DaemonRun, error) { return DaemonRun{commit, at}, nil })
		if err != nil {
			t.Fatal(err)
		}
		return added
	}
	if !rec("abc", t0) {
		t.Fatal("new commit not recorded")
	}
	// The same commit, even after a daemon restart, adds nothing.
	if rec("abc", t0.Add(3*time.Hour)) {
		t.Fatal("known commit recorded again")
	}
	if !rec("def", t0.Add(4*time.Hour)) {
		t.Fatal("second commit not recorded")
	}
	b, _ := os.ReadFile(p)
	if string(b) != `{"commit":"abc","at":"2026-11-12T08:00:00Z"}`+"\n"+`{"commit":"def","at":"2026-11-12T12:00:00Z"}`+"\n" {
		t.Fatalf("history %q", b)
	}
	h, err := ReadRunHistory(p)
	if err != nil || !h["abc"].Equal(t0) {
		t.Fatalf("%v %v", h, err)
	}
	if _, err := RecordDaemonRun("", nil); err == nil {
		t.Fatal("empty path accepted")
	}
}

func TestParseDaemonStatus(t *testing.T) {
	d, err := parseDaemonStatus([]byte(`{"running":true,"commit":"abc","started":"2026-11-12T20:16:37Z"}`))
	if err != nil || d.Commit != "abc" || !d.Started.Equal(time.Date(2026, 11, 12, 20, 16, 37, 0, time.UTC)) {
		t.Fatalf("%+v %v", d, err)
	}
	for _, bad := range []string{`{"running":false,"commit":"abc","started":"2026-11-12T20:16:37Z"}`, `{"running":true,"started":"2026-11-12T20:16:37Z"}`, `{"running":true,"commit":"a","started":""}`} {
		if _, err := parseDaemonStatus([]byte(bad)); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
}

// The release job publishes the tag's annotation as the release body, and
// the timer unit in the repo is the one the VM runbook names.
func TestReleaseWorkflowAndTimer(t *testing.T) {
	rd := func(p string) string {
		b, err := os.ReadFile(filepath.Join("..", "..", p))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	wf := rd(".github/workflows/release.yml")
	for _, want := range []string{"contents:body", "--release-notes"} {
		if !strings.Contains(wf, want) {
			t.Errorf("release.yml lacks %q", want)
		}
	}
	timer := rd("script/systemd/rhiza-release.timer")
	if !strings.Contains(timer, "OnCalendar=hourly") {
		t.Error("timer is not hourly")
	}
	svc := rd("script/systemd/rhiza-release.service")
	for _, want := range []string{"record-run", "tools/release cut", "RHIZA_RELEASE_RUNS"} {
		if !strings.Contains(svc, want) {
			t.Errorf("service lacks %q", want)
		}
	}
	if !strings.Contains(rd("agents/vm.md"), "rhiza-release.timer") {
		t.Error("agents/vm.md does not name rhiza-release.timer")
	}
}
