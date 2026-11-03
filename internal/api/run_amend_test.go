package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

const (
	oldBrief = "# Brief\n\nBuild the red dashboard.\n"
	newBrief = "# Brief\n\nBuild the blue dashboard.\n"
)

// amendFixture is a run of agentWorkflow whose build step's session sits
// idle at its prompt, its worktree holding brief.md, and a file holding
// the brief's new content.
type amendFixture struct {
	*sessionFixture
	sessions *session.Store
	sid      string
	run      string
	brief    string
	from     string
}

// runAmendFixture runs the workflow; at the run's first poll, while the
// step's session is idle, during runs with the rex calls so far
// forgotten; at the second poll the session exits 0 and the run succeeds.
func runAmendFixture(t *testing.T, during func(f *amendFixture)) *amendFixture {
	t.Helper()
	sf := newSessionFixture(t, "amend")
	newChange(t, sf.ctx, `{"intent":"X","title":"Amend"}`)
	writeWorkflow(t, sf.ctx.Root, "w", agentWorkflow+"interval = \"1s\"\n")
	sf.ctx.Start = (&shStart{}).start
	f := &amendFixture{
		sessionFixture: sf,
		sessions:       session.NewStore(session.Dir(sf.common)),
		brief:          filepath.Join(sf.worktree, "brief.md"),
		from:           filepath.Join(t.TempDir(), "new-brief.md"),
	}
	if err := os.WriteFile(f.brief, []byte(oldBrief), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.from, []byte(newBrief), 0o644); err != nil {
		t.Fatal(err)
	}
	sf.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		sf.rex.run(name, args, dir, env)
		if len(args) > 0 && args[0] == "window" {
			f.sid = args[len(args)-1]
			if _, err := f.sessions.Update(f.sid, func(r *session.Record) {
				r.Status.Hook = &session.Report{State: session.StateIdle}
			}); err != nil {
				t.Fatal(err)
			}
			return []byte(rexWindowJSON), nil
		}
		return nil, nil
	}
	useClock(t, &scriptClock{t: t, now: time.Unix(0, 0), script: []func(){
		func() {
			runs, err := run.NewStore(run.Dir(sf.common)).List()
			if err != nil || len(runs) != 1 {
				t.Fatalf("runs %v, %v", runs, err)
			}
			f.run = runs[0].ID
			sf.rex.calls, sf.rex.dirs = nil, nil
			during(f)
		},
		func() {
			code := 0
			if _, err := f.sessions.Update(f.sid, func(r *session.Record) {
				r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
			}); err != nil {
				t.Fatal(err)
			}
		},
	}})
	sf.ctx.repoSession = true // rhiza-<repo> is open in the fake Rex
	resp, err := callOp(t, sf.ctx, "run.start", `{"foreground":true,"workflow":"w","permission_mode":"manual"}`)
	if err != nil {
		t.Fatal(err)
	}
	if info := resp.(RunResponse).Run; info.State != run.Succeeded {
		t.Fatalf("run = %+v", info)
	}
	return f
}

// amended is the run's amended journal entries.
func (f *amendFixture) amended(t *testing.T) []run.Entry {
	t.Helper()
	j, err := run.NewStore(run.Dir(f.common)).Read(f.run)
	if err != nil {
		t.Fatal(err)
	}
	var out []run.Entry
	for _, e := range j.Entries {
		if e.Type == run.TypeAmended {
			out = append(out, e)
		}
	}
	return out
}

// unchanged fails unless the brief holds its old content and mode, the
// journal has no amended entry and rex was not run.
func (f *amendFixture) unchanged(t *testing.T, when string) {
	t.Helper()
	data, err := os.ReadFile(f.brief)
	if err != nil || string(data) != oldBrief {
		t.Errorf("%s: brief %q, %v", when, data, err)
	}
	if fi, err := os.Stat(f.brief); err != nil || fi.Mode().Perm() != 0o600 {
		t.Errorf("%s: brief mode %v, %v", when, fi.Mode(), err)
	}
	if a := f.amended(t); len(a) != 0 {
		t.Errorf("%s: journaled %+v", when, a)
	}
	if len(f.rex.calls) != 0 {
		t.Errorf("%s: rex ran %q", when, f.rex.calls)
	}
	if left, _ := filepath.Glob(filepath.Join(f.worktree, ".brief.md.amend-*")); len(left) != 0 {
		t.Errorf("%s: left temporary files %v", when, left)
	}
}

