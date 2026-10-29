package api

import (
	"bytes"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// shStart answers `sh -c <cmd>` with the exit code codes[cmd] (0 if
// absent) and records each program.
type shStart struct {
	codes    map[string]int
	programs []Program
}

func (s *shStart) start(p Program) (Process, error) {
	s.programs = append(s.programs, p)
	cmd := ""
	if len(p.Args) == 2 {
		cmd = p.Args[1]
	}
	return &fakeProcess{out: strings.NewReader(""), code: s.codes[cmd]}, nil
}

// writeWorkflow writes a named workflow under root's .rhiza/workflows.
func writeWorkflow(t *testing.T, root, name, body string) string {
	t.Helper()
	dir := filepath.Join(root, ".rhiza", "workflows")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name+".toml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

const twoCommands = `
[step.one]
kind = "command"
run = "echo one > step1.txt"
on.pass = "two"

[step.two]
kind = "command"
run = "test -f step1.txt"
`

const agentWorkflow = `
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "done"
`

func TestWorkflowFile(t *testing.T) {
	root := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ ref, want string }{
		{"hello", filepath.Join(root, ".rhiza", "workflows", "hello.toml")},
		{" hello ", filepath.Join(root, ".rhiza", "workflows", "hello.toml")},
		{"flows/hello.toml", filepath.Join(wd, "flows", "hello.toml")},
		{"hello.toml", filepath.Join(wd, "hello.toml")},
		{"/abs/x", "/abs/x"},
	} {
		got, err := workflowFile(root, c.ref)
		if err != nil || got != c.want {
			t.Errorf("workflowFile(%q) = %q, %v; want %q", c.ref, got, err, c.want)
		}
	}
	if got, err := workflowFile(root, ""); err != nil || got != filepath.Join(root, ".rhiza", "workflows", "ship-change.toml") {
		t.Errorf("workflowFile(\"\") = %q, %v; want the ship-change name", got, err)
	}
	for _, ref := range []string{".", ".."} {
		_, err := workflowFile(root, ref)
		wantExit(t, err, 2)
	}
}

func TestRunStartWorkflowRefusals(t *testing.T) {
	cases := []struct {
		name, workflow, body, message string
	}{
		{"MissingName", "nope", "", `no workflow "nope"`},
		{"MissingPath", "elsewhere/nope.toml", "", "does not exist"},
		{"Invalid", "bad", "[step.a]\nkind = \"teleport\"\n", "teleport"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newSessionFixture(t, "x")
			newChange(t, f.ctx, `{"intent":"X"}`)
			if c.body != "" {
				writeWorkflow(t, f.ctx.Root, c.workflow, c.body)
			}
			_, err := callOp(t, f.ctx, "run.start", `{"workflow":"`+c.workflow+`"}`)
			wantExit(t, err, 2)
			if e := AsError(err); !strings.Contains(e.Message, c.message) || e.Next != runStartNext {
				t.Errorf("error %q next %q; want %q and %q", e.Message, e.Next, c.message, runStartNext)
			}
			if len(f.git.calls) != 0 {
				t.Errorf("a refused run ran git %v", f.git.calls)
			}
		})
	}
}

