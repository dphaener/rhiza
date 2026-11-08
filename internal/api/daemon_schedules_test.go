package api

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

type schedHarness struct {
	t     *testing.T
	now   time.Time
	log   *strings.Builder
	progs []Program
	projs []*daemonProj
	rex   []*fakeExec
	last  string
}

// newSchedHarness makes one project per body of dispatch.toml, all on one
// fake clock, one fake Detach and one daemon log.
func newSchedHarness(t *testing.T, bodies ...string) *schedHarness {
	h := &schedHarness{t: t, now: time.Date(2026, 11, 8, 10, 30, 0, 0, time.UTC), log: &strings.Builder{}}
	for i, body := range bodies {
		f := newSessionFixture(t, "s")
		writeDispatch(t, f.ctx.Root, body)
		f.ctx.Now = func() time.Time { return h.now }
		f.ctx.Stderr = h.log
		f.ctx.Detach = func(p Program, _ string) error { h.progs = append(h.progs, p); return nil }
		f.ctx.Start = (&shStart{codes: map[string]int{"false": 1}}).start
		name := string(rune('a' + i))
		f.ctx.ProjectName = name
		f.ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
		common, err := gitCommonDir(f.ctx)
		if err != nil {
			t.Fatal(err)
		}
		h.projs = append(h.projs, &daemonProj{name: name, dir: f.ctx.Root, common: common, ctx: f.ctx,
			rs: run.NewStore(run.Dir(common)), pending: map[string]daemonPending{}, unjournaled: map[string]daemonPending{}})
		h.rex = append(h.rex, f.rex)
	}
	return h
}

// poll is one root poll over the projects with maxRuns machine-wide slots,
// as runDaemon does it. The disk guard is off: it reads the host's real
// free space, and a CI runner under the 10G default would hold every start.
func (h *schedHarness) poll(maxRuns int) {
	h.t.Helper()
	r := DaemonRequest{Autonomous: true, MinFree: "0"}
	turns := make([]daemonTurn, len(h.projs))
	var picks []daemon.ProjectPick
	running := 0
	for i, p := range h.projs {
		turns[i] = p.prepare(r, nil)
		if turns[i].err != nil {
			h.t.Fatal(turns[i].err)
		}
		picks = append(picks, turns[i].pick)
		running += turns[i].pick.Running
	}
	share := map[string]int{}
	picked := daemon.RoundRobin(picks, running, maxRuns, h.last)
	h.last = daemon.LastTurn(picked, h.last)
	for _, pk := range picked {
		share[pk.Project]++
	}
	resp := DaemonResponse{}
	for i, p := range h.projs {
		if err := p.poll(r, &resp, turns[i], maxRuns, share[p.name]); err != nil {
			h.t.Fatal(err)
		}
	}
}

func (h *schedHarness) started() []string {
	var out []string
	for _, p := range h.progs {
		out = append(out, strings.Join(p.Args, " "))
	}
	return out
}

const hourlyCmd = "[[schedules]]\nname = \"tick\"\nevery = \"1h\"\ncommand = \"true\"\n"

