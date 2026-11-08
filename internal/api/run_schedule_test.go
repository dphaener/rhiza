package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

func writeDispatch(t *testing.T, root, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".rhiza", "dispatch.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A scheduled run of a command and of a repo workflow opens no change set and
// no worktree, runs in the main checkout, journals trigger schedule and the
// schedule's name, shows in run ls and run status with no change set, and
// the command's output reaches the run's log.
func TestScheduledRunCommandAndWorkflow(t *testing.T) {
	f := newSessionFixture(t, "sched")
	writeWorkflow(t, f.ctx.Root, "tidy", "[step.one]\nkind = \"command\"\nrun = \"tidy-up\"\n")
	writeDispatch(t, f.ctx.Root, `
[[schedules]]
name = "intake"
every = "15m"
command = "echo hello"

[[schedules]]
name = "nightly"
cron = "0 3 * * *"
workflow = "tidy"
`)
	var log strings.Builder
	f.ctx.Stderr = &log
	var programs []Program
	f.ctx.Start = func(p Program) (Process, error) {
		programs = append(programs, p)
		return &fakeProcess{out: strings.NewReader("hello from the command\n")}, nil
	}
	before, _ := os.ReadDir(filepath.Join(f.ctx.Root, ".rhiza", "log"))
	rs := run.NewStore(run.Dir(f.common))
	for _, c := range []struct{ name, wf string }{{"intake", "intake"}, {"nightly", "tidy"}} {
		resp, err := callOp(t, f.ctx, "run.start", `{"schedule":"`+c.name+`"}`)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		info := resp.(RunResponse).Run
		if info.State != "succeeded" || info.Change != "" || info.Worktree != "" || info.Workflow != c.wf ||
			info.Trigger != "schedule" || info.TriggerName != c.name || !info.Autonomous {
			t.Errorf("%s: run = %+v", c.name, info)
		}
		j, err := rs.Read(info.ID)
		if err != nil {
			t.Fatal(err)
		}
		s := j.Entries[0].Started
		if s.Trigger != run.TriggerSchedule || s.TriggerName != c.name || s.Change != "" || s.Worktree != "" || s.Branch != "" {
			t.Errorf("%s: started entry = %+v", c.name, s)
		}
	}
	if len(programs) != 2 {
		t.Fatalf("%d programs started, want 2", len(programs))
	}
	if programs[0].Dir != f.ctx.Root || programs[1].Dir != f.ctx.Root {
		t.Errorf("steps ran in %q and %q, want the main checkout %q", programs[0].Dir, programs[1].Dir, f.ctx.Root)
	}
	if !strings.HasSuffix(programs[0].Args[1], "echo hello") {
		t.Errorf("command step ran %q", programs[0].Args)
	}
	if !strings.Contains(log.String(), "hello from the command") {
		t.Errorf("run log %q lacks the command's output", log.String())
	}
	after, _ := os.ReadDir(filepath.Join(f.ctx.Root, ".rhiza", "log"))
	if len(after) != len(before) {
		t.Errorf("the model log grew from %d to %d entries: a change set was opened", len(before), len(after))
	}

	ls, err := callOp(t, f.ctx, "run.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	items := ls.(RunLsResponse).Runs
	if len(items) != 2 {
		t.Fatalf("run ls = %+v", items)
	}
	for _, it := range items {
		if it.Change != "" || it.Trigger != "schedule" || (it.TriggerName != "intake" && it.TriggerName != "nightly") {
			t.Errorf("run ls item = %+v", it)
		}
	}
	if h := human(t, ls); !strings.Contains(h, "schedule intake") {
		t.Errorf("run ls text %q does not name the schedule", h)
	}
	st, err := callOp(t, f.ctx, "run.status", `{"run":"`+items[0].ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if si := st.(RunInfo); si.Change != "" || si.Trigger != "schedule" || si.TriggerName != items[0].TriggerName {
		t.Errorf("run status = %+v", si)
	}
}

// A scheduled run's workflow with a leaves step fails the run naming it.
func TestScheduledRunLeavesStepFails(t *testing.T) {
	f := newSessionFixture(t, "sched")
	build, err := os.ReadFile(filepath.Join("..", "workflow", "builtin", "build.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if _, perr := workflow.Parse("build", build); perr != nil {
		t.Fatal(perr)
	}
	writeWorkflow(t, f.ctx.Root, "buildit", string(build))
	writeDispatch(t, f.ctx.Root, "[[schedules]]\nname = \"b\"\nevery = \"1h\"\nworkflow = \"buildit\"\n")
	f.ctx.Stderr = &strings.Builder{}
	f.ctx.Start = (&shStart{}).start
	resp, err := callOp(t, f.ctx, "run.start", `{"schedule":"b","permission_mode":"manual"}`)
	if err != nil {
		t.Fatal(err)
	}
	info := resp.(RunResponse).Run
	if info.State != "failed" || !strings.Contains(info.Detail, `step "build"`) || !strings.Contains(info.Detail, "needs a change set") {
		t.Errorf("run = state %q detail %q, want a failure naming step build", info.State, info.Detail)
	}
}

func TestScheduledRunRefusals(t *testing.T) {
	f := newSessionFixture(t, "sched")
	writeDispatch(t, f.ctx.Root, "[[schedules]]\nname = \"a\"\nevery = \"1h\"\ncommand = \"true\"\n")
	if _, err := callOp(t, f.ctx, "run.start", `{"schedule":"nope"}`); err == nil || !strings.Contains(err.Error(), `no schedule "nope"`) {
		t.Errorf("unknown schedule: %v", err)
	}
	if _, err := callOp(t, f.ctx, "run.start", `{"schedule":"a","workflow":"plan"}`); err == nil || !strings.Contains(err.Error(), "--workflow") {
		t.Errorf("schedule with --workflow: %v", err)
	}
}

// fakeCheckout answers the git reads of a main checkout's snapshot from
// fields a test's command changes.
type fakeCheckout struct {
	branch, head string
	files        map[string]string
}

func (c *fakeCheckout) exec(name string, args []string, dir string, env []string) ([]byte, error) {
	args = args[2:] // -C <dir>
	switch strings.Join(args, " ") {
	case "rev-parse HEAD":
		return []byte(c.head + "\n"), nil
	case "rev-parse --abbrev-ref HEAD":
		return []byte(c.branch + "\n"), nil
	case "diff HEAD --name-only":
		var out string
		for f := range c.files {
			out += f + "\n"
		}
		return []byte(out), nil
	}
	if len(args) == 3 && args[0] == "hash-object" {
		return []byte(c.files[args[2]] + "\n"), nil
	}
	return nil, nil
}

// A scheduled command that leaves the checkout alone records no warning; one
// that switches branch, commits or edits a tracked file ends with its own
// status and a warning in run status and the log naming what changed.
func TestScheduledRunWarnsOnMainCheckoutChange(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*fakeCheckout)
		want   string
	}{
		{"quiet", func(*fakeCheckout) {}, ""},
		{"branch", func(c *fakeCheckout) { c.branch = "feature" }, "branch from main to feature"},
		{"commit", func(c *fakeCheckout) { c.head = "bbbb2222" }, "HEAD from aaaa1111 to bbbb2222"},
		{"edit", func(c *fakeCheckout) { c.files = map[string]string{"main.go": "h2"} }, "main.go"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newSessionFixture(t, "sched")
			writeDispatch(t, f.ctx.Root, "[[schedules]]\nname = \"s\"\nevery = \"1h\"\ncommand = \"go\"\n")
			co := &fakeCheckout{branch: "main", head: "aaaa1111", files: map[string]string{}}
			f.ctx.Exec = co.exec
			var log strings.Builder
			f.ctx.Stderr = &log
			f.ctx.Start = func(p Program) (Process, error) {
				c.mutate(co)
				return &fakeProcess{out: strings.NewReader(""), code: 0}, nil
			}
			resp, err := callOp(t, f.ctx, "run.start", `{"schedule":"s"}`)
			if err != nil {
				t.Fatal(err)
			}
			info := resp.(RunResponse).Run
			if info.State != "succeeded" {
				t.Errorf("state = %q, want the command's own status", info.State)
			}
			st, err := callOp(t, f.ctx, "run.status", `{"run":"`+info.ID+`"}`)
			if err != nil {
				t.Fatal(err)
			}
			warns := st.(RunInfo).Warnings
			if c.want == "" {
				if len(warns) != 0 || strings.Contains(log.String(), "warning") {
					t.Errorf("warnings %v, log %q: want none", warns, log.String())
				}
				return
			}
			if len(warns) != 1 || !strings.Contains(warns[0], c.want) {
				t.Errorf("warnings = %v, want one naming %q", warns, c.want)
			}
			if !strings.Contains(log.String(), "warning") || !strings.Contains(log.String(), c.want) {
				t.Errorf("log %q lacks the warning naming %q", log.String(), c.want)
			}
		})
	}
}
