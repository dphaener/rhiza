package run

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

// fakeClock advances Now on Sleep and records each sleep.
type fakeClock struct {
	now    time.Time
	sleeps []time.Duration
}

func (c *fakeClock) Now() time.Time { return c.now }
func (c *fakeClock) Sleep(d time.Duration) {
	c.sleeps = append(c.sleeps, d)
	c.now = c.now.Add(d)
}

// fakeRunner scripts each step's results and counts calls.
type fakeRunner struct {
	codes    map[string][]int  // command exit codes per step, in call order
	holds    map[string][]bool // wait results per step; past the end, the last
	agentOut map[string][]bool // done per poll per session; past the end, true
	outcome  map[string]string // agent outcome per session
	failOn   map[string]error  // runner error per step
	commands map[string]int
	polls    map[string]int
	starts   map[string]int
	envs     map[string][]string
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		codes: map[string][]int{}, holds: map[string][]bool{}, agentOut: map[string][]bool{},
		outcome: map[string]string{}, failOn: map[string]error{},
		commands: map[string]int{}, polls: map[string]int{}, starts: map[string]int{},
		envs: map[string][]string{},
	}
}

func (f *fakeRunner) Command(step workflow.Step, env []string) (int, string, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return 0, "", err
	}
	n := f.commands[step.Name]
	f.commands[step.Name]++
	codes := f.codes[step.Name]
	if n < len(codes) {
		return codes[n], "", nil
	}
	return 0, "", nil
}

func (f *fakeRunner) Holds(step workflow.Step, env []string) (bool, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return false, err
	}
	n := f.polls[step.Name]
	f.polls[step.Name]++
	hs := f.holds[step.Name]
	switch {
	case len(hs) == 0:
		return false, nil
	case n < len(hs):
		return hs[n], nil
	}
	return hs[len(hs)-1], nil
}

func (f *fakeRunner) StartAgent(step workflow.Step, env []string) (string, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return "", err
	}
	f.starts[step.Name]++
	return "S-" + step.Name, nil
}

func (f *fakeRunner) AgentOutcome(session string) (bool, string, string, error) {
	n := f.polls[session]
	f.polls[session]++
	if seq := f.agentOut[session]; n < len(seq) && !seq[n] {
		return false, "", "", nil
	}
	return true, f.outcome[session], "exit 0", nil
}

const runID = "01RUNAAAAAAAAAAAAAAAAAAAAA"

// setup parses src, starts a run of it and returns an engine over fakes.
func setup(t *testing.T, src string) (*Engine, *fakeRunner, *fakeClock, *workflow.Definition) {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	if err := Start(store, runID, Started{Workflow: "wf", Definition: src, Change: "CS1"}, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	e := &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}
	return e, r, c, def
}

func outcomes(st State) []string {
	var out []string
	for _, s := range st.Steps {
		out = append(out, s.Step+":"+s.Outcome)
	}
	return out
}

func entryTypes(t *testing.T, e *Engine) []string {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, en := range j.Entries {
		out = append(out, en.Type+":"+en.Step)
	}
	return out
}

const waitWF = `
name = "wf"
[step.ci]
kind = "wait"
until = "gh pr checks"
interval = "1s"
timeout = "5s"
on.pass = "done"
on.timeout = "failed"
`

func TestWaitHeldOnThirdPoll(t *testing.T) {
	e, r, c, _ := setup(t, waitWF)
	r.holds["ci"] = []bool{false, false, true}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"ci:pass"}) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if r.polls["ci"] != 3 {
		t.Errorf("polls = %d, want 3", r.polls["ci"])
	}
	if want := []time.Duration{time.Second, time.Second}; !reflect.DeepEqual(c.sleeps, want) {
		t.Errorf("sleeps = %v, want %v", c.sleeps, want)
	}
}

func TestWaitTimesOut(t *testing.T) {
	e, r, c, _ := setup(t, waitWF)
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || !reflect.DeepEqual(outcomes(st), []string{"ci:timeout"}) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	// Checks at 0s..5s: the last one falls on the deadline.
	if r.polls["ci"] != 6 {
		t.Errorf("polls = %d, want 6", r.polls["ci"])
	}
	if len(c.sleeps) != 5 {
		t.Errorf("sleeps = %v, want five", c.sleeps)
	}
	if st.Detail != `step "ci" ended with outcome "timeout"` {
		t.Errorf("detail = %q", st.Detail)
	}
}