// A poll over two projects starts each due schedule's run; one slot full
// makes the other wait, and it fires once when a slot frees.
func TestDaemonPollFiresSchedules(t *testing.T) {
	h := newSchedHarness(t, hourlyCmd, strings.Replace(hourlyCmd, "tick", "tock", 1))
	h.poll(1) // first sight: nothing due
	if len(h.progs) != 0 {
		t.Fatalf("started before a slot: %v", h.started())
	}
	h.now = h.now.Add(31 * time.Minute) // 11:01, slot 11:00
	h.poll(1)
	// One machine-wide slot: only one project's schedule starts.
	if len(h.progs) != 1 || !strings.Contains(h.started()[0], "run start --foreground --schedule ") {
		t.Fatalf("started = %v, want one scheduled run", h.started())
	}
	first := h.started()[0]
	// Its run is journaled and running: the slot stays full, the other waits.
	h.journal(0, h.name(first), true)
	h.now = h.now.Add(time.Minute)
	h.poll(1)
	if len(h.progs) != 1 {
		t.Fatalf("started with no free slot: %v", h.started())
	}
	// It ends: the waiting schedule fires, once, however many polls.
	h.journal(0, h.name(first), false)
	h.now = h.now.Add(time.Minute)
	h.poll(1)
	h.poll(1)
	if len(h.progs) != 2 || h.started()[1] == first {
		t.Fatalf("started = %v, want the waiting schedule once", h.started())
	}
	for i := range h.rex {
		if c := h.rexCalls(i); len(c) != 0 {
			t.Errorf("rex was called: %v", c)
		}
	}
	if !strings.Contains(h.log.String(), "started a run of schedule") {
		t.Errorf("daemon log has no start line:\n%s", h.log)
	}
}

// name is the schedule a started program's args name.
func (h *schedHarness) name(args string) string {
	f := strings.Fields(args)
	for i, a := range f {
		if a == "--schedule" {
			return f[i+1]
		}
	}
	h.t.Fatalf("no --schedule in %q", args)
	return ""
}

// journal writes a scheduled run of schedule name into the project that owns
// it: running (a live process) or ended.
func (h *schedHarness) journal(_ int, name string, running bool) {
	h.t.Helper()
	for _, p := range h.projs {
		for _, p2 := range h.progs {
			if p2.Dir != p.ctx.Root || h.name(strings.Join(p2.Args, " ")) != name {
				continue
			}
			rid := "01M6WSQMY5Q33CWG319V9AR" + strings.ToUpper(name[:3])
			if _, err := p.rs.Read(rid); err != nil {
				at := h.now
				must := func(e run.Entry) {
					if err := p.rs.Append(rid, e); err != nil {
						h.t.Fatal(err)
					}
				}
				if err := p.rs.Create(rid, run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "schedule:" + name, Trigger: run.TriggerSchedule, TriggerName: name}}); err != nil {
					h.t.Fatal(err)
				}
				must(run.Entry{Type: run.TypeProcess, At: at, Process: &run.Process{Host: hostName(), Pid: selfPid()}})
			}
			if !running {
				if err := p.rs.Append(rid, run.Entry{Type: run.TypeFinished, At: h.now, End: run.Failed, Detail: `step "command" failed`}); err != nil {
					h.t.Fatal(err)
				}
			}
			return
		}
	}
	h.t.Fatalf("no started program for %s", name)
}

// A schedule whose workflow does not resolve starts nothing and the daemon
// log names the schedule; a standing problem is logged once.
func TestDaemonScheduleUnresolvedWorkflow(t *testing.T) {
	h := newSchedHarness(t, "[[schedules]]\nname = \"nightly\"\nevery = \"1h\"\nworkflow = \"gone\"\n")
	h.poll(2)
	h.now = h.now.Add(time.Hour)
	h.poll(2)
	h.poll(2)
	if len(h.progs) != 0 {
		t.Fatalf("started %v", h.started())
	}
	if n := strings.Count(h.log.String(), "schedule nightly not started"); n != 1 {
		t.Fatalf("log names the schedule %d times:\n%s", n, h.log)
	}
}

