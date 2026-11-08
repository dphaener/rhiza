package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
)

// fakeExec records the programs a call runs and answers with out, or with
// err when set.
type fakeExec struct {
	calls [][]string
	dirs  []string
	out   string
	err   error
}

func (f *fakeExec) run(name string, args []string, dir string, _ []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{name}, args...))
	f.dirs = append(f.dirs, dir)
	if f.err != nil {
		return nil, f.err
	}
	return []byte(f.out), nil
}

const rexWindowJSON = `{"window_id":"window:w1","block_ids":["block:b1","block:b2"]}`

// sessionFixture is a launcher repository at <tmp>/repo, its fake git and
// rex, and a worktree for slug at <tmp>/repo.<slug> that already exists
// on cs/<slug> and holds a Rhiza repository (so the launch reuses it).
type sessionFixture struct {
	ctx      *Ctx
	git      *fakeGit
	rex      *fakeExec
	common   string
	worktree string
}

func newSessionFixture(t *testing.T, slug string) *sessionFixture {
	t.Helper()
	tmp := t.TempDir()
	root := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	// tmp is a Rhiza root registering repo, as the testscript Setup does,
	// so runs and plan requests (which need a registered project) start.
	if err := os.WriteFile(filepath.Join(tmp, "rhiza.toml"), []byte("[projects.repo]\npath = \"repo\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	f := &sessionFixture{ctx: ctx, common: filepath.Join(root, ".git"), worktree: filepath.Join(tmp, "repo."+slug)}
	if err := os.MkdirAll(f.worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	wctx := testCtx(f.worktree)
	if _, err := callOp(t, wctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	f.git = &fakeGit{out: map[string]string{
		"rev-parse --git-common-dir": ".git\n",
		"rev-parse --show-prefix":    "\n",
		"worktree list --porcelain": "worktree " + root + "\nHEAD 1111\nbranch refs/heads/main\n\n" +
			"worktree " + f.worktree + "\nHEAD 2222\nbranch refs/heads/cs/" + slug + "\n\n",
	}}
	f.rex = &fakeExec{out: rexWindowJSON}
	ctx.Git = f.git.run
	ctx.Exec = f.rex.run
	ctx.Self = "/opt/bin/rhiza"
	return f
}

func startSession(t *testing.T, ctx *Ctx, req string) SessionStartResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "session.start", req)
	if err != nil {
		t.Fatalf("session.start %s: %v", req, err)
	}
	return resp.(SessionStartResponse)
}

func TestSessionStartRefusals(t *testing.T) {
	cases := []struct {
		name, req string
		exit      int
		message   string
	}{
		{"NoMode", `{}`, 2, "--permission-mode"},
		{"UnknownMode", `{"permission_mode":"yolo"}`, 2, "unknown permission mode"},
		{"Bypass", `{"permission_mode":"bypassPermissions"}`, 1, `trust = "dedicated"`},
		{"UnknownAgent", `{"permission_mode":"manual","agent":"nobody"}`, 2, "unknown agent"},
		{"BadStep", `{"permission_mode":"manual","step":"Review Now"}`, 2, "step"},
		{"BadRun", `{"permission_mode":"manual","run":"run 1"}`, 2, "run"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newSessionFixture(t, "x")
			newChange(t, f.ctx, `{"intent":"X"}`)
			_, err := callOp(t, f.ctx, "session.start", c.req)
			wantExit(t, err, c.exit)
			if !strings.Contains(AsError(err).Message, c.message) {
				t.Errorf("message %q lacks %q", AsError(err).Message, c.message)
			}
			if c.exit == 2 && AsError(err).Next != "rhiza describe session.start" {
				t.Errorf("next = %q", AsError(err).Next)
			}
			if len(f.rex.calls) != 0 || len(f.git.calls) != 0 {
				t.Errorf("a refused launch ran rex %v or git %v", f.rex.calls, f.git.calls)
			}
		})
	}
}

func TestSessionStartNoChangeSet(t *testing.T) {
	f := newSessionFixture(t, "x")
	_, err := callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 2)
	_, err = callOp(t, f.ctx, "session.start", `{"permission_mode":"manual","change":"ZZZZZZ"}`)
	wantExit(t, err, 4)
}

