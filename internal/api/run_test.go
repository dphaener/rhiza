package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

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
	for _, ref := range []string{"", "  ", ".", ".."} {
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