func TestWaitSleepsOnlyToDeadline(t *testing.T) {
	e, r, c, _ := setup(t, strings.Replace(waitWF, `interval = "1s"`, `interval = "3s"`, 1))
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || r.polls["ci"] != 3 {
		t.Fatalf("result %q polls %d", st.Result, r.polls["ci"])
	}
	if want := []time.Duration{3 * time.Second, 2 * time.Second}; !reflect.DeepEqual(c.sleeps, want) {
		t.Errorf("sleeps = %v, want %v", c.sleeps, want)
	}
}

const loopWF = `
name = "wf"
[step.build]
kind = "command"
run = "make"
max_visits = 2
on.pass = "test"

[step.test]
kind = "command"
run = "make test"
max_visits = 3
on.fail = "build"
`

func TestMaxVisitsFailsRun(t *testing.T) {
	e, r, _, _ := setup(t, loopWF)
	r.codes["test"] = []int{1, 1, 1}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed {
		t.Fatalf("result = %q, want failed", st.Result)
	}
	if want := `step "build" went over max_visits 2`; st.Detail != want {
		t.Fatalf("detail = %q, want %q", st.Detail, want)
	}
	t.Logf("run %s: %s", st.Result, st.Detail)
	if r.commands["build"] != 2 || r.commands["test"] != 2 {
		t.Errorf("commands = %v", r.commands)
	}
	if st.Visits["build"] != 2 || st.Visits["test"] != 2 {
		t.Errorf("visits = %v", st.Visits)
	}
}

func TestDefaultMaxVisitsIsOne(t *testing.T) {
	e, r, _, _ := setup(t, strings.Replace(loopWF, "max_visits = 2\n", "", 1))
	r.codes["test"] = []int{1}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if want := `step "build" went over max_visits 1`; st.Result != Failed || st.Detail != want {
		t.Fatalf("result %q detail %q", st.Result, st.Detail)
	}
}

const cmdWF = `
name = "wf"
[step.test]
kind = "command"
run = "make test"
on.pass = "ship"
on.fail = "fix"

[step.ship]
kind = "command"
run = "true"

[step.fix]
kind = "command"
run = "true"
on.pass = "failed"
`

func TestCommandPassFollowsOnPass(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"test:pass", "ship:pass"}) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if st.Steps[0].Next != "ship" || st.Steps[0].Detail != "exit 0" {
		t.Errorf("first step = %+v", st.Steps[0])
	}
	want := []string{"RHIZA_CHANGE=CS1", "RHIZA_RUN=" + runID, "RHIZA_STEP=test"}
	if !reflect.DeepEqual(r.envs["test"], want) {
		t.Errorf("env = %v, want %v", r.envs["test"], want)
	}
	if r.commands["fix"] != 0 {
		t.Errorf("fix ran")
	}
}

func TestCommandFailFollowsOnFail(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	r.codes["test"] = []int{2}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(outcomes(st), []string{"test:fail", "fix:pass"}) {
		t.Fatalf("steps %v", outcomes(st))
	}
	if st.Steps[0].Detail != "exit 2" {
		t.Errorf("detail = %q", st.Steps[0].Detail)
	}
	// fix's on.pass = "failed" ends the run failed.
	if st.Result != Failed || st.Detail != `step "fix" ended with outcome "pass"` {
		t.Errorf("result %q detail %q", st.Result, st.Detail)
	}
}

func TestCommandRunnerErrorFailsRun(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	r.failOn["test"] = errors.New("sh: not found")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(outcomes(st), []string{"test:error"}) || st.Result != Failed || st.Detail != "sh: not found" {
		t.Fatalf("steps %v result %q detail %q", outcomes(st), st.Result, st.Detail)
	}
}

