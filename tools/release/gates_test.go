package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGates(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	// testdata/runs.jsonl: abc first ran 54h ago (its earliest line), fresh
	// 2h ago; zzz is not in it.
	hist, err := ReadRunHistory(filepath.Join("testdata", "runs.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !hist["abc"].Equal(now.Add(-54*time.Hour)) || len(hist) != 2 {
		t.Fatalf("history %v", hist)
	}
	base := GateInputs{Now: now, Commit: "abc", Kind: KindPatch, CI: "success", Runs: hist}
	status := func(gs []Gate, name string) Gate {
		for _, g := range gs {
			if g.Name == name {
				return g
			}
		}
		t.Fatalf("no gate %s", name)
		return Gate{}
	}
	mod := func(f func(*GateInputs)) GateInputs { in := base; f(&in); return in }
	for _, c := range []struct {
		name string
		in   GateInputs
		gate string
		pass bool
	}{
		{"all pass", base, "", true},
		{"ci red", mod(func(i *GateInputs) { i.CI = "failure" }), "ci", false},
		{"ci missing", mod(func(i *GateInputs) { i.CI = "" }), "ci", false},
		{"commit missing from history", mod(func(i *GateInputs) { i.Commit = "zzz" }), "soak", false},
		{"soak too short", mod(func(i *GateInputs) { i.Commit = "fresh" }), "soak", false},
		{"bug filed since", mod(func(i *GateInputs) { i.BugRequests = []time.Time{now.Add(-3 * time.Hour)} }), "soak", false},
		{"bug before the soak", mod(func(i *GateInputs) { i.BugRequests = []time.Time{now.Add(-60 * time.Hour)} }), "soak", true},
		{"milestone in flight", mod(func(i *GateInputs) { i.InFlight = []string{"CS1"} }), "milestone", false},
		{"milestone blocked", mod(func(i *GateInputs) { i.Blocking = []string{"Q1"} }), "milestone", false},
		{"second patch in 24h", mod(func(i *GateInputs) { i.LastPatch = now.Add(-5 * time.Hour) }), "patch", false},
		{"patch after 24h", mod(func(i *GateInputs) { i.LastPatch = now.Add(-25 * time.Hour) }), "patch", true},
		{"severe skips soak", mod(func(i *GateInputs) { i.Severe, i.Commit = true, "zzz" }), "soak", true},
		{"severe skips batching", mod(func(i *GateInputs) { i.Severe, i.LastPatch = true, now.Add(-time.Hour) }), "patch", true},
		{"severe never skips ci", mod(func(i *GateInputs) { i.Severe, i.CI = true, "failure" }), "ci", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			gs := Gates(c.in)
			if c.gate == "" {
				if !AllPass(gs) {
					t.Fatalf("%+v", gs)
				}
				return
			}
			g := status(gs, c.gate)
			if g.Pass != c.pass || g.Reason == "" {
				t.Fatalf("%+v", g)
			}
			if AllPass(gs) != (c.pass && AllPass(gs)) {
				t.Fatalf("AllPass disagrees: %+v", gs)
			}
		})
	}
}

func TestReadRunHistory(t *testing.T) {
	p := filepath.Join(t.TempDir(), "runs.jsonl")
	os.WriteFile(p, []byte(`{"commit":"a","at":"2026-11-10T00:00:00Z"}`+"\n\n"+`{"commit":"a","at":"2026-11-09T00:00:00Z"}`+"\n"), 0o644)
	h, err := ReadRunHistory(p)
	if err != nil || !h["a"].Equal(time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("%v %v", h, err)
	}
	if h, err := ReadRunHistory(p + ".none"); err != nil || len(h) != 0 {
		t.Fatalf("%v %v", h, err)
	}
	os.WriteFile(p, []byte("nope\n"), 0o644)
	if _, err := ReadRunHistory(p); err == nil {
		t.Fatal("bad line accepted")
	}
}

// TestGatesCheckReports runs tools/release check over a fixture repo with a
// fake clock, a fake CI and the fixture run history: each gate is printed
// with why it passes or fails.
func TestGatesCheckReports(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		c := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
		c.Dir = dir
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	logs, _ := fixture(t)
	if err := os.CopyFS(filepath.Join(dir, ".rhiza", "log"), os.DirFS(logs)); err != nil {
		t.Fatal(err)
	}
	git("init", "-q")
	git("add", ".")
	git("commit", "-q", "-m", "one")
	head := git("rev-parse", "HEAD")
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	runs := filepath.Join(t.TempDir(), "runs.jsonl")
	write := func(at time.Time) {
		os.WriteFile(runs, []byte(`{"commit":"`+head+`","at":"`+at.Format(time.RFC3339)+`"}`+"\n"), 0o644)
	}
	check := func(ci string) string {
		t.Helper()
		var asked string
		env := CheckEnv{Now: func() time.Time { return now }, Runs: runs,
			CI: func(c string) (string, error) { asked = c; return ci, nil }}
		var b strings.Builder
		if err := runCheck(dir, &b, env); err != nil {
			t.Fatal(err)
		}
		if asked != head {
			t.Fatalf("CI asked for %q, want %s", asked, head)
		}
		return b.String()
	}
	want := func(out string, lines ...string) {
		t.Helper()
		for _, l := range lines {
			if !strings.Contains(out, l) {
				t.Fatalf("missing %q in:\n%s", l, out)
			}
		}
	}
	// The bug request (CS2) was opened 2026-11-10; a soak from 11-11 is clean.
	write(now.Add(-30 * time.Hour))
	out := check("success")
	want(out, "gates for "+head+": ready", "ci        pass", "soak      pass", "milestone pass", "patch     pass  not a patch release")

	// CI red and a soak begun before the bug request.
	write(now.Add(-70 * time.Hour))
	out = check("failure")
	want(out, ": blocked", "ci        FAIL  CI on "+head+" is failure", "soak      FAIL  a bug request was filed at 2026-11-10T00:00:00Z")

	// Missing from the history.
	os.WriteFile(runs, nil, 0o644)
	want(check("success"), "soak      FAIL  "+head+" is not in the daemon run history")

	// An open question blocking a node CS4 (serving the finished m2) created.
	write(now.Add(-30 * time.Hour))
	q := filepath.Join(dir, ".rhiza", "log")
	writeAt(t, q, "CS4", "c", "ops", "2026-11-10T00:00:00Z", `{"ops":[{"op":"node.create","kind":"capability","node":"N4","handle":"leaf4"}]}`)
	writeAt(t, q, "CS7", "a", "ops", "2026-11-11T00:00:00Z", `{"ops":[{"op":"node.create","kind":"question","node":"Q1","handle":"which-way"},
		{"op":"link.create","type":"blocks","link":"L1","from":"Q1","to":"N4"}]}`)
	want(check("success"), "milestone FAIL  questions blocking the milestone: which-way")

	// Answered, it no longer blocks.
	writeAt(t, q, "CS7", "b", "ops", "2026-11-11T01:00:00Z", `{"ops":[{"op":"node.create","kind":"decision","node":"D1","handle":"this-way"},
		{"op":"link.create","type":"answers","link":"L2","from":"D1","to":"Q1"}]}`)
	want(check("success"), "milestone pass")
}

func TestParseCheck(t *testing.T) {
	t.Setenv("RHIZA_RELEASE_RUNS", "/x/runs.jsonl")
	env, err := parseCheck([]string{"-commit", "abc"})
	if err != nil || env.Runs != "/x/runs.jsonl" || env.Commit != "abc" || env.CI == nil || env.Now == nil {
		t.Fatalf("%+v %v", env, err)
	}
}