func TestSessionStartReusesWorktree(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	cs := newChange(t, f.ctx, `{"intent":"Make tabs readable: name them","title":"Readable tabs"}`).ChangeSet
	r := startSession(t, f.ctx, `{"permission_mode":"manual"}`)
	s := r.Session

	if !r.WorktreeReused || s.Worktree != f.worktree || s.Branch != "cs/readable-tabs" {
		t.Errorf("worktree %s on %s, reused %v", s.Worktree, s.Branch, r.WorktreeReused)
	}
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(s.ID) {
		t.Errorf("session ID %q is not a UUID v4", s.ID)
	}
	if s.Change != cs.String() || s.Title != "Readable tabs" || s.Step != "build" || s.PermissionMode != "manual" {
		t.Errorf("session %+v", s)
	}
	if s.Tab != (SessionTab{Window: "window:w1", Block: "block:b1", Title: "repo: Readable tabs"}) {
		t.Errorf("tab %+v", s.Tab)
	}
	if s.Status.State != session.StateStarting || s.Source != "default" {
		t.Errorf("status %+v from %s", s.Status, s.Source)
	}

	// Exactly one rex call, in the worktree, opening the tab there.
	want := []string{"rex", "window", "new", "repo: Readable tabs", "-c", f.worktree, "--json", "--",
		"/opt/bin/rhiza", "session", "run", s.ID}
	if len(f.rex.calls) != 1 || !slices.Equal(f.rex.calls[0], want) || f.rex.dirs[0] != f.worktree {
		t.Fatalf("rex calls %q in %q, want %q", f.rex.calls, f.rex.dirs, want)
	}
	if !slices.Equal(r.RexArgv, want) {
		t.Errorf("rex_argv %q", r.RexArgv)
	}
	if slices.ContainsFunc(f.git.calls, func(c string) bool { return strings.HasPrefix(c, "worktree add") }) {
		t.Errorf("a reused worktree was added again: %v", f.git.calls)
	}

	// The worktree's log got the change set, which is current there, and
	// the agent's actor, written by the launcher.
	ww := ctxWorkspace(t, testCtx(f.worktree))
	if cur, err := ww.current(); err != nil || cur != cs {
		t.Errorf("worktree current = %s, %v; want %s", cur, err, cs)
	}
	if _, ok := ww.changeSet(cs); !ok {
		t.Error("the change set was not copied into the worktree's log")
	}
	if s.Actor != s.Agent+"-"+s.ID[:8] {
		t.Errorf("actor %q is not <agent>-<the session ID's first 8 hex>", s.Actor)
	}
	m := ww.model(id.ID{})
	n, err := m.Resolve(s.Actor)
	if err != nil {
		t.Fatalf("actor %s not in the worktree's model: %v", s.Actor, err)
	}
	if n.Kind != "actor" {
		t.Errorf("actor node kind %s", n.Kind)
	}
	launcher := ctxWorkspace(t, f.ctx)
	if _, err := launcher.model(id.ID{}).Resolve(s.Actor); err == nil {
		t.Error("the agent's actor was recorded in the launcher's model too")
	}
	if s.For == "" {
		t.Error("no for")
	}

	// The record is in the git common dir and resolves to starting.
	rec, err := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Tab.Window != "window:w1" || rec.Actor != s.Actor || rec.PermissionMode != "manual" {
		t.Errorf("record %+v", rec)
	}
}

func TestSessionStartStepAndRexSession(t *testing.T) {
	f := newSessionFixture(t, "carving-keeps-ids")
	f.ctx.RexSession = "session:mine"
	newChange(t, f.ctx, `{"intent":"Carving keeps IDs, always"}`)
	r := startSession(t, f.ctx, `{"permission_mode":"plan","step":"review","run":"run-1"}`)
	if r.Session.Tab.Title != "repo: Carving keeps IDs · run run-1 · review" || r.Session.Run != "run-1" {
		t.Errorf("session %+v", r.Session)
	}
	if got := f.rex.calls[0][3:6]; !slices.Equal(got, []string{"-s", "session:mine", "repo: Carving keeps IDs · run run-1 · review"}) {
		t.Errorf("rex args %q", f.rex.calls[0])
	}
	if r.Session.Tab.Session != "session:mine" {
		t.Errorf("tab session %q", r.Session.Tab.Session)
	}
	r = startSession(t, f.ctx, `{"permission_mode":"plan","rex_session":"session:other"}`)
	if !slices.Contains(f.rex.calls[1], "session:other") || r.Session.Tab.Session != "session:other" {
		t.Errorf("rex args %q", f.rex.calls[1])
	}
}

func TestSessionStartDryRun(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	before := logFiles(t, f.worktree)
	f.ctx.DryRun = true
	r := startSession(t, f.ctx, `{"permission_mode":"manual"}`)
	if !r.DryRun || len(r.RexArgv) == 0 || r.Session.Tab.Title != "repo: Readable tabs" {
		t.Errorf("dry run %+v", r)
	}
	if len(f.rex.calls) != 0 {
		t.Errorf("a dry run ran rex: %v", f.rex.calls)
	}
	if after := logFiles(t, f.worktree); !slices.Equal(before, after) {
		t.Errorf("a dry run wrote the worktree's log: %v", after)
	}
	if _, err := os.Stat(session.Dir(f.common)); err == nil {
		t.Error("a dry run wrote a session")
	}
	if !strings.Contains(human(t, r), "This was a dry run") {
		t.Errorf("human %q", human(t, r))
	}
}