func TestRunStartAgentRefusalsBeforeWrite(t *testing.T) {
	cases := []struct {
		name, body, req string
		exit            int
		message         string
	}{
		{"NoMode", agentWorkflow, `{"workflow":"w"}`, 2, "no permission mode"},
		{"Bypass", agentWorkflow, `{"workflow":"w","permission_mode":"bypassPermissions"}`, 1, `trust = "dedicated"`},
		{"StepBypass", agentWorkflow + "permission_mode = \"bypassPermissions\"\n", `{"workflow":"w","permission_mode":"manual"}`, 1, `trust = "dedicated"`},
		{"UnknownAgent", agentWorkflow, `{"workflow":"w","permission_mode":"manual","agent":"nobody"}`, 2, "unknown agent"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newSessionFixture(t, "x")
			newChange(t, f.ctx, `{"intent":"X"}`)
			writeWorkflow(t, f.ctx.Root, "w", c.body)
			sh := &shStart{}
			f.ctx.Start = sh.start
			_, err := callOp(t, f.ctx, "run.start", c.req)
			wantExit(t, err, c.exit)
			if e := AsError(err); !strings.Contains(e.Message, `step "build"`) || !strings.Contains(e.Message, c.message) {
				t.Errorf("message %q lacks the step or %q", e.Message, c.message)
			}
			if len(f.git.calls) != 0 || len(f.rex.calls) != 0 || len(sh.programs) != 0 {
				t.Errorf("a refused run ran git %v, rex %v or %d programs", f.git.calls, f.rex.calls, len(sh.programs))
			}
			if _, err := os.Stat(run.Dir(f.common)); !os.IsNotExist(err) {
				t.Errorf("a refused run wrote %s: %v", run.Dir(f.common), err)
			}
		})
	}
}

func TestRunStartCommandSteps(t *testing.T) {
	f := newSessionFixture(t, "try-runs")
	cs := newChange(t, f.ctx, `{"intent":"Try runs","title":"Try runs"}`).ChangeSet
	writeWorkflow(t, f.ctx.Root, "hello", twoCommands)
	sh := &shStart{}
	f.ctx.Start = sh.start

	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(RunResponse)
	info := r.Run
	if info.State != run.Succeeded || r.ExitCode() != 0 || info.Workflow != "hello" ||
		info.Change != cs.String() || info.Worktree != f.worktree {
		t.Errorf("run = %+v, exit %d", info, r.ExitCode())
	}
	var got []string
	for _, s := range info.Steps {
		got = append(got, s.Step+":"+s.Outcome+":"+s.Next)
	}
	if want := []string{"one:pass:two", "two:pass:"}; !slices.Equal(got, want) {
		t.Errorf("steps %v, want %v", got, want)
	}
	if len(sh.programs) != 2 {
		t.Fatalf("ran %d programs, want 2", len(sh.programs))
	}
	for _, p := range sh.programs {
		if p.Name != "sh" || p.Args[0] != "-c" || p.Dir != f.worktree || p.Terminal ||
			!slices.Contains(p.Env, "RHIZA_RUN="+info.ID) || !slices.Contains(p.Env, "RHIZA_CHANGE="+cs.String()) {
			t.Errorf("program %+v", p)
		}
	}
	if _, err := os.Stat(filepath.Join(run.Dir(f.common), info.ID+".jsonl")); err != nil {
		t.Errorf("no journal: %v", err)
	}
	if !strings.HasPrefix(human(t, r), "Run "+info.ID+" succeeded\n") {
		t.Errorf("markdown:\n%s", human(t, r))
	}

	st, err := callOp(t, f.ctx, "run.status", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if s := st.(RunInfo); s.ID != info.ID || s.State != run.Succeeded || len(s.Steps) != 2 {
		t.Errorf("status = %+v", s)
	}
	_, err = callOp(t, f.ctx, "run.resume", `{"run":"`+info.ID+`"}`)
	wantExit(t, err, 3)
	if e := AsError(err); e.Next != "rhiza run status "+info.ID {
		t.Errorf("next = %q", e.Next)
	}
}

func TestRunStartFailedExitsNonZero(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "hello", twoCommands)
	f.ctx.Start = (&shStart{codes: map[string]int{"test -f step1.txt": 1}}).start

	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(RunResponse)
	if r.Run.State != run.Failed || r.ExitCode() != 1 || !strings.Contains(r.Run.Detail, `step "two"`) {
		t.Errorf("run = %+v, exit %d", r.Run, r.ExitCode())
	}
}