func TestRunAmendIdleStep(t *testing.T) {
	var resp RunAmendResponse
	f := runAmendFixture(t, func(f *amendFixture) {
		out, err := callOp(t, f.ctx, "run.amend", `{"file":"./brief.md","from":"`+f.from+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		resp = out.(RunAmendResponse)
		pointer := "rhiza: test amended brief.md: read it again and carry on from it."
		if resp.Run != f.run || resp.Step != "build" || resp.Session != f.sid || resp.File != "brief.md" ||
			resp.Path != f.brief || resp.Actor != "test" || resp.Pointer != pointer || resp.Waited != "0s" || resp.DryRun {
			t.Errorf("response %+v", resp)
		}
		want := [][]string{
			{"rex", "send", "-b", "block:b1", pointer},
			{"rex", "send-key", "-b", "block:b1", "enter"},
		}
		if !slices.EqualFunc(f.rex.calls, want, slices.Equal) {
			t.Errorf("rex calls %q, want %q", f.rex.calls, want)
		}
		data, err := os.ReadFile(f.brief)
		if err != nil || string(data) != newBrief {
			t.Errorf("brief %q, %v", data, err)
		}
		if fi, err := os.Stat(f.brief); err != nil || fi.Mode().Perm() != 0o600 {
			t.Errorf("brief mode %v, %v; want 0600 kept", fi.Mode(), err)
		}
		a := f.amended(t)
		if len(a) != 1 || a[0].Step != "build" || a[0].Session != f.sid || a[0].Actor != "test" || a[0].File != "brief.md" {
			t.Errorf("amended entries %+v", a)
		}
		st, err := callOp(t, f.ctx, "run.status", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		info := st.(RunInfo)
		if len(info.Amendments) != 1 || info.Amendments[0].Actor != "test" || info.Amendments[0].File != "brief.md" ||
			info.Amendments[0].Session != f.sid || info.Current != "build" {
			t.Errorf("status %+v", info)
		}
		h := human(t, resp)
		if !strings.HasPrefix(h, "Amended brief.md of step build in run "+f.run+"\n") || !strings.Contains(h, "- Pointer: "+pointer+"\n") {
			t.Errorf("human:\n%s", h)
		}
	})
	st, err := callOp(t, f.ctx, "run.status", `{"run":"`+f.run+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	info := st.(RunInfo)
	if info.State != run.Succeeded || len(info.Amendments) != 1 || info.Amendments[0].Step != "build" {
		t.Errorf("final status %+v", info)
	}
	h := human(t, info)
	line := "- " + info.Amendments[0].At + " test amended brief.md (step build)\n"
	if !strings.Contains(h, "\n## Amendments\n\n"+line) {
		t.Errorf("status markdown lacks %q:\n%s", line, h)
	}
}

func TestRunAmendFromStdin(t *testing.T) {
	runAmendFixture(t, func(f *amendFixture) {
		f.ctx.Stdin = strings.NewReader(newBrief)
		if _, err := callOp(t, f.ctx, "run.amend", `{"run":"`+f.run[:6]+`","step":"build","file":"brief.md","from":"-"}`); err != nil {
			t.Fatal(err)
		}
		if data, _ := os.ReadFile(f.brief); string(data) != newBrief {
			t.Errorf("brief %q", data)
		}
		if a := f.amended(t); len(a) != 1 {
			t.Errorf("amended entries %+v", a)
		}
	})
}

func TestRunAmendRefusals(t *testing.T) {
	cases := []struct {
		name, req string
		exit      int
		message   string
	}{
		{"MissingFile", `{"file":"plan.md","from":"FROM"}`, 4, "no file plan.md"},
		{"Escaping", `{"file":"../escape.md","from":"FROM"}`, 2, "outside the run's worktree"},
		{"EscapingInside", `{"file":"docs/../../escape.md","from":"FROM"}`, 2, "outside the run's worktree"},
		{"Absolute", `{"file":"/etc/hosts","from":"FROM"}`, 2, "absolute"},
		{"UnderGit", `{"file":".git/config","from":"FROM"}`, 2, "under .git"},
		{"NoFile", `{"from":"FROM"}`, 2, "needs --file"},
		{"NoFrom", `{"file":"brief.md"}`, 2, "needs --from"},
		{"EmptyContent", `{"file":"brief.md","from":"EMPTY"}`, 2, "blank"},
		{"UnreadableFrom", `{"file":"brief.md","from":"EMPTY.nope"}`, 2, "cannot read"},
		{"NoSession", `{"step":"done","file":"brief.md","from":"FROM"}`, 3, "has no running session"},
		{"BadTimeout", `{"file":"brief.md","from":"FROM","timeout":"soon"}`, 2, "timeout"},
		{"UnknownRun", `{"run":"ZZZZZZ","file":"brief.md","from":"FROM"}`, 4, "no run matches"},
	}
	runAmendFixture(t, func(f *amendFixture) {
		empty := filepath.Join(t.TempDir(), "empty.md")
		if err := os.WriteFile(empty, []byte(" \n\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, c := range cases {
			req := strings.ReplaceAll(strings.ReplaceAll(c.req, "FROM", f.from), "EMPTY", empty)
			_, err := callOp(t, f.ctx, "run.amend", req)
			if ExitCode(err) != c.exit || !strings.Contains(AsError(err).Message, c.message) {
				t.Errorf("%s: err = %v, exit %d; want %d and %q", c.name, err, ExitCode(err), c.exit, c.message)
			}
			f.unchanged(t, c.name)
		}
	})
}

func TestRunAmendSymlinkRefused(t *testing.T) {
	runAmendFixture(t, func(f *amendFixture) {
		outside := filepath.Join(t.TempDir(), "outside.md")
		if err := os.WriteFile(outside, []byte("outside\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(f.worktree, "link.md")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Dir(outside), filepath.Join(f.worktree, "away")); err != nil {
			t.Fatal(err)
		}
		for _, file := range []string{"link.md", "away/outside.md"} {
			_, err := callOp(t, f.ctx, "run.amend", `{"file":"`+file+`","from":"`+f.from+`"}`)
			wantExit(t, err, 2)
			f.unchanged(t, file)
		}
		if data, _ := os.ReadFile(outside); string(data) != "outside\n" {
			t.Errorf("the file outside the worktree changed: %q", data)
		}
	})
}

func TestRunAmendDryRun(t *testing.T) {
	runAmendFixture(t, func(f *amendFixture) {
		f.ctx.DryRun = true
		defer func() { f.ctx.DryRun = false }()
		out, err := callOp(t, f.ctx, "run.amend", `{"file":"brief.md","from":"`+f.from+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		r := out.(RunAmendResponse)
		if !r.DryRun || r.File != "brief.md" || r.Session != f.sid || r.Actor != "test" || !strings.Contains(r.Pointer, "amended brief.md") {
			t.Errorf("response %+v", r)
		}
		if h := human(t, r); !strings.HasPrefix(h, "Would amend brief.md") || !strings.Contains(h, "This was a dry run: nothing was written, journaled or sent.") {
			t.Errorf("human:\n%s", h)
		}
		f.unchanged(t, "dry run")
	})
}

func TestRunAmendFinishedRun(t *testing.T) {
	f := runAmendFixture(t, func(*amendFixture) {})
	f.rex.calls = nil
	_, err := callOp(t, f.ctx, "run.amend", `{"file":"brief.md","from":"`+f.from+`"}`)
	wantExit(t, err, 3)
	if e := AsError(err); !strings.Contains(e.Message, "run "+f.run+" is finished") || e.Next != "rhiza run status "+f.run {
		t.Errorf("error %+v", e)
	}
	f.unchanged(t, "finished run")
}

func TestRunAmendExitedSession(t *testing.T) {
	runAmendFixture(t, func(f *amendFixture) {
		code := 0
		if _, err := f.sessions.Update(f.sid, func(r *session.Record) {
			r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
		}); err != nil {
			t.Fatal(err)
		}
		_, err := callOp(t, f.ctx, "run.amend", `{"file":"brief.md","from":"`+f.from+`"}`)
		wantExit(t, err, 3)
		f.unchanged(t, "exited session")
	})
}

func TestAmendRel(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"brief.md", "brief.md"},
		{" ./docs/plan.md ", "docs/plan.md"},
		{"docs/../brief.md", "brief.md"},
		{".gitignore", ".gitignore"},
		{"a/.git/x", "a/.git/x"},
	} {
		if got, err := amendRel(c.in); err != nil || got != c.want {
			t.Errorf("amendRel(%q) = %q, %v; want %q", c.in, got, err, c.want)
		}
	}
	for _, in := range []string{"", ".", "..", "../x", "a/../../x", "/abs", ".git", ".git/HEAD", "./.git/config"} {
		_, err := amendRel(in)
		wantExit(t, err, 2)
	}
}