func TestSessionStartCreatesWorktree(t *testing.T) {
	f := newSessionFixture(t, "other")
	newChange(t, f.ctx, `{"intent":"Fresh work"}`)
	path := filepath.Join(filepath.Dir(f.worktree), "repo.fresh-work")
	// git worktree add is faked: it makes the folder a Rhiza repository.
	f.git.out["cat-file -e HEAD:.rhiza/repo.json"] = ""
	add := "worktree add -b cs/fresh-work " + path
	git := f.git.run
	f.ctx.Git = func(args ...string) (string, error) {
		if strings.Join(args, " ") == add {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return "", err
			}
			if _, err := callOp(t, testCtx(path), "init", `{}`); err != nil {
				return "", err
			}
			f.git.calls = append(f.git.calls, add)
			return "", nil
		}
		return git(args...)
	}
	r := startSession(t, f.ctx, `{"permission_mode":"manual"}`)
	if r.WorktreeReused || r.Session.Worktree != path || !slices.Contains(f.git.calls, add) {
		t.Errorf("worktree %s reused %v, git %v", r.Session.Worktree, r.WorktreeReused, f.git.calls)
	}
}

func TestSessionStartNoRepoAtHEAD(t *testing.T) {
	f := newSessionFixture(t, "other")
	newChange(t, f.ctx, `{"intent":"Fresh work"}`)
	_, err := callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 3)
	if len(f.rex.calls) != 0 {
		t.Errorf("rex ran: %v", f.rex.calls)
	}
}

func TestSessionStartSquattedPath(t *testing.T) {
	f := newSessionFixture(t, "other")
	newChange(t, f.ctx, `{"intent":"Squatted"}`)
	if err := os.MkdirAll(filepath.Join(filepath.Dir(f.worktree), "repo.squatted"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 3)

	// The worktree on another branch is a conflict too.
	g := newSessionFixture(t, "readable-tabs")
	newChange(t, g.ctx, `{"intent":"Readable tabs"}`)
	g.git.out["worktree list --porcelain"] = "worktree " + g.worktree + "\nHEAD 2222\nbranch refs/heads/elsewhere\n\n"
	_, err = callOp(t, g.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 3)
	if !strings.Contains(AsError(err).Message, "elsewhere") {
		t.Errorf("message %q", AsError(err).Message)
	}
}

func TestSessionStartRexFails(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	f.rex.err = errors.New("rex window new: rex is not running")
	_, err := callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 1)
	if !strings.Contains(AsError(err).Message, "rex is not running") {
		t.Errorf("message %q", AsError(err).Message)
	}
	recs, lerr := session.NewStore(session.Dir(f.common)).List()
	if lerr != nil || len(recs) != 1 {
		t.Fatalf("sessions %v, %v", recs, lerr)
	}
	rep, src := session.Resolve(recs[0])
	if rep.State != session.StateError || src != session.SourceExit || !strings.Contains(rep.Detail, "not running") {
		t.Errorf("status %+v from %s", rep, src)
	}

	// Unreadable JSON from rex fails the same way.
	f.rex.err, f.rex.out = nil, "nope"
	_, err = callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 1)
}

