package main

import (
	"fmt"
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

// fakeGit records every git call and answers the read-only ones from its
// tags: local and remote hold annotated tag names, published the lightweight
// tags (name to commit) that publishing a GitHub release makes on the
// remote; fetch brings those here. A tag -a makes a local annotated tag and
// a push moves it to the remote.
type fakeGit struct {
	local, remote []string
	published     map[string]string
	fetched       map[string]string
	calls         [][]string
	notes         string
}

func (f *fakeGit) run(args ...string) (string, error) {
	f.calls = append(f.calls, args)
	has := func(ts []string, name string) bool {
		for _, x := range ts {
			if x == name {
				return true
			}
		}
		return false
	}
	ref := func(s string) string {
		return strings.TrimSuffix(strings.TrimPrefix(s, "refs/tags/"), "^{commit}")
	}
	switch args[0] {
	case "fetch":
		for name, c := range f.published {
			if !has(f.local, name) {
				if f.fetched == nil {
					f.fetched = map[string]string{}
				}
				f.fetched[name] = c
			}
		}
	case "tag":
		if args[1] == "--list" {
			if _, ok := f.fetched[args[2]]; ok || has(f.local, args[2]) {
				return args[2], nil
			}
			return "", nil
		}
		if args[1] == "-a" {
			name := args[2]
			if name == "-f" {
				name = args[3]
			} else if _, ok := f.fetched[name]; ok {
				return "", fmt.Errorf("tag %s already exists", name)
			}
			b, err := os.ReadFile(args[len(args)-1])
			f.notes = string(b)
			delete(f.fetched, name)
			f.local = append(f.local, name)
			return "", err
		}
	case "cat-file":
		if has(f.local, ref(args[2])) {
			return "tag", nil
		}
		return "commit", nil
	case "rev-parse":
		return f.fetched[ref(args[1])], nil
	case "ls-remote":
		name := ref(args[3])
		if has(f.remote, name) {
			return "obj\trefs/tags/" + name + "\nc\trefs/tags/" + name + "^{}", nil
		}
		if c, ok := f.published[name]; ok {
			return c + "\trefs/tags/" + name, nil
		}
	case "push":
		name := ref(args[len(args)-1])
		if _, ok := f.published[name]; ok && args[1] != "--force" {
			return "", fmt.Errorf("push %s: rejected, already exists", name)
		}
		delete(f.published, name)
		f.remote = append(f.remote, name)
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
}

func TestMajor(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	runs := filepath.Join(t.TempDir(), "runs.jsonl")
	dir, head := cutRepo(t, "v1.0.0")
	writeEntry(t, filepath.Join(dir, ".rhiza", "log"), "CS5", "zz", note("NB", "breaking: drops --old")+","+discusses("LB", "NB", "CS5"))
	os.WriteFile(runs, []byte(`{"commit":"`+head+`","at":"`+now.Add(-30*time.Hour).Format(time.RFC3339)+`"}`+"\n"), 0o644)
	state := ReleaseNone
	var drafts []string
	env := CheckEnv{Now: func() time.Time { return now }, Runs: runs,
		CI:       func(string) (string, error) { return "success", nil },
		Breaking: func(string) (string, []Break, error) { return "", nil, nil },
		Release:  func(string) (string, error) { return state, nil },
		Draft: func(tag, commit, notes string) error {
			drafts = append(drafts, tag+" "+commit)
			// The release body, as GoReleaser keeps it: the notes alone.
			if notes == "" || strings.HasPrefix(notes, "v2.0.0") {
				t.Errorf("draft notes %q", notes)
			}
			state = ReleaseDraft
			return nil
		}}
	g := &fakeGit{}
	cut := func() (Outcome, string) {
		t.Helper()
		var b strings.Builder
		out, err := Cut(dir, &b, env, g.run)
		if err != nil {
			t.Fatal(err)
		}
		return out, b.String()
	}

	// Held, with one draft prepared, until the owner publishes it.
	for i := 0; i < 2; i++ {
		if out, msg := cut(); out != SignOff || g.pushes() != 0 || !strings.Contains(msg, "v2.0.0") {
			t.Fatalf("%v pushes=%d\n%s", out, g.pushes(), msg)
		}
	}
	if len(drafts) != 1 || drafts[0] != "v2.0.0 "+head {
		t.Fatalf("drafts %v", drafts)
	}
	// Signed off: publishing the draft makes a lightweight tag at the
	// candidate on the remote, which the tool replaces with the annotated
	// one, pushed once.
	state = ReleasePublished
	g.published = map[string]string{"v2.0.0": head}
	if out, msg := cut(); out != Tagged || g.pushes() != 1 {
		t.Fatalf("%v pushes=%d\n%s", out, g.pushes(), msg)
	}
	if !strings.HasPrefix(g.notes, "v2.0.0\n\n") {
		t.Fatalf("tag message %q", g.notes)
	}
	// Once the annotated tag is there, a repeat does nothing.
	if out, _ := cut(); out != Existing || g.pushes() != 1 {
		t.Fatalf("%v pushes=%d", out, g.pushes())
	}
	// Signed off but no tag yet on the remote: a plain tag and push.
	g = &fakeGit{}
	if out, _ := cut(); out != Tagged || g.pushes() != 1 || g.calls[len(g.calls)-1][1] == "--force" {
		t.Fatalf("%v calls %v", out, g.calls)
	}
	// A published release tagged at another commit is never moved.
	g = &fakeGit{published: map[string]string{"v2.0.0": "other"}}
	var b strings.Builder
	if _, err := Cut(dir, &b, env, g.run); err == nil || g.pushes() != 0 {
		t.Fatalf("err %v pushes=%d", err, g.pushes())
	}
	// A failing gate holds a major too, with no draft.
	env.CI = func(string) (string, error) { return "failure", nil }
	g, state, drafts = &fakeGit{}, ReleaseNone, nil
	if out, _ := cut(); out != Blocked || len(drafts) != 0 {
		t.Fatalf("%v drafts %v", out, drafts)
	}
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

// TestCutRealGitKeepsHeadings tags with real git, plainly and by replacing
// a published lightweight tag, and reads the annotation back: the notes'
// "## " section headings survive, as release.yml publishes them.
func TestCutRealGitKeepsHeadings(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	gitIn := func(dir string) GitFunc {
		return func(args ...string) (string, error) {
			c := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "tag.gpgsign=false"}, args...)...)
			c.Dir = dir
			out, err := c.CombinedOutput()
			s := strings.TrimSpace(string(out))
			if err != nil {
				return s, fmt.Errorf("git %v: %w: %s", args, err, s)
			}
			return s, nil
		}
	}
	setup := func(t *testing.T, tag string) (dir, head string, git GitFunc, env CheckEnv) {
		dir, head = cutRepo(t, tag)
		git = gitIn(dir)
		origin := t.TempDir()
		if _, err := gitIn(origin)("init", "-q", "--bare"); err != nil {
			t.Fatal(err)
		}
		if _, err := git("remote", "add", "origin", origin); err != nil {
			t.Fatal(err)
		}
		runs := filepath.Join(t.TempDir(), "runs.jsonl")
		os.WriteFile(runs, []byte(`{"commit":"`+head+`","at":"`+now.Add(-30*time.Hour).Format(time.RFC3339)+`"}`+"\n"), 0o644)
		env = CheckEnv{Now: func() time.Time { return now }, Runs: runs,
			CI:       func(string) (string, error) { return "success", nil },
			Breaking: func(string) (string, []Break, error) { return "", nil, nil }}
		return
	}
	body := func(t *testing.T, git GitFunc, v string) string {
		t.Helper()
		out, err := git("for-each-ref", "--format=%(contents:body)", "refs/tags/"+v)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "## ") {
			t.Fatalf("%s body lost its headings:\n%s", v, out)
		}
		return out
	}

	t.Run("plain tag", func(t *testing.T) {
		dir, _, git, env := setup(t, "v0.1.0")
		var b strings.Builder
		if out, err := Cut(dir, &b, env, git); err != nil || out != Tagged {
			t.Fatalf("%v %v\n%s", out, err, b.String())
		}
		body(t, git, "v0.2.0")
	})

	t.Run("replacing a published lightweight tag", func(t *testing.T) {
		dir, head, git, env := setup(t, "v1.0.0")
		writeEntry(t, filepath.Join(dir, ".rhiza", "log"), "CS5", "zz", note("NB", "breaking: drops --old")+","+discusses("LB", "NB", "CS5"))
		env.Release = func(string) (string, error) { return ReleasePublished, nil }
		env.Draft = func(string, string, string) error { t.Fatal("drafted a signed-off release"); return nil }
		// Publishing the release made a lightweight tag at the candidate on
		// the remote.
		for _, args := range [][]string{{"tag", "v2.0.0", head}, {"push", "-q", "origin", "refs/tags/v2.0.0"}, {"tag", "-d", "v2.0.0"}} {
			if _, err := git(args...); err != nil {
				t.Fatal(err)
			}
		}
		// The lightweight tag stays on the remote only (no fetch), so this
		// exercises the tag -a -f replace and its forced push with real git.
		noFetch := func(args ...string) (string, error) {
			if args[0] == "fetch" {
				return "", nil
			}
			return git(args...)
		}
		var b strings.Builder
		if out, err := Cut(dir, &b, env, noFetch); err != nil || out != Tagged {
			t.Fatalf("%v %v\n%s", out, err, b.String())
		}
		if !strings.Contains(body(t, git, "v2.0.0"), "## Breaking changes") {
			t.Fatal("breaking changes heading missing")
		}
	})
}