func TestRunStopsAtWaitAndResumes(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "w", `
[step.prep]
kind = "command"
run = "prep"
on.pass = "gate"

[step.gate]
kind = "wait"
until = "ready"
timeout = "1m"
`)
	sh := &shStart{codes: map[string]int{"ready": 1}}
	f.ctx.Start = sh.start

	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"w","no_wait":true}`)
	if err != nil {
		t.Fatal(err)
	}
	info := resp.(RunResponse).Run
	if info.State != runStopped || info.Current != "gate" || info.Detail != "waiting: ready" {
		t.Fatalf("run = %+v", info)
	}
	if !strings.Contains(human(t, resp), "next: rhiza run resume "+info.ID) {
		t.Errorf("markdown:\n%s", human(t, resp))
	}
	st, err := callOp(t, f.ctx, "run.status", `{"run":"`+info.ID[:6]+`"}`)
	if err != nil || st.(RunInfo).State != runStopped {
		t.Fatalf("status = %+v, %v", st, err)
	}

	sh.codes["ready"] = 0
	resp, err = callOp(t, f.ctx, "run.resume", `{"run":"`+info.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	info = resp.(RunResponse).Run
	if info.State != run.Succeeded || len(info.Steps) != 2 {
		t.Errorf("resumed run = %+v", info)
	}
	prep := 0
	for _, p := range sh.programs {
		if p.Args[1] == "prep" {
			prep++
		}
	}
	if prep != 1 {
		t.Errorf("prep ran %d times, want 1", prep)
	}
}

func TestRunStatusAndResumeRefusals(t *testing.T) {
	f := newSessionFixture(t, "x")
	_, err := callOp(t, f.ctx, "run.status", `{}`)
	wantExit(t, err, 4)
	_, err = callOp(t, f.ctx, "run.status", `{"run":"ZZZZZZ"}`)
	wantExit(t, err, 4)
	_, err = callOp(t, f.ctx, "run.resume", `{"run":""}`)
	wantExit(t, err, 2)
}

func TestRunStartDryRunWritesNothing(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "hello", twoCommands)
	sh := &shStart{}
	f.ctx.Start = sh.start
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(RunResponse)
	if !r.DryRun || r.Run.ID == "" || r.Run.Current != "one" {
		t.Errorf("dry run = %+v", r)
	}
	if len(sh.programs) != 0 {
		t.Errorf("a dry run ran %d programs", len(sh.programs))
	}
	if _, err := os.Stat(run.Dir(f.common)); !os.IsNotExist(err) {
		t.Errorf("a dry run wrote %s: %v", run.Dir(f.common), err)
	}
}