func TestSessionLsShow(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	resp, err := callOp(t, f.ctx, "session.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if ls := resp.(SessionLsResponse); len(ls.Sessions) != 0 || !strings.Contains(human(t, ls), "No sessions yet") {
		t.Errorf("empty ls %+v", ls)
	}

	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	s := startSession(t, f.ctx, `{"permission_mode":"manual"}`).Session
	// Reports from a live session's sources: the hook wins over Rex.
	if _, err := session.NewStore(session.Dir(f.common)).Update(s.ID, func(r *session.Record) {
		r.Status.Rex = &session.Report{State: session.StateWorking, At: r.StartedAt}
		r.Status.Hook = &session.Report{State: session.StateBlocked, Kind: session.KindPermission, Detail: "Bash", At: r.StartedAt}
	}); err != nil {
		t.Fatal(err)
	}
	// ls and show never run anything.
	f.rex.err = errors.New("rex must not be called")

	resp, err = callOp(t, f.ctx, "session.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	ls := resp.(SessionLsResponse)
	if len(ls.Sessions) != 1 {
		t.Fatalf("ls %+v", ls)
	}
	row := ls.Sessions[0]
	if row.ID != s.ID || row.Status != "blocked" || row.Source != "hook" || row.Worktree != f.worktree ||
		row.PermissionMode != "manual" || row.Tab.Window != "window:w1" || row.Change != s.Change {
		t.Errorf("row %+v", row)
	}
	if !strings.Contains(human(t, ls), s.ID+" blocked (hook)") {
		t.Errorf("human %q", human(t, ls))
	}

	resp, err = callOp(t, f.ctx, "session.show", `{"id":"`+s.ID[:6]+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	show := resp.(SessionInfo)
	if show.Status.State != "blocked" || show.Status.Kind != "permission" || show.Source != "hook" {
		t.Errorf("status %+v from %s", show.Status, show.Source)
	}
	if show.Reports.Rex == nil || show.Reports.Rex.State != "working" || show.Reports.Exit != nil {
		t.Errorf("reports %+v", show.Reports)
	}
	h := human(t, show)
	for _, want := range []string{"- Status: blocked permission: Bash", "- rex: working", "- exit: none"} {
		if !strings.Contains(h, want) {
			t.Errorf("human lacks %q:\n%s", want, h)
		}
	}
	if len(f.rex.calls) != 1 {
		t.Errorf("rex calls %v", f.rex.calls)
	}

	_, err = callOp(t, f.ctx, "session.show", `{"id":"deadbeef"}`)
	wantExit(t, err, 4)
	_, err = callOp(t, f.ctx, "session.show", `{"id":""}`)
	wantExit(t, err, 2)
}

func TestSessionLsNeedsGit(t *testing.T) {
	ctx := initCtx(t)
	_, err := callOp(t, ctx, "session.ls", `{}`)
	wantExit(t, err, 1)
}

func TestNewSessionIDDeterministic(t *testing.T) {
	a, err := newSessionID(testCtx("").IDs)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := newSessionID(testCtx("").IDs)
	if a != b {
		t.Errorf("same seed, different IDs: %s %s", a, b)
	}
	if a[14] != '4' {
		t.Errorf("not a version 4 UUID: %s", a)
	}
}

func TestParseRexWindow(t *testing.T) {
	tab, err := parseRexWindow(session.Tab{Title: "t"}, []byte(rexWindowJSON))
	if err != nil || tab.Window != "window:w1" || tab.Block != "block:b1" || tab.Title != "t" {
		t.Errorf("tab %+v, %v", tab, err)
	}
	for _, out := range []string{"", "{}", `{"window_id":"w"}`, `{"window_id":"w","block_ids":[]}`, "not json"} {
		if _, err := parseRexWindow(session.Tab{}, []byte(out)); err == nil {
			t.Errorf("parseRexWindow(%q) = nil error", out)
		}
	}
}

// The agent's actor lives in the worktree only; its log is otherwise the
// launcher's.
func TestSessionStartWorktreeLogFolds(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	startSession(t, f.ctx, `{"permission_mode":"manual"}`)
	l := store.ReadLog(os.DirFS(filepath.Join(f.worktree, store.Dir)))
	if len(l.Problems) != 0 {
		t.Errorf("log problems %v", l.Problems)
	}
}

// fakeProcess is a started program that ends with code, and whose output,
// for a program that is not on the terminal, is out.
type fakeProcess struct {
	out    io.Reader
	code   int
	err    error
	killed bool
}

func (p *fakeProcess) Output() io.Reader  { return p.out }
func (p *fakeProcess) Wait() (int, error) { return p.code, p.err }
func (p *fakeProcess) Kill()              { p.killed = true }

// fakeStart records the programs a call starts. The agent (a Terminal
// program) ends with agent; rex events prints events.
type fakeStart struct {
	programs []Program
	agent    *fakeProcess
	events   string
	rex      *fakeProcess
	err      error
}

func (f *fakeStart) start(p Program) (Process, error) {
	f.programs = append(f.programs, p)
	if f.err != nil {
		return nil, f.err
	}
	if p.Terminal {
		return f.agent, nil
	}
	f.rex = &fakeProcess{out: strings.NewReader(f.events)}
	return f.rex, nil
}

// statusEvent is a `rex events --json` line reporting state for block.
func statusEvent(block, state, kind, msg string) string {
	return `{"type":"session_event","payload":{"session_id":"session:s1","event":{"type":"block_event","payload":{` +
		`"block_id":"` + block + `","creator_name":"com.superlogical.terminal","name":"program_status_changed",` +
		`"payload":{"id":"","record":{"state":"` + state + `","kind":"` + kind + `","msg":"` + msg + `",` +
		`"updated_at":"2026-10-24T12:00:00Z"},"reason":"report"}}}}}`
}

// runFixture is a launched session ready for session.run.
func runFixture(t *testing.T, startReq string) (*sessionFixture, SessionInfo, *fakeStart) {
	t.Helper()
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	s := startSession(t, f.ctx, startReq).Session
	fs := &fakeStart{agent: &fakeProcess{code: 3}}
	f.ctx.Start = fs.start
	return f, s, fs
}

func TestSessionRun(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"manual","run":"run-1"}`)
	f.ctx.RexSession, f.ctx.RexBlock = "session:s1", "block:b1"
	fs.events = statusEvent("block:other", "error", "", "not this tab") + "\n" +
		"not json\n" +
		statusEvent("block:b1", "working", "", "Running tests") + "\n"

	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID[:8]+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionRunResponse)
	if r.Code != 3 || r.ExitCode() != 3 || len(r.Warnings) != 0 {
		t.Errorf("code %d, exit %d, warnings %v", r.Code, r.ExitCode(), r.Warnings)
	}
	if len(fs.programs) != 2 {
		t.Fatalf("programs %+v", fs.programs)
	}
	ev, agent := fs.programs[0], fs.programs[1]
	if ev.Name != "rex" || strings.Join(ev.Args, " ") != "events -s session:s1 --json" || ev.Terminal {
		t.Errorf("rex program %+v", ev)
	}
	if !fs.rex.killed {
		t.Error("rex events was not stopped")
	}
	argv := strings.Join(append([]string{agent.Name}, agent.Args...), " ")
	for _, want := range []string{"--session-id " + s.ID, "--permission-mode manual", "/opt/bin/rhiza session hook"} {
		if !strings.Contains(argv, want) {
			t.Errorf("argv %q lacks %q", argv, want)
		}
	}
	if !agent.Terminal || agent.Dir != f.worktree {
		t.Errorf("agent terminal %v, dir %q", agent.Terminal, agent.Dir)
	}
	wantEnv := []string{"RHIZA_CHANGE=" + s.Change, "RHIZA_ACTOR=" + s.Actor, "RHIZA_FOR=" + s.For, "RHIZA_SESSION=" + s.ID, "RHIZA_RUN=run-1", "GOTMPDIR=" + run.TmpDir(run.Dir(f.common), "run-1")}
	if fi, err := os.Stat(run.TmpDir(run.Dir(f.common), "run-1")); err != nil || !fi.IsDir() {
		t.Errorf("the run's GOTMPDIR folder does not exist when the agent starts: %v", err)
	}
	if !slices.Equal(agent.Env, wantEnv) || !slices.Equal(agent.Unset, []string{"RHIZA_AUTONOMOUS", "RHIZA_MILESTONE", "RHIZA_LEAF", "RHIZA_ATTEMPT", "RHIZA_CHECK", "RHIZA_CYCLE", "RHIZA_ROLE", "RHIZA_COMMIT", "RHIZA_ROOT"}) {
		t.Errorf("env %v unset %v, want %v", agent.Env, agent.Unset, wantEnv)
	}

	rec, err := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if x := rec.Status.Exit; x == nil || x.State != session.StateExited || x.Code == nil || *x.Code != 3 {
		t.Errorf("exit report %+v", x)
	}
	if x := rec.Status.Rex; x == nil || x.State != session.StateWorking || x.Detail != "Running tests" {
		t.Errorf("rex report %+v", x)
	}
	if rec.Status.Hook != nil {
		t.Errorf("hook report %+v", rec.Status.Hook)
	}
	// The Rex report came first, then the exit: the exit is the status.
	if r.Session.Source != "exit" || r.Session.Status.State != session.StateExited {
		t.Errorf("resolved %s from %s", r.Session.Status.State, r.Session.Source)
	}

	// A hook report after the exit is ignored: it is not recorded and the
	// session stays exited.
	f.ctx.Stdin = strings.NewReader(`{"session_id":"` + s.ID + `","hook_event_name":"UserPromptSubmit"}`)
	hresp, herr := callOp(t, f.ctx, "session.hook", `{"event":"UserPromptSubmit"}`)
	if herr != nil {
		t.Fatal(herr)
	}
	if hr := hresp.(SessionHookResponse); hr.Recorded {
		t.Errorf("a hook after the exit was recorded: %+v", hr)
	}
	after, err := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status.Hook != nil {
		t.Errorf("hook report after the exit %+v", after.Status.Hook)
	}
	if got, src := session.Resolve(after); got.State != session.StateExited || src != session.SourceExit {
		t.Errorf("after a late hook: %s from %s", got.State, src)
	}
	if out := human(t, r); !strings.Contains(out, "exited with code 3") {
		t.Errorf("human %q", out)
	}
}

func TestSessionRunNoRunNoRex(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"plan"}`)
	fs.agent.code = 0
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.(SessionRunResponse).ExitCode() != 0 {
		t.Errorf("exit %d", resp.(SessionRunResponse).ExitCode())
	}
	if len(fs.programs) != 1 {
		t.Fatalf("outside Rex only the agent starts: %+v", fs.programs)
	}
	agent := fs.programs[0]
	if slices.ContainsFunc(agent.Env, func(kv string) bool { return strings.HasPrefix(kv, "RHIZA_RUN=") }) ||
		!slices.Equal(agent.Unset, []string{"RHIZA_RUN", "RHIZA_AUTONOMOUS", "RHIZA_MILESTONE", "RHIZA_LEAF", "RHIZA_ATTEMPT", "RHIZA_CHECK", "RHIZA_CYCLE", "RHIZA_ROLE", "RHIZA_COMMIT", "RHIZA_ROOT"}) {
		t.Errorf("env %v unset %v", agent.Env, agent.Unset)
	}
}

