package run

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
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
	codes      map[string][]int         // command exit codes per step, in call order
	holds      map[string][]bool        // wait results per step; past the end, the last
	agentOut   map[string][]bool        // done per poll per session; past the end, true
	agentSeq   map[string][]AgentStatus // status per poll per session; past the end, agentOut
	agentErr   map[string]error         // error per session once its scripted polls run out
	onPoll     func(session string, n int)
	outcome    map[string]string // agent outcome per session
	failOn     map[string]error  // runner error per step
	commands   map[string]int
	polls      map[string]int
	starts     map[string]int
	models     []string // each StartAgent call's model, in order
	envs       map[string][]string
	answers    map[string]Answer // Answer per question; missing is not OK
	answerN    map[string]int    // Answer calls per question
	arbiters   []string          // "<step> <question>" per StartArbiter call
	asked      []string          // "<session> <text>" per AddQuestion call
	askErr     error             // AddQuestion's error, when set
	delivered  []string          // "<session> <question> <decision> <statement>" per Deliver call
	deliverErr error             // Deliver's error, when set
	ended      []string          // "step:session" per AgentEnded call, in order
	endedErr   error             // what AgentEnded returns
	onCommand  func(step string) // called as each command step runs
	leaves     []Leaf            // what Leaves returns
	leavesErr  error             // Leaves' error, when set
	leafSpecs  []LeafSpec        // each StartLeaf call's spec, in order
	leafErr    map[string]error  // StartLeaf's error per "<handle> <role>"
	// leafCommit is LeafCommit's commit and detail per leaf handle; a
	// missing handle has commit "C-<handle>".
	leafCommit map[string][2]string
	leafSince  []string // "<handle> <since>" per LeafCommit call
	// evidence is Evidence's ok and detail per leaf handle; a missing
	// handle has the verifier's own Evidence (ok).
	evidence      map[string]fakeEvidence
	evidenceCalls []string // "<handle> <session>" per Evidence call
	// evidenceCommitErr is CommitEvidence's error per leaf handle; a
	// missing handle commits "E-<handle>".
	evidenceCommitErr map[string]error
	commitCalls       []string // "<handle> <session>" per CommitEvidence call
	checkouts         []string // "open <dir> <commit>" and "remove <dir>" per call, in order
	checkoutErr       error    // OpenCheckout's error, when set
	removeErr         error    // RemoveCheckout's error, when set
	heads             []string // Head's result per call, in order; past the end, fakeHead
	// leafQuestions are "<handle> <text>" per AddLeafQuestion call; the
	// nth question is QL<n>, handle leaf-failed-<n>.
	leafQuestions   []string
	leafQuestionErr error // AddLeafQuestion's error, when set
	// leafAnswers are LeafAnswer's decision and statement per question;
	// a missing question has no answer.
	leafAnswers map[string][2]string
	leafAnswerN map[string]int // LeafAnswer calls per question
	// store is the engine's journal store, where a check session's fake
	// agent records its findings round as StartCheck launches it.
	store *Store
	// noRound names the check sessions whose agent records no findings
	// round; every other check session records an empty one.
	noRound    map[string]bool
	closed     []string    // the session of each ClosePane call, in order
	closeErr   error       // ClosePane's error, when set
	checkSpecs []CheckSpec // each StartCheck call's spec, in order
	checkErr   error       // StartCheck's error, when set
}

// fakeEvidence is one leaf's scripted Evidence result.
type fakeEvidence struct {
	ok     bool
	detail string
	err    error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		codes: map[string][]int{}, holds: map[string][]bool{}, agentOut: map[string][]bool{},
		agentSeq: map[string][]AgentStatus{}, agentErr: map[string]error{},
		outcome: map[string]string{}, failOn: map[string]error{},
		commands: map[string]int{}, polls: map[string]int{}, starts: map[string]int{},
		envs: map[string][]string{}, answers: map[string]Answer{}, answerN: map[string]int{},
		leafAnswers: map[string][2]string{}, leafAnswerN: map[string]int{},
		noRound: map[string]bool{},
	}
}

func (f *fakeRunner) Command(step workflow.Step, env []string) (int, string, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return 0, "", err
	}
	n := f.commands[step.Name]
	f.commands[step.Name]++
	if f.onCommand != nil {
		f.onCommand(step.Name)
	}
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

func (f *fakeRunner) StartAgent(step workflow.Step, model string, env []string) (string, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return "", err
	}
	f.starts[step.Name]++
	f.models = append(f.models, model)
	return "S-" + step.Name, nil
}

func (f *fakeRunner) AgentEnded(step workflow.Step, session string) error {
	f.ended = append(f.ended, step.Name+":"+session)
	return f.endedErr
}

func (f *fakeRunner) AgentStatus(session string) (AgentStatus, error) {
	n := f.polls[session]
	f.polls[session]++
	if f.onPoll != nil {
		f.onPoll(session, n)
	}
	if seq := f.agentSeq[session]; n < len(seq) {
		return seq[n], nil
	}
	if seq := f.agentOut[session]; n < len(seq) && !seq[n] {
		return AgentStatus{}, nil
	}
	if err := f.agentErr[session]; err != nil {
		return AgentStatus{}, err
	}
	return AgentStatus{Done: true, Outcome: f.outcome[session], Detail: "exit 0"}, nil
}

func (f *fakeRunner) StartArbiter(step workflow.Step, question string, env []string) (string, error) {
	f.envs[step.Name] = env
	if err := f.failOn[step.Name]; err != nil {
		return "", err
	}
	f.arbiters = append(f.arbiters, step.Name+" "+question)
	return "S-" + step.Name + "-" + question, nil
}

func (f *fakeRunner) AddQuestion(session, text string) (Question, error) {
	if f.askErr != nil {
		return Question{}, f.askErr
	}
	f.asked = append(f.asked, session+" "+text)
	n := strconv.Itoa(len(f.asked))
	return Question{ID: "Q" + n, Handle: "question-s-" + n}, nil
}

func (f *fakeRunner) Deliver(session, question, decision, statement string) error {
	f.delivered = append(f.delivered, session+" "+question+" "+decision+" "+statement)
	return f.deliverErr
}