func TestRunAgentStepEndsWithTheSessionExit(t *testing.T) {
	for _, c := range []struct {
		code    int
		state   string
		outcome string
	}{{0, run.Succeeded, "pass"}, {2, run.Failed, "fail"}} {
		f := newSessionFixture(t, "x")
		newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
		writeWorkflow(t, f.ctx.Root, "w", agentWorkflow+"on.fail = \"failed\"\n")
		f.ctx.Start = (&shStart{}).start
		sessions := session.NewStore(session.Dir(f.common))
		// The fake rex runs the tab's command at once: the agent exits
		// with c.code before the tab is reported open.
		f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
			f.rex.run(name, args, dir, env)
			sid := args[len(args)-1]
			code := c.code
			if _, err := sessions.Update(sid, func(r *session.Record) {
				r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
			}); err != nil {
				t.Fatal(err)
			}
			return []byte(rexWindowJSON), nil
		}
		resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"w","permission_mode":"manual"}`)
		if err != nil {
			t.Fatal(err)
		}
		info := resp.(RunResponse).Run
		if info.State != c.state || len(info.Steps) != 1 || info.Steps[0].Outcome != c.outcome || info.Steps[0].Session == "" {
			t.Fatalf("exit %d: run = %+v", c.code, info)
		}
		rec, err := sessions.Read(info.Steps[0].Session)
		if err != nil || rec.Run != info.ID || rec.Step != "build" || rec.Worktree != f.worktree ||
			rec.Prompt != "build "+info.Change {
			t.Errorf("session = %+v, %v", rec, err)
		}
		if len(f.rex.calls) != 1 {
			t.Errorf("rex ran %d times, want 1", len(f.rex.calls))
		}
	}
}

// scriptClock is a run.Clock whose every Sleep (one per poll of an agent
// step) runs the next step of a script, so a test drives a session's
// record between polls.
type scriptClock struct {
	t      *testing.T
	now    time.Time
	script []func()
	sleeps int
}

func (c *scriptClock) Now() time.Time { return c.now }

func (c *scriptClock) Sleep(d time.Duration) {
	c.now = c.now.Add(d)
	if c.sleeps >= len(c.script) {
		c.t.Fatalf("poll %d: the script has ended but the run still polls", c.sleeps+1)
	}
	c.sleeps++
	c.script[c.sleeps-1]()
}

// useClock makes run.start and run.resume poll by c for the test.
func useClock(t *testing.T, c run.Clock) {
	t.Helper()
	old := runClock
	runClock = c
	t.Cleanup(func() { runClock = old })
}

// rexReports splits what was written to Rex into OSC 7501 reports, each
// a map of its fields with title and msg base64-decoded.
func rexReports(t *testing.T, out string) []map[string]string {
	t.Helper()
	var reports []map[string]string
	for _, part := range strings.Split(out, "\x1b]7501;")[1:] {
		body, ok := strings.CutSuffix(part, "\x1b\\")
		if !ok {
			t.Fatalf("report %q is not ended by ST", part)
		}
		fields := map[string]string{}
		for _, kv := range strings.Split(body, ":") {
			k, v, _ := strings.Cut(kv, "=")
			if k == "title" || k == "msg" {
				d, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					t.Fatalf("%s %q: %v", k, v, err)
				}
				v = string(d)
			}
			fields[k] = v
		}
		reports = append(reports, fields)
	}
	return reports
}

// A run's agent step is driven blocked, still blocked, answered, blocked
// again and exited: run.status reads blocked with the tab and question
// only while it is, and Rex gets one blocked report per block and one
// working report per end.
func TestRunBlockedStatusAndRexReports(t *testing.T) {
	f := newSessionFixture(t, "blocked-runs")
	newChange(t, f.ctx, `{"intent":"X","title":"Blocked runs"}`)
	writeWorkflow(t, f.ctx.Root, "w", agentWorkflow+"interval = \"1s\"\n")
	f.ctx.Start = (&shStart{}).start
	rex := &bytes.Buffer{}
	f.ctx.Rex = rex
	sessions := session.NewStore(session.Dir(f.common))
	var sid string
	hook := func(state, kind, detail string) {
		t.Helper()
		if _, err := sessions.Update(sid, func(r *session.Record) {
			r.Status.Hook = &session.Report{State: state, Kind: kind, Detail: detail}
		}); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		f.rex.run(name, args, dir, env)
		if sid != "" {
			// The advisor the block launched: it asks nothing.
			return []byte(rexWindowJSON), nil
		}
		sid = args[len(args)-1]
		hook(session.StateBlocked, session.KindQuestion, "Ship it?")
		return []byte(rexWindowJSON), nil
	}
	status := func() RunInfo {
		t.Helper()
		resp, err := callOp(t, f.ctx, "run.status", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(RunInfo)
	}
	blockedReports := func() int {
		n := 0
		for _, r := range rexReports(t, rex.String()) {
			if r["state"] == "blocked" {
				n++
			}
		}
		return n
	}
	clock := &scriptClock{t: t, now: time.Unix(0, 0)}
	clock.script = []func(){
		func() {
			info := status()
			b := info.Blocked
			if info.State != runBlocked || b == nil || b.Step != "build" || b.Session != sid ||
				b.Kind != session.KindQuestion || b.Question != "Ship it?" ||
				b.Tab.Window != "window:w1" || b.Tab.Block != "block:b1" || b.Tab.Title == "" {
				t.Fatalf("blocked status = %+v, blocked %+v", info, b)
			}
			out := human(t, info)
			for _, want := range []string{
				"Run " + info.ID + " blocked\n",
				"- Blocked: build waits on the owner in tab " + b.Tab.Title + " (window window:w1)\n",
				"- Question: Ship it?\n",
				"next: rex focus -b block:b1\n",
			} {
				if !strings.Contains(out, want) {
					t.Errorf("markdown lacks %q:\n%s", want, out)
				}
			}
			if n := blockedReports(); n != 1 {
				t.Errorf("%d blocked reports after the first poll, want 1", n)
			}
		},
		func() {}, // still blocked
		func() { hook(session.StateWorking, "", "") },
		func() {
			if info := status(); info.State != runRunning || info.Blocked != nil {
				t.Errorf("answered status = %+v", info)
			}
			hook(session.StateBlocked, session.KindPermission, "Allow rm?")
		},
		func() {
			code := 0
			if _, err := sessions.Update(sid, func(r *session.Record) {
				r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
			}); err != nil {
				t.Fatal(err)
			}
		},
	}
	useClock(t, clock)

	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"w","permission_mode":"manual"}`)
	if err != nil {
		t.Fatal(err)
	}
	info := resp.(RunResponse).Run
	if info.State != run.Succeeded || info.Blocked != nil || clock.sleeps != len(clock.script) {
		t.Fatalf("run = %+v after %d polls", info, clock.sleeps)
	}
	reports := rexReports(t, rex.String())
	var states []string
	for _, r := range reports {
		states = append(states, r["state"])
		if r["app"] != "rhiza" {
			t.Errorf("report %v: app is not rhiza", r)
		}
	}
	// The run's status (working at its step, then done) frames the
	// blocked notices; each block's end restores that status in full, as
	// a report replaces its record.
	if want := []string{"working", "blocked", "working", "blocked", "working", "done"}; !slices.Equal(states, want) {
		t.Fatalf("reports %v, want %v", states, want)
	}
	first, title := reports[1], reports[1]["title"]
	wantMsg := "Blocked runs: build waits on you in tab " + title + ": Ship it?"
	if title == "" || first["kind"] != "question" || first["msg"] != wantMsg {
		t.Errorf("first block report %v; want msg %q", first, wantMsg)
	}
	if reports[3]["kind"] != "permission" || !strings.HasSuffix(reports[3]["msg"], ": Allow rm?") {
		t.Errorf("second block report %v", reports[3])
	}
	for _, i := range []int{0, 2, 4} {
		if r := reports[i]; r["msg"] != "step build" || r["title"] != "Blocked runs" || r["id"] != "" {
			t.Errorf("working report %v, want the run's root status at step build", r)
		}
	}
	if r := reports[5]; r["msg"] != "run succeeded" || r["title"] != "Blocked runs" {
		t.Errorf("final report %v", r)
	}
	if s := status(); s.State != run.Succeeded || s.Blocked != nil {
		t.Errorf("final status = %+v", s)
	}
}