func TestSessionRunDryRun(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"manual"}`)
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionRunResponse)
	if !r.DryRun || len(r.Argv) == 0 || len(fs.programs) != 0 {
		t.Errorf("dry run %+v started %v", r, fs.programs)
	}
	if !strings.Contains(human(t, r), "nothing was started or written") {
		t.Errorf("human %q", human(t, r))
	}
}

func TestSessionRunRefusals(t *testing.T) {
	cases := []struct {
		name    string
		edit    func(*session.Record)
		exit    int
		message string
	}{
		{"NoMode", func(r *session.Record) { r.PermissionMode = "" }, 1, "no permission mode"},
		{"Bypass", func(r *session.Record) { r.PermissionMode = "bypassPermissions" }, 1, `trust = "dedicated"`},
		{"UnknownMode", func(r *session.Record) { r.PermissionMode = "yolo" }, 1, "unknown permission mode"},
		{"UnknownAgent", func(r *session.Record) { r.Agent = "nobody" }, 1, "does not know"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, s, fs := runFixture(t, `{"permission_mode":"manual"}`)
			st := session.NewStore(session.Dir(f.common))
			if _, err := st.Update(s.ID, c.edit); err != nil {
				t.Fatal(err)
			}
			_, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
			wantExit(t, err, c.exit)
			if !strings.Contains(AsError(err).Message, c.message) {
				t.Errorf("message %q lacks %q", AsError(err).Message, c.message)
			}
			if len(fs.programs) != 0 {
				t.Errorf("a refused run started %v", fs.programs)
			}
		})
	}

	t.Run("Unknown", func(t *testing.T) {
		f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
		_, err := callOp(t, f.ctx, "session.run", `{"id":"deadbeef"}`)
		wantExit(t, err, 4)
	})
	t.Run("NoStart", func(t *testing.T) {
		f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
		f.ctx.Start = nil
		_, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
		wantExit(t, err, 1)
	})
}

func TestSessionRunStartFails(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"manual"}`)
	fs.err = errors.New("executable file not found in $PATH")
	_, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	wantExit(t, err, 1)
	rec, rerr := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if x := rec.Status.Exit; x == nil || x.State != session.StateError || !strings.Contains(x.Detail, "not found") {
		t.Errorf("exit report %+v", x)
	}
}