func (f *fakeRunner) Answer(question, _ string) (Answer, error) {
	f.answerN[question]++
	if a, ok := f.answers[question]; ok {
		return a, nil
	}
	return Answer{Why: "no live Decision answers question " + question}, nil
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
	if err := Start(store, runID, Started{Workflow: "wf", Definition: src, Change: "CS1", Title: "Blocked runs"}, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.store = store
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

const ladderWF = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build {change}"
models = ["sonnet", "opus"]
max_visits = 3
on.pass = "done"
on.fail = "build"
`

// An agent step revisited climbs its model ladder by visit, the last
// model kept past the end of the list.
func TestAgentModelLadder(t *testing.T) {
	e, r, _, _ := setup(t, ladderWF)
	r.outcome["S-build"] = workflow.OutcomeFail
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || st.Visits["build"] != 3 {
		t.Fatalf("result %q visits %v", st.Result, st.Visits)
	}
	if want := []string{"sonnet", "opus", "opus"}; !reflect.DeepEqual(r.models, want) {
		t.Errorf("models = %q, want %q", r.models, want)
	}
	t.Logf("visits 1, 2, 3 started on %q", r.models)
}

// An agent step with no models starts on the agent's default.
func TestAgentModelDefault(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	r.outcome["S-build"] = workflow.OutcomePass
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if want := []string{""}; !reflect.DeepEqual(r.models, want) {
		t.Errorf("models = %q, want %q", r.models, want)
	}
}

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

// Fake session statuses for the blocked-step tests.
var (
	idle     = AgentStatus{}
	question = &Block{
		Kind: "question", Question: "Ship it?",
		TabSession: "rex-1", TabWindow: "w-2", TabBlock: "b-3", TabTitle: "build",
	}
	blocked = AgentStatus{Blocked: question}
)

// noticeLog collects the engine's notices.
func noticeLog(e *Engine) *[]Notice {
	var ns []Notice
	e.Notify = func(n Notice) { ns = append(ns, n) }
	return &ns
}

func countType(t *testing.T, e *Engine, typ string) int {
	t.Helper()
	n := 0
	for _, en := range entryTypes(t, e) {
		if strings.HasPrefix(en, typ+":") {
			n++
		}
	}
	return n
}

// TestAgentBlockedAndBack drives a fake session idle, blocked, blocked,
// idle and exited: one blocked and one unblocked entry, one notice for the
// block and one for its end, and the folded state reads blocked between.
func TestAgentBlockedAndBack(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	ns := noticeLog(e)
	var lines []string
	e.Log = func(l string) { lines = append(lines, l) }
	r.agentSeq["S-build"] = []AgentStatus{idle, blocked, blocked, idle}
	r.outcome["S-build"] = "pass"
	seen := map[int]*Block{}
	r.onPoll = func(session string, n int) {
		j, err := e.Store.Read(runID)
		if err != nil {
			t.Fatal(err)
		}
		st := j.State()
		seen[n] = st.Blocked
		if st.Blocked != nil && (st.BlockedStep != "build" || st.BlockedSession != session) {
			t.Errorf("poll %d: blocked step %q session %q", n, st.BlockedStep, st.BlockedSession)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || st.Blocked != nil {
		t.Fatalf("result %q blocked %+v", st.Result, st.Blocked)
	}
	// State at each poll, before the poll's own status is acted on.
	for n, want := range []bool{false, false, true, true, false} {
		if got := seen[n] != nil; got != want {
			t.Errorf("poll %d: state blocked = %v, want %v", n, got, want)
		}
	}
	if !reflect.DeepEqual(seen[2], question) {
		t.Errorf("folded block = %+v, want %+v", seen[2], question)
	}
	want := []string{
		"started:build", "entered:build", "session:build", "blocked:build",
		"unblocked:build", "ended:build", "finished:build",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	wantNotices := []Notice{
		{Run: runID, Change: "CS1", Title: "Blocked runs", Step: "build", Session: "S-build", Blocked: true, Block: question},
		{Run: runID, Change: "CS1", Title: "Blocked runs", Step: "build", Session: "S-build", Blocked: false, Block: question},
	}
	if !reflect.DeepEqual(*ns, wantNotices) {
		t.Errorf("notices = %+v\nwant %+v", *ns, wantNotices)
	}
	var progress []string
	for _, l := range lines {
		if strings.Contains(l, "blocked") {
			progress = append(progress, l)
		}
	}
	if want := []string{"step build: blocked on question", "step build: unblocked"}; !reflect.DeepEqual(progress, want) {
		t.Errorf("log = %v, want %v", progress, want)
	}
}

// TestAgentResumeWhileBlocked resumes a step whose journal already shows
// its session blocked: no second blocked entry or notice, one for the end.
func TestAgentResumeWhileBlocked(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "build", Visit: 1},
		{Type: TypeSession, At: at(1), Step: "build", Session: "S-earlier"},
		{Type: TypeBlocked, At: at(1), Step: "build", Session: "S-earlier", Block: question},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	ns := noticeLog(e)
	r.agentSeq["S-earlier"] = []AgentStatus{blocked, blocked, blocked}
	r.outcome["S-earlier"] = "pass"
	blockedNotices := func() int {
		n := 0
		for _, x := range *ns {
			if x.Blocked {
				n++
			}
		}
		return n
	}
	r.onPoll = func(_ string, n int) {
		if n == 3 { // still blocked through three polls
			if b := blockedNotices(); b != 0 {
				t.Errorf("blocked notices while still blocked = %d, want 0", b)
			}
			if c := countType(t, e, TypeBlocked); c != 1 {
				t.Errorf("blocked entries while still blocked = %d, want 1", c)
			}
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.starts["build"] != 0 || r.polls["S-earlier"] != 4 {
		t.Fatalf("result %q starts %v polls %v", st.Result, r.starts, r.polls)
	}
	t.Logf("blocked notices on resume: %d; blocked entries: %d", blockedNotices(), countType(t, e, TypeBlocked))
	if blockedNotices() != 0 || countType(t, e, TypeBlocked) != 1 {
		t.Fatalf("blocked notices %d, blocked entries %d; want 0 and 1", blockedNotices(), countType(t, e, TypeBlocked))
	}
	if len(*ns) != 1 || (*ns)[0].Blocked || (*ns)[0].Session != "S-earlier" || !reflect.DeepEqual((*ns)[0].Block, question) {
		t.Fatalf("notices = %+v, want one end of the block", *ns)
	}
}

// TestAgentBlockedThenExit ends the session while it is blocked: the
// unblocked entry comes before the ended one.
func TestAgentBlockedThenExit(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	ns := noticeLog(e)
	r.agentSeq["S-build"] = []AgentStatus{blocked, {Done: true, Outcome: "pass", Detail: "exit 0", Blocked: question}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || st.Blocked != nil || st.Steps[0].Outcome != "pass" {
		t.Fatalf("result %q blocked %+v steps %+v", st.Result, st.Blocked, st.Steps)
	}
	want := []string{
		"started:build", "entered:build", "session:build", "blocked:build",
		"unblocked:build", "ended:build", "finished:build",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	if len(*ns) != 2 || !(*ns)[0].Blocked || (*ns)[1].Blocked {
		t.Errorf("notices = %+v", *ns)
	}
}

// TestAgentBlockedThenRunnerError ends the step with an error outcome
// while blocked: the block still ends before the step does.
func TestAgentBlockedThenRunnerError(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	ns := noticeLog(e)
	r.agentSeq["S-build"] = []AgentStatus{blocked}
	r.agentErr["S-build"] = errors.New("session record unreadable")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || st.Detail != "session record unreadable" || st.Blocked != nil {
		t.Fatalf("result %q detail %q blocked %+v", st.Result, st.Detail, st.Blocked)
	}
	want := []string{
		"started:build", "entered:build", "session:build", "blocked:build",
		"unblocked:build", "ended:build", "finished:build",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	if len(*ns) != 2 || (*ns)[1].Blocked {
		t.Errorf("notices = %+v", *ns)
	}
}

// TestAgentSecondBlockNotifiesAgain blocks, answers and blocks again in
// one step: each block is journalled and notified.
func TestAgentSecondBlockNotifiesAgain(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	ns := noticeLog(e)
	permission := &Block{Kind: "permission", Question: "Run rm?", TabTitle: "build"}
	r.agentSeq["S-build"] = []AgentStatus{blocked, idle, {Blocked: permission}, {Blocked: permission}}
	r.outcome["S-build"] = "pass"
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"started:build", "entered:build", "session:build", "blocked:build", "unblocked:build",
		"blocked:build", "unblocked:build", "ended:build", "finished:build",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	var got []string
	for _, n := range *ns {
		got = append(got, fmt.Sprintf("%v:%s", n.Blocked, n.Block.Question))
	}
	if want := []string{"true:Ship it?", "false:Ship it?", "true:Run rm?", "false:Run rm?"}; !reflect.DeepEqual(got, want) {
		t.Errorf("notices = %v, want %v", got, want)
	}
}

// TestAgentBlockWithoutNotify journals the block when no Notify is set.
func TestAgentBlockWithoutNotify(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	r.agentSeq["S-build"] = []AgentStatus{blocked, idle}
	r.outcome["S-build"] = "pass"
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if countType(t, e, TypeBlocked) != 1 || countType(t, e, TypeUnblocked) != 1 {
		t.Errorf("journal = %v", entryTypes(t, e))
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

const agentUntilWF = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build {change}"
until = "pr-open"
interval = "2s"
on.pass = "done"
`

// An idle session whose until holds passes the step and is left open: the
// engine never waits for its exit.
func TestAgentUntilHoldsWhenIdle(t *testing.T) {
	e, r, _, _ := setup(t, agentUntilWF)
	r.agentSeq["S-build"] = []AgentStatus{{}, {Idle: true}}
	r.agentErr["S-build"] = errors.New("polled past the pass")
	r.holds["build"] = []bool{true}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["S-build"] != 2 || r.polls["build"] != 1 {
		t.Fatalf("result %q detail %q polls %v", st.Result, st.Detail, r.polls)
	}
	if got := st.Steps[0].Detail; got != "idle and until held; session left open" {
		t.Errorf("detail = %q", got)
	}
}

// An idle session whose until does not hold keeps the step waiting, with
// the condition checked once per idle poll; it still ends on the exit.
func TestAgentUntilFailsKeepsWaiting(t *testing.T) {
	e, r, c, _ := setup(t, agentUntilWF)
	r.agentSeq["S-build"] = []AgentStatus{{Idle: true}, {Idle: true}, {}}
	r.holds["build"] = []bool{false}
	r.outcome["S-build"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["S-build"] != 4 || r.polls["build"] != 2 {
		t.Fatalf("result %q polls %v", st.Result, r.polls)
	}
	if st.Steps[0].Detail != "exit 0" || len(c.sleeps) != 3 {
		t.Errorf("detail %q sleeps %v", st.Steps[0].Detail, c.sleeps)
	}
}

// A blocked session is never passed over, even when until holds and the
// session also reads idle.
func TestAgentUntilNotCheckedWhileBlocked(t *testing.T) {
	e, r, _, _ := setup(t, agentUntilWF)
	b := &Block{Kind: "question", Question: "Ship it?"}
	r.agentSeq["S-build"] = []AgentStatus{{Blocked: b}, {Idle: true, Blocked: b}}
	r.holds["build"] = []bool{true}
	r.outcome["S-build"] = "fail"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || r.polls["build"] != 0 {
		t.Fatalf("result %q polls %v: until was checked while blocked", st.Result, r.polls)
	}
}

// A runner error from until ends the step with the error outcome.
func TestAgentUntilErrorFailsRun(t *testing.T) {
	e, r, _, _ := setup(t, agentUntilWF)
	r.agentSeq["S-build"] = []AgentStatus{{Idle: true}}
	// Only until fails: the session has started by the first poll.
	r.onPoll = func(string, int) { r.failOn["build"] = errors.New("sh missing") }
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || r.starts["build"] != 1 || !strings.Contains(st.Steps[0].Detail, "sh missing") {
		t.Fatalf("result %q steps %+v", st.Result, st.Steps)
	}
}

const agentDoneWF = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build {change}"
until = "pr-open"
until_when = "done"
interval = "2s"
on.pass = "done"
`

// A step that waits for done never checks until while the session sits
// idle, however long; once the agent says done and until holds, it passes
// and leaves the session open.
func TestAgentUntilDoneIgnoresIdle(t *testing.T) {
	e, r, c, _ := setup(t, agentDoneWF)
	r.agentSeq["S-build"] = []AgentStatus{{Idle: true}, {Idle: true}, {Idle: true}, {Idle: true, SaidDone: true}}
	r.agentErr["S-build"] = errors.New("polled past the pass")
	r.holds["build"] = []bool{true}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["S-build"] != 4 || r.polls["build"] != 1 || len(c.sleeps) != 3 {
		t.Fatalf("result %q detail %q polls %v sleeps %v", st.Result, st.Detail, r.polls, c.sleeps)
	}
	if got := st.Steps[0].Detail; got != "said done and until held; session left open" {
		t.Errorf("detail = %q", got)
	}
}

// A session that said done but whose until
// does not hold keeps the step waiting, checking until once per poll.
func TestAgentUntilDoneKeepsWaitingUntilHolds(t *testing.T) {
	e, r, _, _ := setup(t, agentDoneWF)
	r.agentSeq["S-build"] = []AgentStatus{{SaidDone: true}, {SaidDone: true}}
	r.agentErr["S-build"] = errors.New("polled past the pass")
	r.holds["build"] = []bool{false, true}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["S-build"] != 2 || r.polls["build"] != 2 {
		t.Fatalf("result %q polls %v", st.Result, r.polls)
	}
}

// A session that said done but is blocked on the owner is not passed over.
func TestAgentUntilDoneNotCheckedWhileBlocked(t *testing.T) {
	e, r, _, _ := setup(t, agentDoneWF)
	b := &Block{Kind: "permission"}
	r.agentSeq["S-build"] = []AgentStatus{{SaidDone: true, Blocked: b}}
	r.holds["build"] = []bool{true}
	r.outcome["S-build"] = "fail"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || r.polls["build"] != 0 {
		t.Fatalf("result %q polls %v: until was checked while blocked", st.Result, r.polls)
	}
}

// When the session of a step that waits for done exits 0, until is
// checked: the step passes only if it holds. A non-zero exit fails without
// checking it.
func TestAgentUntilDoneOnExit(t *testing.T) {
	for _, tt := range []struct {
		name    string
		outcome string
		holds   []bool
		result  string
		detail  string
		checks  int
	}{
		{"exit 0, until holds", "pass", []bool{true}, Succeeded, "exit 0 and until held", 1},
		{"exit 0, until does not hold", "pass", []bool{false}, Failed, "exit 0 but until did not hold", 1},
		{"non-zero exit", "fail", []bool{true}, Failed, "exit 0", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e, r, _, _ := setup(t, agentDoneWF)
			r.agentSeq["S-build"] = []AgentStatus{{Idle: true}}
			r.outcome["S-build"] = tt.outcome
			r.holds["build"] = tt.holds
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if st.Result != tt.result || st.Steps[0].Detail != tt.detail || r.polls["build"] != tt.checks {
				t.Fatalf("result %q detail %q polls %v", st.Result, st.Steps[0].Detail, r.polls)
			}
		})
	}
}

// A step without until_when keeps today's behaviour on exit: exit 0
// passes without checking until.
func TestAgentUntilIdleExitSkipsUntil(t *testing.T) {
	e, r, _, _ := setup(t, agentUntilWF)
	r.outcome["S-build"] = "pass"
	r.holds["build"] = []bool{false}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["build"] != 0 || st.Steps[0].Detail != "exit 0" {
		t.Fatalf("result %q detail %q polls %v", st.Result, st.Steps[0].Detail, r.polls)
	}
}

// A run whose command step failed is retried at that step: the step runs
// again and the run succeeds. The journal keeps the failure, then the
// retried entry, then the new visit.
func TestRetryRerunsFailedStep(t *testing.T) {
	e, r, _, _ := setup(t, `
name = "wf"
[step.one]
kind = "command"
run = "true"
on.pass = "two"

[step.two]
kind = "command"
run = "test -f x"
on.pass = "done"
on.fail = "failed"
`)
	r.codes["two"] = []int{1, 0}
	st, err := e.Run(runID)
	if err != nil || st.Result != Failed {
		t.Fatalf("first run: %v, %q", err, st.Result)
	}
	if _, err := e.Run(runID); !errors.Is(err, ErrFinished) {
		t.Fatalf("plain resume of a failed run: %v", err)
	}
	st, err = e.Retry(runID[:6])
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("result after retry = %q (%s)", st.Result, st.Detail)
	}
	if r.commands["one"] != 1 || r.commands["two"] != 2 {
		t.Errorf("commands = %v", r.commands)
	}
	if got, want := outcomes(st), []string{"one:pass", "two:fail", "two:pass"}; !reflect.DeepEqual(got, want) {
		t.Errorf("outcomes = %v, want %v", got, want)
	}
	if last := st.Steps[2]; !last.Retried || last.Visit != 2 {
		t.Errorf("retried visit = %+v", last)
	}
	want := []string{
		"started:one", "entered:one", "ended:one", "entered:two", "ended:two", "finished:two",
		"retried:two", "entered:two", "ended:two", "finished:two",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("entries = %v\nwant %v", got, want)
	}
}

// A run that went over a step's max_visits is retried at that step: the
// cap counts only visits since the retry, so the step is entered again.
func TestRetryOverCapStep(t *testing.T) {
	e, r, _, _ := setup(t, loopWF)
	r.codes["test"] = []int{1, 1, 0}
	st, err := e.Run(runID)
	if err != nil || st.Result != Failed || st.FinishedStep != "build" {
		t.Fatalf("first run: %v, %+v", err, st)
	}
	st, err = e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("result after retry = %q (%s)", st.Result, st.Detail)
	}
	if r.commands["build"] != 3 || r.commands["test"] != 3 {
		t.Errorf("commands = %v", r.commands)
	}
	if st.Visits["build"] != 3 || st.CapVisits("build") != 1 {
		t.Errorf("build visits %d, cap visits %d", st.Visits["build"], st.CapVisits("build"))
	}
}

// Retrying a run that succeeded is refused and writes nothing.
func TestRetrySucceededIsErrNotFailed(t *testing.T) {
	e, r, _, _ := setup(t, cmdWF)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	before := entryTypes(t, e)
	_, err := e.Retry(runID)
	if !errors.Is(err, ErrNotFailed) || !strings.Contains(err.Error(), "only a failed run can be retried") {
		t.Fatalf("err = %v, want ErrNotFailed", err)
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, before) {
		t.Errorf("entries changed: %v", got)
	}
	if r.commands["test"] != 1 {
		t.Errorf("test ran %d times", r.commands["test"])
	}
}

// Retry on a run that is not finished runs it as Run does, with no
// retried entry.
func TestRetryUnfinishedRunResumes(t *testing.T) {
	e, r, _, _ := setup(t, waitWF)
	e.NoWait = true
	if st, err := e.Run(runID); err != nil || st.Stopped == "" {
		t.Fatalf("first run: %v, %+v", err, st)
	}
	e.NoWait = false
	r.holds["ci"] = []bool{true}
	st, err := e.Retry(runID)
	if err != nil || st.Result != Succeeded {
		t.Fatalf("retry of a stopped run: %v, %+v", err, st)
	}
	if countType(t, e, TypeRetried) != 0 || len(st.Retries) != 0 {
		t.Errorf("a stopped run was journalled as retried")
	}
}

// Decide counts max_visits from the last retry.
func TestDecideCapSinceRetry(t *testing.T) {
	def, err := workflow.Parse("wf", []byte(loopWF))
	if err != nil {
		t.Fatal(err)
	}
	st := State{Current: "build", Visits: map[string]int{"build": 2}, RetryBase: map[string]int{"build": 2}}
	got, err := Decide(def, st)
	if err != nil {
		t.Fatal(err)
	}
	if want := (Action{Kind: ActEnter, Step: "build", Visit: 3}); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	st.Visits["build"], st.RetryBase["build"] = 4, 2
	if got, _ := Decide(def, st); got.Kind != ActOverCap || got.Visit != 5 {
		t.Fatalf("over cap since retry: %+v", got)
	}
}

// arbiterWF has an agent step and a declared arbiter step.
const arbiterWF = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build"
[step.decide]
kind = "arbiter"
interval = "3s"
`

// arbitrateOn runs the arbiter of def on question Q1 as asked from step
// build, without taking the run's lock, and returns its result.
func arbitrateOn(t *testing.T, e *Engine, def *workflow.Definition) result {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	r := &running{e: e, j: j, def: def}
	res, err := r.arbitrate(def.Steps["build"], "Q1")
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// arbitrations reads the run's folded arbitrations.
func arbitrations(t *testing.T, e *Engine) []Arbitration {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	return j.State().Arbitrations
}

// TestArbitratePasses polls the arbiter's session at the arbiter step's
// interval until it says done, then passes on a valid answer, journalling
// arbiter and arbitrated entries.
func TestArbitratePasses(t *testing.T) {
	e, r, c, def := setup(t, arbiterWF)
	r.agentSeq["S-decide-Q1"] = []AgentStatus{idle, {SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "assume-sqlite", Statement: "Use SQLite"}
	res := arbitrateOn(t, e, def)
	if res.outcome != workflow.OutcomePass || res.detail != "answered by assume-sqlite" {
		t.Fatalf("result %+v", res)
	}
	if !reflect.DeepEqual(r.arbiters, []string{"decide Q1"}) || r.answerN["Q1"] != 1 {
		t.Errorf("arbiters %v answers %v", r.arbiters, r.answerN)
	}
	if want := []string{"RHIZA_CHANGE=CS1", "RHIZA_RUN=" + runID, "RHIZA_STEP=decide"}; !reflect.DeepEqual(r.envs["decide"], want) {
		t.Errorf("arbiter env %v, want %v", r.envs["decide"], want)
	}
	if want := []time.Duration{3 * time.Second}; !reflect.DeepEqual(c.sleeps, want) {
		t.Errorf("sleeps = %v, want %v", c.sleeps, want)
	}
	if got, want := entryTypes(t, e), []string{"started:build", "arbiter:build", "arbitrated:build"}; !reflect.DeepEqual(got, want) {
		t.Errorf("journal %v, want %v", got, want)
	}
	want := []Arbitration{{Step: "build", Session: "S-decide-Q1", Question: "Q1", Outcome: "pass",
		Detail: "answered by assume-sqlite", Decision: "assume-sqlite", At: at(1)}}
	if got := arbitrations(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("arbitrations %+v\nwant %+v", got, want)
	}
}

// TestArbitrateFailsWithoutAnswer fails when the arbiter wrote nothing,
// or a Decision the Runner does not accept (no reasoning), with the
// Runner's why; an exited session is judged the same as one said done.
func TestArbitrateFailsWithoutAnswer(t *testing.T) {
	for _, tt := range []struct {
		name   string
		answer *Answer
		why    string
	}{
		{"NoDecision", nil, "no live Decision answers question Q1"},
		{"NoReasoning", &Answer{Decision: "assume-sqlite", Why: "decision assume-sqlite has no reasoning"},
			"decision assume-sqlite has no reasoning"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e, r, _, def := setup(t, arbiterWF)
			if tt.answer != nil {
				r.answers["Q1"] = *tt.answer
			}
			res := arbitrateOn(t, e, def)
			if res.outcome != workflow.OutcomeFail || res.detail != tt.why {
				t.Fatalf("result %+v, want fail %q", res, tt.why)
			}
			a := arbitrations(t, e)
			if len(a) != 1 || a[0].Outcome != "fail" || a[0].Detail != tt.why || a[0].Decision != "" || a[0].Open() {
				t.Errorf("arbitrations %+v", a)
			}
		})
	}
}

// TestArbitrateFailsWhenArbiterAsks fails as soon as the arbiter's own
// session blocks on a question, without asking for the answer (D10).
func TestArbitrateFailsWhenArbiterAsks(t *testing.T) {
	e, r, _, def := setup(t, arbiterWF)
	ns := noticeLog(e)
	r.agentSeq["S-decide-Q1"] = []AgentStatus{idle, blocked, {SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "assume-sqlite"}
	res := arbitrateOn(t, e, def)
	if res.outcome != workflow.OutcomeFail || res.detail != "the arbiter asked a question: Ship it?" {
		t.Fatalf("result %+v", res)
	}
	if r.answerN["Q1"] != 0 || r.polls["S-decide-Q1"] != 2 {
		t.Errorf("answer calls %d, polls %d", r.answerN["Q1"], r.polls["S-decide-Q1"])
	}
	if countType(t, e, TypeBlocked) != 0 || len(*ns) != 0 {
		t.Errorf("the arbiter's block was journalled or notified: %v", *ns)
	}
}

// TestArbitrateResumeRejoins rejoins the arbiter session an open
// arbitration names instead of launching a second one.
func TestArbitrateResumeRejoins(t *testing.T) {
	e, r, _, def := setup(t, arbiterWF)
	if err := e.Store.Append(runID, Entry{Type: TypeArbiter, At: at(1), Step: "build", Session: "S-earlier", Question: "Q1"}); err != nil {
		t.Fatal(err)
	}
	r.agentSeq["S-earlier"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "assume-sqlite"}
	res := arbitrateOn(t, e, def)
	if res.outcome != workflow.OutcomePass || len(r.arbiters) != 0 || r.polls["S-earlier"] != 1 {
		t.Fatalf("result %+v arbiters %v polls %v", res, r.arbiters, r.polls)
	}
	if got := countType(t, e, TypeArbiter); got != 1 {
		t.Errorf("arbiter entries = %d, want 1", got)
	}
	a := arbitrations(t, e)
	if len(a) != 1 || a[0].Session != "S-earlier" || a[0].Outcome != "pass" {
		t.Errorf("arbitrations %+v", a)
	}
	// A closed arbitration is not rejoined: the next one launches anew.
	if _, ok := (Journal{ID: runID, Entries: mustEntries(t, e)}).State().OpenArbitration("Q1"); ok {
		t.Error("a judged arbitration is still open")
	}
}

// TestArbitrateBuiltinAndStartError uses the built-in arbiter when the
// file declares none, and journals a Runner failure to launch it as an
// error verdict.
func TestArbitrateBuiltinAndStartError(t *testing.T) {
	e, r, c, def := setup(t, agentWF)
	r.agentSeq["S-arbiter-Q1"] = []AgentStatus{idle, {Done: true, Outcome: "pass"}}
	r.answers["Q1"] = Answer{OK: true, Decision: "d"}
	if res := arbitrateOn(t, e, def); res.outcome != workflow.OutcomePass {
		t.Fatalf("result %+v", res)
	}
	if !reflect.DeepEqual(r.arbiters, []string{"arbiter Q1"}) || !reflect.DeepEqual(c.sleeps, []time.Duration{workflow.DefaultInterval}) {
		t.Errorf("arbiters %v sleeps %v", r.arbiters, c.sleeps)
	}

	e, r, _, def = setup(t, agentWF)
	r.failOn["arbiter"] = errors.New("adapter refused")
	res := arbitrateOn(t, e, def)
	if res.outcome != OutcomeError || res.detail != "adapter refused" {
		t.Fatalf("result %+v", res)
	}
	a := arbitrations(t, e)
	if len(a) != 1 || a[0].Session != "" || a[0].Outcome != OutcomeError || countType(t, e, TypeArbiter) != 0 {
		t.Errorf("arbitrations %+v", a)
	}
}

func mustEntries(t *testing.T, e *Engine) []Entry {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	return j.Entries
}

// autonomousSetup is setup for an autonomous run of src.
func autonomousSetup(t *testing.T, src string) (*Engine, *fakeRunner) {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	store := newStore(t)
	if err := Start(store, runID, Started{Workflow: "wf", Definition: src, Change: "CS1", Title: "Blocked runs", Autonomous: true}, t0, def); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	return &Engine{Store: store, Runner: r, Clock: c, Now: func() time.Time { return at(1) }}, r
}

// TestAutonomousQuestionGoesToArbiter: an autonomous run's session blocks
// on a question; the engine adds the question once, journals asked, has
// the arbiter answer it, polls on while the session is still blocked
// without asking again, and the session's exit decides the step. The
// owner is never notified.
func TestAutonomousQuestionGoesToArbiter(t *testing.T) {
	e, r := autonomousSetup(t, arbiterWF)
	ns := noticeLog(e)
	r.agentSeq["S-build"] = []AgentStatus{idle, blocked, idle, idle}
	r.outcome["S-build"] = "pass"
	r.agentSeq["S-decide-Q1"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "assume-sqlite", Statement: "Use SQLite."}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"build:pass"}) {
		t.Fatalf("result %q steps %v detail %q", st.Result, outcomes(st), st.Detail)
	}
	if !reflect.DeepEqual(r.asked, []string{"S-build Ship it?"}) || !reflect.DeepEqual(r.arbiters, []string{"decide Q1"}) {
		t.Errorf("asked %v arbiters %v", r.asked, r.arbiters)
	}
	want := []string{
		"started:build", "entered:build", "session:build", "blocked:build", "asked:build",
		"arbiter:build", "arbitrated:build", "answered:build", "unblocked:build", "ended:build", "finished:build",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	if !reflect.DeepEqual(r.delivered, []string{"S-build Q1 assume-sqlite Use SQLite."}) {
		t.Errorf("delivered %q", r.delivered)
	}
	wantAsked := []Asked{{Step: "build", Session: "S-build", Question: "Q1", Handle: "question-s-1", Text: "Ship it?", At: at(1), Answered: "assume-sqlite"}}
	if !reflect.DeepEqual(st.Asked, wantAsked) {
		t.Errorf("asked %+v\nwant %+v", st.Asked, wantAsked)
	}
	if len(*ns) != 0 {
		t.Errorf("notices = %+v, want none", *ns)
	}
}

// TestAutonomousArbiterFailFailsRun: the arbiter fails, so the step ends
// with an error (its on.fail is not followed) and the run fails naming
// the question.
func TestAutonomousArbiterFailFailsRun(t *testing.T) {
	src := strings.Replace(arbiterWF, "[step.decide]", "on.fail = \"done\"\n[step.decide]", 1)
	e, r := autonomousSetup(t, src)
	r.agentSeq["S-build"] = []AgentStatus{blocked}
	r.agentSeq["S-decide-Q1"] = []AgentStatus{{SaidDone: true}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	wantDetail := "arbiter failed on question question-s-1 (Q1): no live Decision answers question Q1"
	if st.Result != Failed || st.Detail != wantDetail || !reflect.DeepEqual(outcomes(st), []string{"build:error"}) {
		t.Fatalf("result %q detail %q steps %v", st.Result, st.Detail, outcomes(st))
	}
	if r.polls["S-build"] != 1 || countType(t, e, TypeUnblocked) != 1 {
		t.Errorf("polls %v unblocked %d", r.polls, countType(t, e, TypeUnblocked))
	}

	e, r = autonomousSetup(t, arbiterWF)
	r.agentSeq["S-build"] = []AgentStatus{blocked}
	r.askErr = errors.New("change set CS1 is landed")
	st, err = e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || !strings.Contains(st.Detail, "change set CS1 is landed") || len(r.arbiters) != 0 {
		t.Fatalf("result %q detail %q arbiters %v", st.Result, st.Detail, r.arbiters)
	}
}

// TestNormalRunQuestionNotArbitrated: without autonomous, the question
// is the owner's: no question node, no arbiter, the owner is notified.
func TestNormalRunQuestionNotArbitrated(t *testing.T) {
	e, r, _, _ := setup(t, arbiterWF)
	ns := noticeLog(e)
	r.agentSeq["S-build"] = []AgentStatus{blocked, blocked, idle}
	r.outcome["S-build"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.asked) != 0 || len(r.arbiters) != 0 || len(st.Asked) != 0 {
		t.Fatalf("result %q asked %v arbiters %v", st.Result, r.asked, r.arbiters)
	}
	if len(*ns) != 2 || countType(t, e, TypeAsked) != 0 {
		t.Errorf("notices %+v asked entries %d", *ns, countType(t, e, TypeAsked))
	}
}

// TestAutonomousOtherBlockNotArbitrated: a block that is not a question
// is not arbitrated, and the owner is still told of it.
func TestAutonomousOtherBlockNotArbitrated(t *testing.T) {
	e, r := autonomousSetup(t, arbiterWF)
	ns := noticeLog(e)
	r.agentSeq["S-build"] = []AgentStatus{{Blocked: &Block{Kind: "permission"}}, idle}
	r.outcome["S-build"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.asked) != 0 || len(r.arbiters) != 0 || len(*ns) != 2 {
		t.Fatalf("result %q asked %v arbiters %v notices %d", st.Result, r.asked, r.arbiters, len(*ns))
	}
}

// TestAutonomousResumeRejoinsArbiter resumes a step whose journal shows
// the block asked and its arbitration open: no second question, no
// second arbiter; the open one is rejoined. A second block of the same
// session is asked anew.
func TestAutonomousResumeRejoinsArbiter(t *testing.T) {
	e, r := autonomousSetup(t, arbiterWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "build", Visit: 1},
		{Type: TypeSession, At: at(1), Step: "build", Session: "S-earlier"},
		{Type: TypeBlocked, At: at(1), Step: "build", Session: "S-earlier", Block: question},
		{Type: TypeAsked, At: at(1), Step: "build", Session: "S-earlier", Question: "Q9", Handle: "question-s-9", Detail: "Ship it?"},
		{Type: TypeArbiter, At: at(1), Step: "build", Session: "S-arb", Question: "Q9"},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	r.agentSeq["S-earlier"] = []AgentStatus{blocked, idle, blocked}
	r.outcome["S-earlier"] = "pass"
	r.agentSeq["S-arb"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q9"] = Answer{OK: true, Decision: "d9"}
	r.agentSeq["S-decide-Q1"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "d1"}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || r.polls["S-arb"] != 1 {
		t.Fatalf("result %q detail %q polls %v", st.Result, st.Detail, r.polls)
	}
	// Only the second block was asked and launched an arbiter.
	if !reflect.DeepEqual(r.asked, []string{"S-earlier Ship it?"}) || !reflect.DeepEqual(r.arbiters, []string{"decide Q1"}) {
		t.Errorf("asked %v arbiters %v", r.asked, r.arbiters)
	}
	if countType(t, e, TypeAsked) != 2 || countType(t, e, TypeArbitrated) != 2 {
		t.Errorf("journal %v", entryTypes(t, e))
	}
}

// TestAutonomousAnswerDelivered: once the arbiter passes, the engine
// delivers its Decision to the asking session once, journals answered
// then unblocked, and keeps polling: the session's own outcome (here a
// fail) ends the step.
func TestAutonomousAnswerDelivered(t *testing.T) {
	e, r := autonomousSetup(t, arbiterWF)
	r.agentSeq["S-build"] = []AgentStatus{blocked, idle, idle}
	r.outcome["S-build"] = "fail"
	r.agentSeq["S-decide-Q1"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "d1", Statement: "Ship it."}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || !reflect.DeepEqual(outcomes(st), []string{"build:fail"}) {
		t.Fatalf("result %q steps %v detail %q", st.Result, outcomes(st), st.Detail)
	}
	if !reflect.DeepEqual(r.delivered, []string{"S-build Q1 d1 Ship it."}) || r.polls["S-build"] != 4 {
		t.Errorf("delivered %q polls %d", r.delivered, r.polls["S-build"])
	}
	var answered []Entry
	for _, en := range mustEntries(t, e) {
		if en.Type == TypeAnswered {
			answered = append(answered, en)
		}
	}
	want := []Entry{{Type: TypeAnswered, At: at(1), Step: "build", Session: "S-build", Question: "Q1", Decision: "d1"}}
	if !reflect.DeepEqual(answered, want) {
		t.Errorf("answered %+v\nwant %+v", answered, want)
	}
	if len(st.Asked) != 1 || st.Asked[0].Answered != "d1" {
		t.Errorf("asked %+v", st.Asked)
	}
}

// TestAutonomousDeliverErrorFailsRun: a delivery that fails ends the
// step with an error naming the question, its block ended first.
func TestAutonomousDeliverErrorFailsRun(t *testing.T) {
	e, r := autonomousSetup(t, arbiterWF)
	r.agentSeq["S-build"] = []AgentStatus{blocked}
	r.agentSeq["S-decide-Q1"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q1"] = Answer{OK: true, Decision: "d1", Statement: "Ship it."}
	r.deliverErr = errors.New("session S-build exited")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	wantDetail := "could not deliver the answer to question question-s-1 (Q1): session S-build exited"
	if st.Result != Failed || st.Detail != wantDetail || !reflect.DeepEqual(outcomes(st), []string{"build:error"}) {
		t.Fatalf("result %q detail %q steps %v", st.Result, st.Detail, outcomes(st))
	}
	if countType(t, e, TypeAnswered) != 0 || countType(t, e, TypeUnblocked) != 1 {
		t.Errorf("journal %v", entryTypes(t, e))
	}
}

// TestAutonomousResumeRedelivers: a resume whose journal has a passing
// arbitration and no answered entry delivers again; one whose journal has
// the answered entry only ends the block.
func TestAutonomousResumeRedelivers(t *testing.T) {
	for _, c := range []struct {
		name      string
		answered  bool
		delivered int
	}{{"no answered", false, 1}, {"answered", true, 0}} {
		t.Run(c.name, func(t *testing.T) {
			e, r := autonomousSetup(t, arbiterWF)
			entries := []Entry{
				{Type: TypeEntered, At: at(1), Step: "build", Visit: 1},
				{Type: TypeSession, At: at(1), Step: "build", Session: "S-earlier"},
				{Type: TypeBlocked, At: at(1), Step: "build", Session: "S-earlier", Block: question},
				{Type: TypeAsked, At: at(1), Step: "build", Session: "S-earlier", Question: "Q9", Handle: "question-s-9", Detail: "Ship it?"},
				{Type: TypeArbiter, At: at(1), Step: "build", Session: "S-arb", Question: "Q9"},
				{Type: TypeArbitrated, At: at(1), Step: "build", Question: "Q9", Outcome: "pass", Detail: "answered by d9", Decision: "d9"},
			}
			if c.answered {
				entries = append(entries, Entry{Type: TypeAnswered, At: at(1), Step: "build", Session: "S-earlier", Question: "Q9", Decision: "d9"})
			}
			for _, en := range entries {
				if err := e.Store.Append(runID, en); err != nil {
					t.Fatal(err)
				}
			}
			r.agentSeq["S-earlier"] = []AgentStatus{blocked, idle}
			r.outcome["S-earlier"] = "pass"
			r.answers["Q9"] = Answer{OK: true, Decision: "d9", Statement: "Yes."}
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if st.Result != Succeeded || len(r.delivered) != c.delivered || len(r.arbiters) != 0 || len(r.asked) != 0 {
				t.Fatalf("result %q delivered %q arbiters %v asked %v", st.Result, r.delivered, r.arbiters, r.asked)
			}
			if countType(t, e, TypeAnswered) != 1 || countType(t, e, TypeUnblocked) != 1 {
				t.Errorf("journal %v", entryTypes(t, e))
			}
		})
	}
}

// Each agent step that ends with a session reports it to the Runner
// exactly once, whatever its outcome: an exit, a pass that leaves the
// session open, a fail, or an error after the session started.
func TestAgentEndedOncePerEndedStep(t *testing.T) {
	for _, tt := range []struct {
		name   string
		wf     string
		script func(r *fakeRunner)
		result string
	}{
		{"exit", agentWF, func(r *fakeRunner) {
			r.agentOut["S-build"] = []bool{false, false}
			r.outcome["S-build"] = "pass"
		}, Succeeded},
		{"fail", agentWF, func(r *fakeRunner) { r.outcome["S-build"] = "fail" }, Failed},
		{"said done, left open", agentDoneWF, func(r *fakeRunner) {
			r.agentSeq["S-build"] = []AgentStatus{{Idle: true}, {Idle: true, SaidDone: true}}
			r.agentErr["S-build"] = errors.New("polled past the pass")
			r.holds["build"] = []bool{true}
		}, Succeeded},
		{"idle, left open", agentUntilWF, func(r *fakeRunner) {
			r.agentSeq["S-build"] = []AgentStatus{{Idle: true}}
			r.agentErr["S-build"] = errors.New("polled past the pass")
			r.holds["build"] = []bool{true}
		}, Succeeded},
		{"status error", agentWF, func(r *fakeRunner) {
			r.agentSeq["S-build"] = []AgentStatus{{}, blocked}
			r.agentErr["S-build"] = errors.New("record unreadable")
		}, Failed},
		{"until error", agentUntilWF, func(r *fakeRunner) {
			r.agentSeq["S-build"] = []AgentStatus{{Idle: true}}
			r.onPoll = func(string, int) { r.failOn["build"] = errors.New("sh missing") }
		}, Failed},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e, r, _, _ := setup(t, tt.wf)
			tt.script(r)
			st, err := e.Run(runID)
			if err != nil {
				t.Fatal(err)
			}
			if st.Result != tt.result {
				t.Fatalf("result %q detail %q", st.Result, st.Detail)
			}
			if want := []string{"build:S-build"}; !reflect.DeepEqual(r.ended, want) {
				t.Errorf("AgentEnded calls = %v, want %v", r.ended, want)
			}
		})
	}
}

// A resumed agent step reports the session it rejoined.
func TestAgentEndedOnResume(t *testing.T) {
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
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if want := []string{"build:S-earlier"}; !reflect.DeepEqual(r.ended, want) {
		t.Errorf("AgentEnded calls = %v, want %v", r.ended, want)
	}
}

// An agent step whose session never started reports nothing.
func TestAgentEndedNotWithoutSession(t *testing.T) {
	e, r, _, _ := setup(t, agentWF)
	r.failOn["build"] = errors.New("adapter refused")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Failed || len(r.ended) != 0 {
		t.Fatalf("result %q AgentEnded calls %v", st.Result, r.ended)
	}
}

// A run that stops reports nothing for the stopped step: only the agent
// step that ended before it is reported, and command steps never are.
func TestAgentEndedNotForStoppedStep(t *testing.T) {
	const wf = `
name = "wf"
[step.build]
kind = "agent"
prompt = "build {change}"
interval = "2s"
on.pass = "lint"
[step.lint]
kind = "command"
run = "lint"
on.pass = "ci"
[step.ci]
kind = "wait"
until = "gh pr checks"
interval = "1s"
timeout = "5s"
on.pass = "done"
`
	e, r, _, _ := setup(t, wf)
	e.NoWait = true
	r.outcome["S-build"] = "pass"
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Stopped == "" {
		t.Fatalf("run did not stop: %+v", st)
	}
	if want := []string{"build:S-build"}; !reflect.DeepEqual(r.ended, want) {
		t.Errorf("AgentEnded calls = %v, want %v", r.ended, want)
	}
}

// An AgentEnded error is logged and changes nothing: the step's outcome,
// its detail and the journal are those of a run where it succeeded.
func TestAgentEndedErrorChangesNothing(t *testing.T) {
	run := func(endedErr error) (State, []string, []string) {
		e, r, _, _ := setup(t, agentWF)
		var logs []string
		e.Log = func(l string) { logs = append(logs, l) }
		r.outcome["S-build"] = "pass"
		r.endedErr = endedErr
		st, err := e.Run(runID)
		if err != nil {
			t.Fatal(err)
		}
		return st, entryTypes(t, e), logs
	}
	okSt, okTypes, _ := run(nil)
	st, types, logs := run(errors.New("record token use: no transcript"))
	if st.Result != okSt.Result || !reflect.DeepEqual(outcomes(st), outcomes(okSt)) ||
		st.Steps[0].Detail != okSt.Steps[0].Detail || !reflect.DeepEqual(types, okTypes) {
		t.Fatalf("with an error: %q %v %v; without: %q %v %v", st.Result, outcomes(st), types, okSt.Result, outcomes(okSt), okTypes)
	}
	if !slices.Contains(logs, "step build: record token use: no transcript") {
		t.Errorf("logs = %q", logs)
	}
}

const askedWF = `
name = "wf"
[step.check]
kind = "command"
run = "rhiza run review-check"
on.pass = "next"
on.fail = "next"
[step.next]
kind = "command"
run = "true"
[step.decide]
kind = "arbiter"
interval = "3s"
`

// journalAsked makes the check command append, lock-free as a child
// process does, a session-less asked entry for question node q.
func journalAsked(t *testing.T, e *Engine, r *fakeRunner, q string) {
	t.Helper()
	r.onCommand = func(step string) {
		if step != "check" {
			return
		}
		if err := e.Store.Append(runID, Entry{
			Type: TypeAsked, At: at(1), Step: "check", Question: q, Handle: "review-x-r3-f1", Detail: "How?",
		}); err != nil {
			t.Fatal(err)
		}
	}
}

// TestCommandAskedArbitrated: a command step journals a session-less asked
// entry; in an autonomous run the arbiter decides it right after the step,
// with no delivery, and the run goes on to its next step.
func TestCommandAskedArbitrated(t *testing.T) {
	e, r := autonomousSetup(t, askedWF)
	journalAsked(t, e, r, "Q9")
	r.agentSeq["S-decide-Q9"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q9"] = Answer{OK: true, Decision: "assume-x", Statement: "Do x."}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"check:pass", "next:pass"}) {
		t.Fatalf("result %q steps %v detail %q", st.Result, outcomes(st), st.Detail)
	}
	want := []string{
		"started:check", "entered:check", "asked:check", "arbiter:check", "arbitrated:check",
		"ended:check", "entered:next", "ended:next", "finished:next",
	}
	if got := entryTypes(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journal = %v\nwant      %v", got, want)
	}
	if !reflect.DeepEqual(r.arbiters, []string{"decide Q9"}) || len(r.delivered) != 0 {
		t.Errorf("arbiters %v delivered %v", r.arbiters, r.delivered)
	}
	if a, ok := st.Arbitration("Q9"); !ok || a.Outcome != "pass" || a.Decision != "assume-x" {
		t.Errorf("arbitration %+v", a)
	}
}

// TestCommandAskedArbiterFailFailsRun: a failed arbiter fails the run,
// whatever the step's on. entries say, naming the question.
func TestCommandAskedArbiterFailFailsRun(t *testing.T) {
	e, r := autonomousSetup(t, askedWF)
	journalAsked(t, e, r, "Q9")
	r.agentSeq["S-decide-Q9"] = []AgentStatus{{SaidDone: true}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	wantDetail := "arbiter failed on question review-x-r3-f1 (Q9): no live Decision answers question Q9"
	if st.Result != Failed || st.Detail != wantDetail || !reflect.DeepEqual(outcomes(st), []string{"check:error"}) {
		t.Fatalf("result %q detail %q steps %v", st.Result, st.Detail, outcomes(st))
	}
	if r.commands["next"] != 0 {
		t.Error("next ran")
	}
}

// TestCommandAskedNormalRunLeftOpen: a normal run runs no arbiter; the
// question stays open for the owner and the run goes on.
func TestCommandAskedNormalRunLeftOpen(t *testing.T) {
	e, r, _, _ := setup(t, askedWF)
	journalAsked(t, e, r, "Q9")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.arbiters) != 0 || countType(t, e, TypeArbiter) != 0 || countType(t, e, TypeArbitrated) != 0 {
		t.Fatalf("result %q arbiters %v", st.Result, r.arbiters)
	}
}

// TestCommandAskedNotArbitratedTwice: a question already decided is not
// decided again by a later command step, nor by a resume; a session's
// question (it has a Session) is not this path's.
func TestCommandAskedNotArbitratedTwice(t *testing.T) {
	src := strings.Replace(askedWF, "[step.next]\nkind = \"command\"\nrun = \"true\"", "[step.next]\nkind = \"command\"\nrun = \"true\"\non.pass = \"again\"\n[step.again]\nkind = \"command\"\nrun = \"true\"", 1)
	e, r := autonomousSetup(t, src)
	if err := e.Store.Append(runID, Entry{Type: TypeAsked, At: at(1), Step: "build", Session: "S-build", Question: "Q1", Handle: "question-s-1"}); err != nil {
		t.Fatal(err)
	}
	journalAsked(t, e, r, "Q9")
	r.agentSeq["S-decide-Q9"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q9"] = Answer{OK: true, Decision: "assume-x"}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"check:pass", "next:pass", "again:pass"}) {
		t.Fatalf("result %q steps %v detail %q", st.Result, outcomes(st), st.Detail)
	}
	if !reflect.DeepEqual(r.arbiters, []string{"decide Q9"}) {
		t.Errorf("arbiters %v", r.arbiters)
	}
}

// TestCommandAskedResume: a run killed after the arbiter launched, before
// the step ended, reruns the command and rejoins the open arbitration; a
// run killed after the verdict does not arbitrate again.
func TestCommandAskedResume(t *testing.T) {
	e, r := autonomousSetup(t, askedWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "check", Visit: 1},
		{Type: TypeAsked, At: at(1), Step: "check", Question: "Q9", Handle: "review-x-r3-f1"},
		{Type: TypeArbiter, At: at(1), Step: "check", Session: "S-earlier", Question: "Q9"},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	r.agentSeq["S-earlier"] = []AgentStatus{{SaidDone: true}}
	r.answers["Q9"] = Answer{OK: true, Decision: "assume-x"}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.arbiters) != 0 || r.polls["S-earlier"] != 1 || countType(t, e, TypeArbitrated) != 1 {
		t.Fatalf("result %q detail %q arbiters %v polls %v", st.Result, st.Detail, r.arbiters, r.polls)
	}

	e, r = autonomousSetup(t, askedWF)
	for _, en := range []Entry{
		{Type: TypeEntered, At: at(1), Step: "check", Visit: 1},
		{Type: TypeAsked, At: at(1), Step: "check", Question: "Q9", Handle: "review-x-r3-f1"},
		{Type: TypeArbiter, At: at(1), Step: "check", Session: "S-earlier", Question: "Q9"},
		{Type: TypeArbitrated, At: at(1), Step: "check", Question: "Q9", Outcome: "pass", Decision: "assume-x"},
	} {
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	st, err = e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.arbiters) != 0 || r.polls["S-earlier"] != 0 || countType(t, e, TypeArbitrated) != 1 {
		t.Fatalf("result %q arbiters %v polls %v", st.Result, r.arbiters, r.polls)
	}
}

// TestCommandAskedRetryArbitratesAgain: after a failed arbitration fails
// the run, a retry of the step has the arbiter decide the question anew.
func TestCommandAskedRetryArbitratesAgain(t *testing.T) {
	e, r := autonomousSetup(t, askedWF)
	journalAsked(t, e, r, "Q9")
	r.agentSeq["S-decide-Q9"] = []AgentStatus{{SaidDone: true}}
	if st, err := e.Run(runID); err != nil || st.Result != Failed {
		t.Fatalf("first run: %v %q", err, st.Result)
	}
	r.onCommand = nil // the question is journalled once, as review-check does
	r.answers["Q9"] = Answer{OK: true, Decision: "assume-x"}
	st, err := e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(r.arbiters) != 2 {
		t.Fatalf("result %q detail %q arbiters %v", st.Result, st.Detail, r.arbiters)
	}
}

// TestAskedWithoutSessionIsNotTheBlocks: a session-less asked entry never
// becomes the asked entry of a block.
func TestAskedWithoutSessionIsNotTheBlocks(t *testing.T) {
	j := Journal{ID: runID, Entries: []Entry{
		startedEntry("build", 0),
		{Type: TypeEntered, At: at(1), Step: "build", Visit: 1},
		{Type: TypeBlocked, At: at(1), Step: "build", Block: &Block{Kind: "question"}},
		{Type: TypeAsked, At: at(1), Step: "check", Question: "Q9", Handle: "review-x-r3-f1"},
	}}
	if a, ok := j.State().BlockAsked(""); ok {
		t.Errorf("block asked %+v", a)
	}
}