// Without Rex nothing is written: the run still reports blocked.
func TestRunBlockedWithoutRexWritesNothing(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "w", agentWorkflow)
	f.ctx.Start = (&shStart{}).start
	stderr := &bytes.Buffer{}
	f.ctx.Stderr = stderr
	sessions := session.NewStore(session.Dir(f.common))
	var sid string
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		f.rex.run(name, args, dir, env)
		if sid != "" {
			// The advisor the block launched: it asks nothing.
			return []byte(rexWindowJSON), nil
		}
		sid = args[len(args)-1]
		if _, err := sessions.Update(sid, func(r *session.Record) {
			r.Status.Hook = &session.Report{State: session.StateBlocked, Kind: session.KindQuestion}
		}); err != nil {
			t.Fatal(err)
		}
		return []byte(rexWindowJSON), nil
	}
	useClock(t, &scriptClock{t: t, script: []func(){func() {
		code := 0
		if _, err := sessions.Update(sid, func(r *session.Record) {
			r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
		}); err != nil {
			t.Fatal(err)
		}
	}}})
	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"w","permission_mode":"manual"}`)
	if err != nil || resp.(RunResponse).Run.State != run.Succeeded {
		t.Fatalf("run = %+v, %v", resp, err)
	}
	if strings.Contains(stderr.String(), "\x1b") {
		t.Errorf("stderr holds an escape sequence:\n%q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "step build: blocked on question") {
		t.Errorf("stderr lacks the blocked progress line:\n%s", stderr.String())
	}
}

// A blocked run with no question, in a tab not yet opened, reads blocked
// without a question line and points at the session.
func TestRunInfoHumanBlockedWithoutQuestion(t *testing.T) {
	info := RunInfo{
		ID: "R1", State: runBlocked, Current: "build", Steps: []RunStep{},
		Blocked: &RunBlocked{Step: "build", Session: "S1", Kind: "auth", Tab: SessionTab{Title: "x-build"}},
	}
	out := human(t, info)
	if !strings.Contains(out, "- Blocked: build waits on the owner in tab x-build\n") ||
		strings.Contains(out, "Question") || !strings.Contains(out, "next: rhiza session show S1\n") {
		t.Errorf("markdown:\n%s", out)
	}
}

func TestRunInfoAcceptsAndQuestions(t *testing.T) {
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	st := run.State{
		ID: "R1", Autonomous: true, Ended: true, Result: run.Succeeded,
		Accepts: []run.Accept{{Change: "C1", LogEntry: "E1", Verb: run.TypeAccepted, At: at}},
		Asked: []run.Asked{
			{Step: "build", Question: "Q1", Handle: "open-q", Text: "Which?", At: at},
			{Step: "build", Question: "Q2", Handle: "failed-q", Text: "How?", At: at},
			{Step: "build", Question: "Q3", Handle: "answered-q", Text: "Why?", At: at, Answered: "use-x"},
			{Step: "build", Question: "Q4", Handle: "passed-q", Text: "When?", At: at},
		},
		Arbitrations: []run.Arbitration{
			{Step: "build", Session: "S1", Question: "Q1", At: at},
			{Step: "build", Session: "S2", Question: "Q2", Outcome: "fail", At: at},
			{Step: "build", Session: "S3", Question: "Q3", Outcome: "pass", Decision: "use-x", At: at},
			{Step: "build", Session: "S4", Question: "Q4", Outcome: "pass", Decision: "use-y", At: at},
		},
	}
	info := runInfo(st, false, nil)
	want := []RunAccept{{Change: "C1", Entry: "E1", Verb: "accepted", At: "2026-01-02T03:04:05Z"}}
	if !slices.Equal(info.Accepts, want) {
		t.Errorf("accepts = %+v", info.Accepts)
	}
	got := []RunQuestion{}
	for _, q := range info.Questions {
		got = append(got, RunQuestion{Question: q.Question, Outcome: q.Outcome, Decision: q.Decision, Answered: q.Answered})
	}
	wantQ := []RunQuestion{
		{Question: "Q1"},
		{Question: "Q2", Outcome: "fail"},
		{Question: "Q3", Outcome: "pass", Decision: "use-x", Answered: true},
		{Question: "Q4", Outcome: "pass", Decision: "use-y"},
	}
	if !slices.Equal(got, wantQ) {
		t.Errorf("questions = %+v", got)
	}
	out := human(t, info)
	for _, s := range []string{
		"## Accepted in the owner's place\n\n- 2026-01-02T03:04:05Z accepted change set C1 (log entry E1)\n",
		"## Questions put to the arbiter\n\n",
		"- open-q (Q1): Which? — open\n",
		"- failed-q (Q2): How? — not decided: the arbiter's verdict was fail\n",
		"- answered-q (Q3): Why? — decision use-x\n",
		"- passed-q (Q4): When? — decision use-y, not yet delivered\n",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("markdown lacks %q:\n%s", s, out)
		}
	}
	plain := human(t, runInfo(run.State{ID: "R2"}, false, nil))
	if strings.Contains(plain, "owner's place") || strings.Contains(plain, "put to the arbiter") {
		t.Errorf("a normal run's markdown:\n%s", plain)
	}
}

func TestRexReportEncoder(t *testing.T) {
	long := strings.Repeat("é", 200) // 400 bytes
	msg := "line one\nline two\x1b[31m\tend" + strings.Repeat("x", 3000)
	out := string(rexReport("blocked", "question", long, msg))
	if !strings.HasPrefix(out, "\x1b]7501;state=blocked:kind=question:app=rhiza:title=") || !strings.HasSuffix(out, "\x1b\\") {
		t.Fatalf("report %q", out)
	}
	body := strings.TrimSuffix(strings.TrimPrefix(out, "\x1b]7501;"), "\x1b\\")
	for _, r := range body {
		if unicode.IsControl(r) {
			t.Fatalf("report body holds control character %U", r)
		}
	}
	r := rexReports(t, out)[0]
	if title := r["title"]; len(title) > rexTitleMax || !utf8.ValidString(title) || title != strings.Repeat("é", 96) {
		t.Errorf("title is %d bytes, valid %v", len(title), utf8.ValidString(title))
	}
	m := r["msg"]
	if len(m) != rexMsgMax || !strings.HasPrefix(m, "line one line two [31m end") {
		t.Errorf("msg is %d bytes: %.40q", len(m), m)
	}
	for _, c := range m {
		if unicode.IsControl(c) {
			t.Fatalf("msg holds control character %U", c)
		}
	}
	for _, c := range []struct{ state, kind, want string }{
		{"blocked", "nonsense", "\x1b]7501;state=blocked:app=rhiza\x1b\\"},
		{"working", "question", "\x1b]7501;state=working:app=rhiza\x1b\\"},
	} {
		if got := string(rexReport(c.state, c.kind, "", "")); got != c.want {
			t.Errorf("rexReport(%s, %s) = %q, want %q", c.state, c.kind, got, c.want)
		}
	}
	end := string(rexNotice(run.Notice{Step: "build", Block: &run.Block{Kind: "question"}}))
	if end != "\x1b]7501;state=working:app=rhiza\x1b\\" {
		t.Errorf("end notice %q", end)
	}
	n := rexNotice(run.Notice{Change: "CS1", Step: "build", Session: "S1", Blocked: true, Block: &run.Block{}})
	if got := rexReports(t, string(n))[0]["msg"]; got != "CS1: build waits on you in tab S1" {
		t.Errorf("msg with no title, tab or question = %q", got)
	}
}

func TestResolveWorkflow(t *testing.T) {
	root := t.TempDir()
	for _, ref := range []string{"", "ship-change"} {
		path, def, data, err := resolveWorkflow(root, ref)
		if err != nil || path != "builtin:ship-change" || def.Name != "ship-change" || len(data) == 0 {
			t.Errorf("resolveWorkflow(%q) = %q, %v, %v; want the built-in", ref, path, def, err)
		}
	}
	writeWorkflow(t, root, "ship-change", twoCommands)
	path, def, _, err := resolveWorkflow(root, "")
	want := filepath.Join(root, ".rhiza", "workflows", "ship-change.toml")
	if err != nil || path != want || def.Start != "one" {
		t.Errorf("a repo file did not replace the built-in: %q, %v, %v", path, def, err)
	}
	_, _, _, err = resolveWorkflow(root, "nope")
	wantExit(t, err, 2)
	if e := AsError(err); !strings.Contains(e.Message, "ship-change") {
		t.Errorf("error %q does not name the built-ins", e.Message)
	}
}

func TestStepRunnerDirFallsBackToMainCheckout(t *testing.T) {
	main, gone := t.TempDir(), filepath.Join(t.TempDir(), "gone")
	// Git cannot run from a removed folder; the main checkout comes from
	// the common dir, with the project's folder within it.
	ctx := &Ctx{Root: filepath.Join(gone, "app"), Git: func(args ...string) (string, error) {
		return "", errors.New("chdir: no such file or directory")
	}}
	if got := (&stepRunner{ctx: ctx, worktree: main, main: main}).dir(); got != main {
		t.Errorf("dir with a worktree = %q; want %q", got, main)
	}
	if got, want := (&stepRunner{ctx: ctx, worktree: gone, main: main}).dir(), filepath.Join(main, "app"); got != want {
		t.Errorf("dir without a worktree = %q; want the main checkout's project folder %q", got, want)
	}
}

// A failed run is refused by a plain resume, with --retry as the next
// command; a dry run of --retry writes nothing and names the step; resume
// --retry reruns the failed step and the run succeeds, the retry shown in
// its steps and retries. A retry of the succeeded run is a conflict.
func TestRunResumeRetry(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "hello", twoCommands)
	sh := &shStart{codes: map[string]int{"test -f step1.txt": 1}}
	f.ctx.Start = sh.start

	resp, err := callOp(t, f.ctx, "run.start", `{"workflow":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	id := resp.(RunResponse).Run.ID
	_, err = callOp(t, f.ctx, "run.resume", `{"run":"`+id+`"}`)
	wantExit(t, err, 3)
	if e := AsError(err); e.Next != "rhiza run resume --retry "+id || !strings.Contains(e.Message, "is finished: it failed") {
		t.Errorf("refusal = %q, next %q", e.Message, e.Next)
	}

	f.ctx.DryRun = true
	resp, err = callOp(t, f.ctx, "run.resume", `{"run":"`+id+`","retry":true}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RunResponse); !r.DryRun || r.Retry != "two" || r.Run.Current != "two" || r.Run.State != runStopped {
		t.Errorf("dry run = %+v", r)
	}
	if !strings.Contains(human(t, resp), "would be retried at step two") {
		t.Errorf("markdown:\n%s", human(t, resp))
	}
	st, err := callOp(t, f.ctx, "run.status", `{"run":"`+id+`"}`)
	if err != nil || st.(RunInfo).State != run.Failed || len(st.(RunInfo).Retries) != 0 {
		t.Fatalf("a dry run changed the run: %+v, %v", st, err)
	}
	f.ctx.DryRun = false

	sh.codes["test -f step1.txt"] = 0
	ran := len(sh.programs)
	resp, err = callOp(t, f.ctx, "run.resume", `{"run":"`+id+`","retry":true}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(RunResponse)
	if r.Run.State != run.Succeeded || r.Retry != "two" || r.ExitCode() != 0 {
		t.Fatalf("retried run = %+v", r)
	}
	if got := len(sh.programs) - ran; got != 1 {
		t.Errorf("retry ran %d programs, want 1 (step two only)", got)
	}
	steps := r.Run.Steps
	if len(steps) != 3 || steps[1].Outcome != "fail" || steps[2].Step != "two" || !steps[2].Retry || steps[2].Visit != 2 {
		t.Errorf("steps = %+v", steps)
	}
	if len(r.Run.Retries) != 1 || r.Run.Retries[0].Step != "two" || !strings.Contains(r.Run.Retries[0].Detail, `step "two"`) {
		t.Errorf("retries = %+v", r.Run.Retries)
	}
	out := human(t, resp)
	for _, want := range []string{"- two (visit 1): fail, then end", "- two (visit 2, retried): pass, then end", "## Retries", "retried at two"} {
		if !strings.Contains(out, want) {
			t.Errorf("markdown lacks %q:\n%s", want, out)
		}
	}

	_, err = callOp(t, f.ctx, "run.resume", `{"run":"`+id+`","retry":true}`)
	wantExit(t, err, 3)
	if e := AsError(err); !strings.Contains(e.Message, "only a failed run can be retried") {
		t.Errorf("refusal = %q", e.Message)
	}
	_, err = callOp(t, f.ctx, "run.resume", `{"run":"`+id+`"}`)
	wantExit(t, err, 3)
	if e := AsError(err); e.Next != "rhiza run status "+id {
		t.Errorf("next = %q", e.Next)
	}
}