// removedEvent is a `rex events --json` line removing block's status
// records for reason.
func removedEvent(block, reason string) string {
	return `{"type":"session_event","payload":{"session_id":"session:s1","event":{"type":"block_event","payload":{` +
		`"block_id":"` + block + `","creator_name":"com.superlogical.terminal","name":"program_status_removed",` +
		`"payload":{"ids":[""],"reason":"` + reason + `"}}}}}`
}

func TestRexStatus(t *testing.T) {
	cases := []struct {
		name, line        string
		ok, exited        bool
		state, kind, text string
	}{
		{"Working", statusEvent("block:b1", "working", "", "Building"), true, false, "working", "", "Building"},
		{"Blocked", statusEvent("block:b1", "blocked", "question", "Which one?"), true, false, "blocked", "question", "Which one?"},
		{"UnknownKind", statusEvent("block:b1", "blocked", "coffee", ""), true, false, "blocked", "", ""},
		{"UnknownState", statusEvent("block:b1", "napping", "", ""), true, false, "idle", "", ""},
		{"ChildExited", removedEvent("block:b1", "child_exited"), true, true, "exited", "", "rex: the terminal's program exited"},
		{"RemovedOtherReason", removedEvent("block:b1", "prompt"), false, false, "", "", ""},
		{"RemovedOtherBlock", removedEvent("block:b2", "child_exited"), false, false, "", "", ""},
		{"OtherBlock", statusEvent("block:b2", "working", "", ""), false, false, "", "", ""},
		{"OtherEvent", strings.Replace(statusEvent("block:b1", "working", "", ""), "program_status_changed", "cwd_changed", 1), false, false, "", "", ""},
		{"NotJSON", "nope", false, false, "", "", ""},
		{"Empty", "", false, false, "", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rep, exited, ok := rexStatus([]byte(c.line), "block:b1")
			if ok != c.ok || exited != c.exited || rep.State != c.state || rep.Kind != c.kind || rep.Detail != c.text {
				t.Errorf("rexStatus = %+v, %v, %v", rep, exited, ok)
			}
		})
	}
}

// TestRecordRexExitIsFinal feeds Rex reports in order: a working report,
// then Rex saying the terminal's program exited, then a late report. Rex
// alone ends the session, and nothing after the exit revives it.
func TestRecordRexExitIsFinal(t *testing.T) {
	var r session.Record
	feed := func(line string) {
		t.Helper()
		rep, exited, ok := rexStatus([]byte(line), "block:b1")
		if !ok {
			t.Fatalf("rexStatus ignored %s", line)
		}
		recordRex(&r, rep, exited)
	}
	feed(statusEvent("block:b1", "working", "", "Running tests"))
	if got, src := session.Resolve(r); got.State != session.StateWorking || src != session.SourceRex {
		t.Fatalf("live: %s from %s", got.State, src)
	}
	feed(removedEvent("block:b1", "child_exited"))
	if got, src := session.Resolve(r); got.State != session.StateExited || src != session.SourceExit || got.Code != nil {
		t.Fatalf("after exit: %+v from %s", got, src)
	}
	feed(statusEvent("block:b1", "blocked", "question", "late"))
	feed(removedEvent("block:b1", "child_exited"))
	if r.Status.Rex == nil || r.Status.Rex.State != session.StateWorking {
		t.Errorf("a report after the exit changed rex: %+v", r.Status.Rex)
	}
	if got, src := session.Resolve(r); got.State != session.StateExited || src != session.SourceExit {
		t.Errorf("a late report revived the session: %s from %s", got.State, src)
	}
}