// A scheduled command run that starts and fails makes no Rex call; the
// daemon log carries its start and its failure.
func TestDaemonScheduledRunFailureLogged(t *testing.T) {
	h := newSchedHarness(t, "[[schedules]]\nname = \"bad\"\nevery = \"1h\"\ncommand = \"false\"\n")
	h.poll(2)
	h.now = h.now.Add(time.Hour)
	h.poll(2)
	if len(h.progs) != 1 {
		t.Fatalf("started %v", h.started())
	}
	p := h.projs[0]
	// The detached process, run here: rhiza run start --schedule bad.
	resp, err := callOp(t, p.ctx, "run.start", `{"schedule":"bad","permission_mode":"manual"}`)
	if err != nil {
		t.Fatal(err)
	}
	if st := resp.(RunResponse).Run.State; st != "failed" {
		t.Fatalf("run state = %q", st)
	}
	if c := h.rexCalls(0); len(c) != 0 {
		t.Errorf("a failing scheduled command called rex: %v", c)
	}
	// Its failure is the daemon log's only: no inbox item, no Rex notice.
	var notices bytes.Buffer
	p.ctx.Rex = &notices
	h.now = h.now.Add(time.Minute)
	h.poll(2)
	if strings.Contains(notices.String(), "run failed") || strings.Contains(notices.String(), "blocked") {
		t.Errorf("a failed scheduled run made a Rex notice: %q", notices.String())
	}
	inbox, err := callOp(t, p.ctx, "inbox", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if items := inbox.(InboxResponse).Items; len(items) != 0 {
		t.Errorf("a failed scheduled run is in the inbox: %+v", items)
	}
	out := h.log.String()
	if !strings.Contains(out, "started a run of schedule bad") || !strings.Contains(out, "of schedule bad failed") {
		t.Errorf("daemon log lacks the start or the failure:\n%s", out)
	}
	if c := h.rexCalls(0); len(c) != 0 {
		t.Errorf("rex was called: %v", c)
	}
}

// A scheduled run interrupted (its process gone, no finished entry) does not
// stop the schedule: the next slot after it fires.
func TestDaemonInterruptedScheduledRunNextSlotFires(t *testing.T) {
	h := newSchedHarness(t, hourlyCmd)
	h.poll(2)
	h.now = h.now.Add(31 * time.Minute) // 11:01, slot 11:00
	h.poll(2)
	if len(h.progs) != 1 {
		t.Fatalf("started %v", h.started())
	}
	p := h.projs[0]
	rid := "01M6WSQMY5Q33CWG319V9ARINT"
	at := h.now
	if err := p.rs.Create(rid, run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "schedule:tick", Trigger: run.TriggerSchedule, TriggerName: "tick"}}); err != nil {
		t.Fatal(err)
	}
	// A pid that is not running: the run was interrupted.
	if err := p.rs.Append(rid, run.Entry{Type: run.TypeProcess, At: at, Process: &run.Process{Host: hostName(), Pid: 1 << 30}}); err != nil {
		t.Fatal(err)
	}
	for _, step := range []time.Duration{time.Minute, 30 * time.Minute} {
		h.now = h.now.Add(step)
		h.poll(2)
	}
	if len(h.progs) != 1 {
		t.Fatalf("fired before the next slot: %v", h.started())
	}
	h.now = h.now.Add(30 * time.Minute) // 12:02, slot 12:00
	h.poll(2)
	if len(h.progs) != 2 {
		t.Fatalf("the slot after an interrupted run did not fire: %v\n%s", h.started(), h.log)
	}
}

func hostName() string { h, _ := os.Hostname(); return h }
func selfPid() int     { return os.Getpid() }

// rexCalls are the calls project i made to rex, the fake exec's other
// calls being git's.
func (h *schedHarness) rexCalls(i int) [][]string {
	var out [][]string
	for _, c := range h.rex[i].calls {
		if c[0] == "rex" {
			out = append(out, c)
		}
	}
	return out
}

// With slots for both, one poll over two projects starts each project's due
// schedule.
func TestDaemonPollStartsEachProjectsSchedule(t *testing.T) {
	h := newSchedHarness(t, hourlyCmd, strings.Replace(hourlyCmd, "tick", "tock", 1))
	h.poll(2)
	h.now = h.now.Add(31 * time.Minute)
	h.poll(2)
	if len(h.progs) != 2 || h.progs[0].Dir == h.progs[1].Dir {
		t.Fatalf("started = %v, want one run per project", h.started())
	}
}