func TestCommandJournalUsesNowNotClock(t *testing.T) {
	e, _, _, _ := setup(t, cmdWF)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	j, _ := e.Store.Read(runID)
	for _, en := range j.Entries[1:] {
		if !en.At.Equal(at(1)) {
			t.Errorf("%s entry at %v, want %v", en.Type, en.At, at(1))
		}
	}
}

const resumeWF = `
name = "wf"
[step.prep]
kind = "command"
run = "echo one >> count"
on.pass = "ci"

[step.ci]
kind = "wait"
until = "test -f go"
interval = "1s"
timeout = "5s"
`

func TestResumeAfterNoWaitDoesNotRedo(t *testing.T) {
	e, r, _, _ := setup(t, resumeWF)
	e.NoWait = true
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Ended || st.Current != "ci" || st.Stopped != "waiting: test -f go" {
		t.Fatalf("ended %v current %q stopped %q", st.Ended, st.Current, st.Stopped)
	}
	if r.polls["ci"] != 1 {
		t.Errorf("polls under NoWait = %d, want 1", r.polls["ci"])
	}

	// Stopped again: still one visit of ci, prep still run once.
	if st, err = e.Run(runID); err != nil || st.Visits["ci"] != 1 || st.Ended {
		t.Fatalf("second stop: err %v visits %v ended %v", err, st.Visits, st.Ended)
	}

	r.holds["ci"] = []bool{true}
	e.NoWait = false
	st, err = e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"prep:pass", "ci:pass"}) {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if r.commands["prep"] != 1 {
		t.Errorf("prep ran %d times, want 1", r.commands["prep"])
	}
	if st.Visits["ci"] != 1 || st.Visits["prep"] != 1 {
		t.Errorf("visits = %v", st.Visits)
	}
	want := []string{
		"started:prep", "entered:prep", "ended:prep", "entered:ci", "stopped:ci",
		"stopped:ci", "ended:ci", "finished:ci",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
}

func TestResumeFinishedRunIsErrFinished(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	st, err := e.Run(runID[:6])
	if !errors.Is(err, ErrFinished) {
		t.Fatalf("err = %v, want ErrFinished", err)
	}
	if !strings.Contains(err.Error(), runID) || st.Result != Succeeded {
		t.Errorf("err %q result %q", err, st.Result)
	}
	if r.commands["test"] != 1 {
		t.Errorf("test ran %d times", r.commands["test"])
	}
}

func TestResumeLiveRunIsErrLive(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	unlock, ok, err := e.Store.TryLock(runID)
	if err != nil || !ok {
		t.Fatalf("lock: ok %v err %v", ok, err)
	}
	defer unlock()
	_, err = e.Run(runID)
	if !errors.Is(err, ErrLive) {
		t.Fatalf("err = %v, want ErrLive", err)
	}
	if want := "run " + runID + " is running in another process"; err.Error() != want {
		t.Errorf("err = %q, want %q", err, want)
	}
	if len(r.commands) != 0 {
		t.Errorf("ran %v while live", r.commands)
	}
}

func TestResumeKilledMidStepContinuesVisit(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	// A process entered test and died before ending it.
	if err := e.Store.Append(runID, Entry{Type: TypeEntered, At: at(1), Step: "test", Visit: 1}); err != nil {
		t.Fatal(err)
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || st.Visits["test"] != 1 || r.commands["test"] != 1 {
		t.Fatalf("result %q visits %v commands %v", st.Result, st.Visits, r.commands)
	}
}

func TestResumeEndedWithoutFinishedWritesFinished(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "test", Visit: 1},
		{Type: TypeEnded, At: at(1), Step: "test", Visit: 1, Outcome: "fail", Next: "fix"},
		{Type: TypeEntered, At: at(1), Step: "fix", Visit: 1},
		{Type: TypeEnded, At: at(1), Step: "fix", Visit: 1, Outcome: "pass", End: Failed},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || st.Detail != `step "fix" ended with outcome "pass"` || len(r.commands) != 0 {
		t.Fatalf("result %q detail %q commands %v", st.Result, st.Detail, r.commands)
	}
}

const agentWF = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build {change}"
interval = "2s"
on.pass = "done"
`

func TestAgentPollsUntilDone(t *testing.T) {
	e, r, c, _ := setup(t, agentWF)
	r.agentOut["S-build"] = []bool{false, false}
	r.outcome["S-build"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.starts["build"] != 1 || r.polls["S-build"] != 3 {
		t.Fatalf("result %q starts %v polls %v", st.Result, r.starts, r.polls)
	}
	if st.Steps[0].Session != "S-build" {
		t.Errorf("session = %q", st.Steps[0].Session)
	}
	if want := []time.Duration{2 * time.Second, 2 * time.Second}; !reflect.DeepEqual(c.sleeps, want) {
		t.Errorf("sleeps = %v", c.sleeps)
	}
}

func TestAgentFailEndsRunFailed(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	r.outcome["S-build"] = "fail"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || st.Detail != `step "build" ended with outcome "fail"` {
		t.Fatalf("result %q detail %q", st.Result, st.Detail)
	}
}

func TestAgentResumedDoesNotStartAgain(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "build", Visit: 1},
		{Type: TypeSession, At: at(1), Step: "build", Session: "S-earlier"},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	r.outcome["S-earlier"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if r.starts["build"] != 0 {
		t.Fatalf("StartAgent called %d times on resume", r.starts["build"])
	}
	if st.Result != Succeeded || st.Steps[0].Session != "S-earlier" || st.Visits["build"] != 1 {
		t.Fatalf("result %q steps %+v visits %v", st.Result, st.Steps, st.Visits)
	}
}

func TestAgentStartErrorFailsRun(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	r.failOn["build"] = errors.New("adapter refused")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || st.Detail != "adapter refused" {
		t.Fatalf("result %q detail %q", st.Result, st.Detail)
	}
}

// TestDecide covers the pure transition apart from any execution.
func TestDecide(t *testing.T) {
	def, err := workflow.Parse("wf", []byte(cmdWF+`
[step.loop]
kind = "command"
run = "x"
max_visits = 2
`))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		st   State
		want Action
		err  string
	}{
		{"finished", State{Ended: true, Result: Succeeded}, Action{Kind: ActFinished}, ""},
		{"fresh enters start", State{Current: "test", Visits: map[string]int{}},
			Action{Kind: ActEnter, Step: "test", Visit: 1}, ""},
		{"entered continues", State{Current: "test", CurrentVisit: 1, Visits: map[string]int{"test": 1}},
			Action{Kind: ActContinue, Step: "test", Visit: 1}, ""},
		{"second visit under cap", State{Current: "loop", Visits: map[string]int{"loop": 1}},
			Action{Kind: ActEnter, Step: "loop", Visit: 2}, ""},
		{"over cap", State{Current: "loop", Visits: map[string]int{"loop": 2}},
			Action{Kind: ActOverCap, Step: "loop", Visit: 3, End: workflow.EndFailed,
				Detail: `step "loop" went over max_visits 2`}, ""},
		{"default cap is one", State{Current: "test", Visits: map[string]int{"test": 1}},
			Action{Kind: ActOverCap, Step: "test", Visit: 2, End: workflow.EndFailed,
				Detail: `step "test" went over max_visits 1`}, ""},
		{"ended pass with no entry finishes succeeded",
			State{Steps: []StepDone{{Step: "ship", Outcome: "pass"}}},
			Action{Kind: ActFinish, Step: "ship", End: workflow.EndSucceeded}, ""},
		{"ended to failed target",
			State{Steps: []StepDone{{Step: "fix", Outcome: "pass"}}},
			Action{Kind: ActFinish, Step: "fix", End: workflow.EndFailed,
				Detail: `step "fix" ended with outcome "pass"`}, ""},
		{"ended with error",
			State{Steps: []StepDone{{Step: "test", Outcome: OutcomeError, Detail: "boom"}}},
			Action{Kind: ActFinish, Step: "test", End: workflow.EndFailed, Detail: "boom"}, ""},
		{"unknown step", State{ID: "R", Current: "gone"}, Action{}, `step "gone" is not in workflow`},
		{"no start", State{ID: "R"}, Action{}, "no step to start at"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Decide(def, tc.st)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("err = %v, want %q", err, tc.err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v\nwant %+v", got, tc.want)
			}
		})
	}
}