// hook runs session.hook for event with payload on stdin and returns the
// response and what it wrote to stderr.
func hook(t *testing.T, ctx *Ctx, event, payload string) (SessionHookResponse, string) {
	t.Helper()
	var stderr strings.Builder
	ctx.Stdin, ctx.Stderr = strings.NewReader(payload), &stderr
	resp, err := callOp(t, ctx, "session.hook", `{"event":"`+event+`"}`)
	if err != nil {
		t.Fatalf("session.hook never fails, got %v", err)
	}
	r := resp.(SessionHookResponse)
	if out := human(t, r); out != "" {
		t.Errorf("a hook printed %q", out)
	}
	return r, stderr.String()
}

func TestSessionHook(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	st := session.NewStore(session.Dir(f.common))
	read := func() session.Record {
		t.Helper()
		rec, err := st.Read(s.ID)
		if err != nil {
			t.Fatal(err)
		}
		return rec
	}

	r, errOut := hook(t, f.ctx, "Notification",
		`{"session_id":"`+s.ID+`","notification_type":"permission_prompt","message":"May I?"}`)
	if !r.Recorded || r.Session != s.ID || r.State != "blocked" || r.Kind != "permission" || errOut != "" {
		t.Errorf("hook %+v, stderr %q", r, errOut)
	}
	if h := read().Status.Hook; h == nil || h.State != "blocked" || h.Kind != "permission" || h.Detail != "May I?" {
		t.Errorf("hook report %+v", h)
	}

	// An event that changes nothing leaves the report alone.
	r, errOut = hook(t, f.ctx, "Notification", `{"session_id":"`+s.ID+`","notification_type":"auth_success"}`)
	if r.Recorded || r.Session != s.ID || errOut != "" {
		t.Errorf("hook %+v, stderr %q", r, errOut)
	}
	if h := read().Status.Hook; h == nil || h.State != "blocked" {
		t.Errorf("hook report %+v", h)
	}

	// Only the hook source changes.
	r, _ = hook(t, f.ctx, "Stop", `{"session_id":"`+s.ID+`"}`)
	rec := read()
	if !r.Recorded || rec.Status.Hook.State != "idle" || rec.Status.Rex != nil || rec.Status.Exit != nil {
		t.Errorf("hook %+v, status %+v", r, rec.Status)
	}

	// A session Rhiza did not launch, or a prefix of one, is ignored
	// without a word and nothing is written.
	for _, sid := range []string{"00000000-0000-4000-8000-000000000000", s.ID[:8], ""} {
		r, errOut = hook(t, f.ctx, "Stop", `{"session_id":"`+sid+`"}`)
		if r.Recorded || r.Session != "" || errOut != "" {
			t.Errorf("session %q: hook %+v, stderr %q", sid, r, errOut)
		}
	}
	if !read().Status.Hook.At.Equal(rec.Status.Hook.At) {
		t.Error("an ignored hook wrote")
	}

	// Failures go to stderr only.
	for _, c := range []struct{ event, payload, want string }{
		{"Stop", "not json", "session hook Stop"},
		{"Bogus", `{"session_id":"` + s.ID + `"}`, "unknown hook event"},
		{"", `{}`, "no hook event"},
	} {
		r, errOut = hook(t, f.ctx, c.event, c.payload)
		if r.Recorded || !strings.Contains(errOut, c.want) {
			t.Errorf("%s %q: hook %+v, stderr %q", c.event, c.payload, r, errOut)
		}
	}

	// Outside a git repository the hook still succeeds.
	f.ctx.Git = nil
	if _, errOut = hook(t, f.ctx, "Stop", `{"session_id":"`+s.ID+`"}`); !strings.Contains(errOut, "git") {
		t.Errorf("stderr %q", errOut)
	}
}

func TestSessionDone(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	st := session.NewStore(session.Dir(f.common))

	// Neither an argument nor RHIZA_SESSION names a session: usage.
	f.ctx.Getenv = func(string) string { return "" }
	_, err := callOp(t, f.ctx, "session.done", `{}`)
	if e := AsError(err); e == nil || e.Code != CodeUsage || !strings.Contains(e.Message, "RHIZA_SESSION") || e.Next == "" {
		t.Fatalf("no session named: %v", err)
	}

	// A dry run writes nothing.
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "session.done", `{"id":"`+s.ID[:8]+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionDoneResponse); !r.DryRun || r.Session != s.ID {
		t.Errorf("dry run %+v", r)
	}
	if rec, _ := st.Read(s.ID); rec.Status.SaidDone() {
		t.Fatal("a dry run marked the session done")
	}
	f.ctx.DryRun = false

	// RHIZA_SESSION names it: the mark is recorded, and Resolve is as before.
	f.ctx.Getenv = func(k string) string {
		if k == "RHIZA_SESSION" {
			return s.ID
		}
		return ""
	}
	resp, err = callOp(t, f.ctx, "session.done", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionDoneResponse)
	if r.Session != s.ID || r.Already || r.At == "" {
		t.Errorf("done %+v", r)
	}
	rec, err := st.Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d := rec.Status.Done; d == nil || d.State != session.StateDone {
		t.Fatalf("done mark %+v", d)
	}
	if rec.Status.Hook != nil || rec.Status.Rex != nil || rec.Status.Exit != nil {
		t.Errorf("session done wrote another source: %+v", rec.Status)
	}
	if out := human(t, r); !strings.Contains(out, "Session "+s.ID+" is done.") {
		t.Errorf("human %q", out)
	}

	// Saying it again keeps the first mark.
	first := *rec.Status.Done
	resp, err = callOp(t, f.ctx, "session.done", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionDoneResponse); !r.Already || r.At != rfc3339(first.At) {
		t.Errorf("again %+v", r)
	}
	if again, _ := st.Read(s.ID); !again.Status.Done.At.Equal(first.At) {
		t.Errorf("the first mark changed: %+v", again.Status.Done)
	}

	// An unknown session is not found.
	_, err = callOp(t, f.ctx, "session.done", `{"id":"ffffffff-ffff"}`)
	if e := AsError(err); e == nil || e.Code != CodeNotFound {
		t.Errorf("unknown session: %v", err)
	}
}

func TestSessionDoneAfterExit(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	st := session.NewStore(session.Dir(f.common))
	code := 0
	if _, err := st.Update(s.ID, func(r *session.Record) {
		r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
	}); err != nil {
		t.Fatal(err)
	}
	_, err := callOp(t, f.ctx, "session.done", `{"id":"`+s.ID+`"}`)
	if e := AsError(err); e == nil || e.Code != CodeConflict || !strings.Contains(e.Message, "has exited") {
		t.Fatalf("done after exit: %v", err)
	}
	if rec, _ := st.Read(s.ID); rec.Status.SaidDone() {
		t.Error("a session that exited was marked done")
	}
}

// The answers are captured before the done mark is written, so a step that
// waits on the mark finds the Decisions; a capture that fails (here, no
// transcript) is a warning and the mark stands.
func TestSessionDoneCapturesBeforeMark(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	st := session.NewStore(session.Dir(f.common))
	f.ctx.Home = t.TempDir()
	f.ctx.Getenv = func(string) string { return "" }

	var markedAtCapture, called bool
	beforeDoneMark = func(rec session.Record) {
		called = true
		cur, err := st.Read(rec.ID)
		if err != nil {
			t.Error(err)
			return
		}
		markedAtCapture = cur.Status.SaidDone()
	}
	t.Cleanup(func() { beforeDoneMark = nil })

	resp, err := callOp(t, f.ctx, "session.done", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionDoneResponse)
	if !called || markedAtCapture {
		t.Errorf("hook called %v, done mark already set at capture %v", called, markedAtCapture)
	}
	if len(r.Warnings) == 0 || !strings.Contains(r.Warnings[0], "no transcript of session "+s.ID) {
		t.Errorf("warnings %v", r.Warnings)
	}
	if rec, _ := st.Read(s.ID); !rec.Status.SaidDone() {
		t.Error("the done mark was not kept")
	}
}

// Done lists the handles of the Decisions it proposed, in the response and
// in its human output, so the planner and the run see what was captured.
func TestSessionDoneReportsProposed(t *testing.T) {
	f, s, _ := interviewFixture(t)
	f.ctx.Getenv = func(string) string { return "" }
	writeTranscript(t, f.ctx.Home, s.ID, pickedAnswer, typedAnswer)
	resp, err := callOp(t, f.ctx, "session.done", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionDoneResponse)
	want := []string{"answer-" + s.ID[:8] + "-1", "answer-" + s.ID[:8] + "-2"}
	if !slices.Equal(r.Decisions, want) || len(r.Warnings) != 0 {
		t.Errorf("decisions %v warnings %v", r.Decisions, r.Warnings)
	}
	out := human(t, r)
	for _, h := range want {
		if !strings.Contains(out, "- Proposed decision "+h+"\n") {
			t.Errorf("human %q lacks %s", out, h)
		}
	}
}

func TestRecordAgentActorReusesSession(t *testing.T) {
	ww := testWorkspace(t, false)
	launcher, err := ww.ids.New()
	if err != nil {
		t.Fatal(err)
	}
	const sid = "0123abcd-0000-4000-8000-000000000000"
	cs, err := ww.ids.New()
	if err != nil {
		t.Fatal(err)
	}
	h1, cs1, err := recordAgentActor(ww, "claude-code", sid, cs, launcher, launcher)
	if err != nil {
		t.Fatal(err)
	}
	before := len(ww.changeSets())
	h2, cs2, err := recordAgentActor(ww, "claude-code", sid, cs, launcher, launcher)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 || cs1 != cs2 {
		t.Errorf("second call = %s %s, want %s %s", h2, cs2, h1, cs1)
	}
	if n := len(ww.changeSets()); n != before {
		t.Errorf("second call wrote change sets: %d -> %d", before, n)
	}
	h3, _, err := recordAgentActor(ww, "claude-code", "4567abcd-0000-4000-8000-000000000000", cs, launcher, launcher)
	if err != nil {
		t.Fatal(err)
	}
	if h3 == h1 || len(ww.changeSets()) != before+1 {
		t.Errorf("a different session got %s, change sets %d", h3, len(ww.changeSets()))
	}
}
