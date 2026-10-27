package run

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func (f *fakeRunner) StartCheck(spec CheckSpec) (string, error) {
	if f.checkErr != nil {
		return "", f.checkErr
	}
	f.checkSpecs = append(f.checkSpecs, spec)
	return checkSID(spec.Check.Name, spec.Role, spec.Cycle), nil
}

// checkSID is the fake runner's session ID for a check session.
func checkSID(check, role string, cycle int) string {
	return "S-" + check + "-" + role + "-" + strconv.Itoa(cycle)
}

const checksWF = `[step.review]
kind = "checks"
interval = "5s"
[[step.review.check]]
name = "alpha"
prompt = "check {check} of {change}"
models = ["sonnet", "opus"]
[[step.review.check]]
name = "beta"
prompt = "check {check}"
[[step.review.check]]
name = "gamma"
prompt = "check {check}"
[step.review.critic]
prompt = "criticise {check}"
[step.review.fixer]
prompt = "fix {leaf}"
[step.review.verifier]
prompt = "verify {leaf}"
`

var threeChecks = []string{"alpha", "beta", "gamma"}

// passChecks makes every check session of cycle exit 0.
func passChecks(r *fakeRunner, cycle int, checks ...string) {
	for _, c := range checks {
		r.outcome[checkSID(c, workflow.RoleCheck, cycle)] = workflow.OutcomePass
	}
}

// checkEntries are the journal's checks step entries as "type check role
// cycle outcome", a check-cycle's as "check-cycle cycle checks...".
func checkEntries(t *testing.T, e *Engine) []string {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, en := range j.Entries {
		switch en.Type {
		case TypeCheckCycle:
			out = append(out, strings.Join(append([]string{en.Type, attemptStr(en.Cycle)}, en.Checks...), " "))
		case TypeCheckSession, TypeCheckJudged, TypeCheck:
			out = append(out, strings.Join(strings.Fields(strings.Join([]string{en.Type, en.Check, en.Role, attemptStr(en.Cycle), en.Outcome}, " ")), " "))
		}
	}
	return out
}

// A checks step over three checks starts every check's session before
// it polls any, judges each, and passes when all three pass in cycle 1.
func TestChecksStepParallel(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	polled := false
	r.onPoll = func(session string, n int) {
		if !polled && len(r.checkSpecs) != 3 {
			t.Errorf("first poll (%s) after %d of 3 check sessions started", session, len(r.checkSpecs))
		}
		polled = true
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	want := []string{
		"check-cycle 1 alpha beta gamma",
		"check-session alpha check 1", "check-session beta check 1", "check-session gamma check 1",
		"check-judged alpha check 1 pass", "check-judged beta check 1 pass", "check-judged gamma check 1 pass",
		"check alpha 1 passed", "check beta 1 passed", "check gamma 1 passed",
	}
	if got := checkEntries(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("check entries =\n%v\nwant\n%v", got, want)
	}
	if d := st.Steps[0].Detail; d != "3 of 3 checks passed" {
		t.Errorf("detail = %q", d)
	}
	if len(st.Checks) != 3 || len(st.Cycles) != 1 {
		t.Fatalf("Checks = %+v, Cycles = %+v", st.Checks, st.Cycles)
	}
	for _, c := range st.Checks {
		if o, ok := c.Latest(); !ok || o.Outcome != CheckPassed || c.Step != "review" {
			t.Errorf("check %s = %+v", c.Check, c)
		}
	}
	// Each session got its check, cycle and role, and the check's model
	// for the cycle.
	sp := r.checkSpecs[0]
	if sp.Check.Name != "alpha" || sp.Role != workflow.RoleCheck || sp.Cycle != 1 || sp.Model != "sonnet" {
		t.Errorf("alpha's spec = %+v", sp)
	}
	for _, v := range []string{"RHIZA_CHANGE=CS1", "RHIZA_STEP=review", EnvCheck + "=alpha", EnvCycle + "=1", EnvRole + "=check"} {
		if !slices.Contains(sp.Env, v) {
			t.Errorf("alpha's env %v lacks %s", sp.Env, v)
		}
	}
	// Every session's end was reported.
	if len(r.ended) != 3 {
		t.Errorf("AgentEnded calls = %v", r.ended)
	}
}

// A check that fails in cycle 1 is the only one started in cycle 2, on
// its ladder's next model; once it passes the step passes.
func TestChecksStepReRunsFailed(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, "beta", "gamma")
	r.outcome[checkSID("alpha", workflow.RoleCheck, 1)] = workflow.OutcomeFail
	passChecks(r, 2, "alpha")
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	var started []string
	for _, s := range r.checkSpecs {
		started = append(started, s.Check.Name+" "+strconv.Itoa(s.Cycle)+" "+s.Model)
	}
	if want := []string{"alpha 1 sonnet", "beta 1 ", "gamma 1 ", "alpha 2 opus"}; !reflect.DeepEqual(started, want) {
		t.Errorf("started %q, want %q", started, want)
	}
	got := checkEntries(t, e)
	if want := []string{"check-cycle 2 alpha", "check-session alpha check 2", "check-judged alpha check 2 pass", "check alpha 2 passed"}; !reflect.DeepEqual(got[len(got)-4:], want) {
		t.Errorf("cycle 2 entries = %v, want %v", got[len(got)-4:], want)
	}
	a, _ := st.Check("review", "alpha")
	if len(a.Outcomes) != 2 || a.Outcomes[0].Outcome != CheckFailed || a.Outcomes[1].Outcome != CheckPassed {
		t.Errorf("alpha's outcomes = %+v", a.Outcomes)
	}
	if c, ok := st.LatestCycle("review"); !ok || c.Cycle != 2 || !reflect.DeepEqual(c.Checks, []string{"alpha"}) {
		t.Errorf("latest cycle = %+v", c)
	}
}

// A check that fails every cycle ends the step fail after 1 +
// MaxFixCycles cycles, naming it.
func TestChecksStepCycleLimit(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, "alpha", "gamma")
	for n := 1; n <= 1+MaxFixCycles; n++ {
		r.outcome[checkSID("beta", workflow.RoleCheck, n)] = workflow.OutcomeFail
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:fail"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if want := "2 of 3 checks passed after 3 cycles; failed: beta"; st.Steps[0].Detail != want {
		t.Errorf("detail = %q, want %q", st.Steps[0].Detail, want)
	}
	if len(r.checkSpecs) != 3+MaxFixCycles {
		t.Errorf("%d check sessions started, want %d", len(r.checkSpecs), 3+MaxFixCycles)
	}
	if len(st.Cycles) != 1+MaxFixCycles {
		t.Errorf("cycles = %+v", st.Cycles)
	}
}

// A check session is judged pass once its agent says done; idle alone
// never judges it, and a live block is journalled.
func TestChecksStepSaidDoneNotIdle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, "beta", "gamma")
	a := checkSID("alpha", workflow.RoleCheck, 1)
	r.agentSeq[a] = []AgentStatus{{Idle: true}, {Idle: true}, {Blocked: &Block{Kind: "permission"}}, {SaidDone: true}}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if r.polls[a] != 4 {
		t.Errorf("alpha polled %d times, want 4", r.polls[a])
	}
	c, _ := st.Check("review", "alpha")
	if j, ok := c.Judgement(workflow.RoleCheck, 1); !ok || j.Outcome != workflow.OutcomePass || j.Detail != "said done" {
		t.Errorf("alpha's judgement = %+v", j)
	}
	if !slices.Contains(entryTypes(t, e), TypeBlocked+":review") {
		t.Errorf("no blocked entry")
	}
}

// A resume rejoins the check sessions the journal holds, never starting
// them again, and starts only the cycle's checks not yet started.
func TestChecksStepResumeRejoins(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, threeChecks...)
	for _, en := range []Entry{
		{Type: TypeEntered, Step: "review", Visit: 1},
		{Type: TypeCheckCycle, Step: "review", Cycle: 1, Checks: threeChecks},
		{Type: TypeCheckSession, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Session: checkSID("alpha", "check", 1)},
		{Type: TypeCheckJudged, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Outcome: "pass"},
		{Type: TypeCheckSession, Step: "review", Check: "beta", Role: "check", Cycle: 1, Session: checkSID("beta", "check", 1)},
	} {
		en.At = at(1)
		if err := e.Store.Append(runID, en); err != nil {
			t.Fatal(err)
		}
	}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	if len(r.checkSpecs) != 1 || r.checkSpecs[0].Check.Name != "gamma" {
		t.Errorf("started %+v, want gamma only", r.checkSpecs)
	}
	if r.polls[checkSID("beta", "check", 1)] != 1 || r.polls[checkSID("alpha", "check", 1)] != 0 {
		t.Errorf("polls = %v", r.polls)
	}
}

// Check sessions open in panes of the change set's tab: the first opens
// it, each later one splits the latest open check pane; once the step
// ends every judged session's pane is closed, journalled.
func TestChecksStepPanes(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	passChecks(r, 1, "alpha", "gamma")
	r.outcome[checkSID("beta", "check", 1)] = workflow.OutcomeFail
	passChecks(r, 2, "beta")
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	var anchors []string
	for _, s := range r.checkSpecs {
		anchors = append(anchors, s.Anchor)
	}
	a1, b1, g1, b2 := checkSID("alpha", "check", 1), checkSID("beta", "check", 1), checkSID("gamma", "check", 1), checkSID("beta", "check", 2)
	if want := []string{"", a1, b1, g1}; !reflect.DeepEqual(anchors, want) {
		t.Errorf("anchors = %q, want %q", anchors, want)
	}
	if want := []string{a1, b1, g1, b2}; !reflect.DeepEqual(r.closed, want) {
		t.Errorf("closed = %v, want %v", r.closed, want)
	}
	st, _ := e.Store.Read(runID)
	for _, c := range st.State().Checks {
		for _, s := range c.Sessions {
			if !s.PaneClosed {
				t.Errorf("session %s's pane is open in the journal", s.Session)
			}
		}
	}
	if len(st.State().Leaves) != 0 {
		t.Errorf("a check's pane-closed entry made a leaf: %+v", st.State().Leaves)
	}
}

// checkPane anchors a check session on the latest open check pane, or
// opens the tab when none is open.
func TestCheckPane(t *testing.T) {
	states := map[string]CheckState{
		"alpha": {Step: "review", Check: "alpha", Sessions: []CheckSession{{Role: "check", Cycle: 1, Session: "SA", Seq: 3}}},
		"beta":  {Step: "review", Check: "beta", Sessions: []CheckSession{{Role: "check", Cycle: 1, Session: "SB", Seq: 5}}},
		"gamma": {Step: "review", Check: "gamma", Sessions: []CheckSession{{Role: "check", Cycle: 1, Session: "SG", Seq: 7, PaneClosed: true}}},
	}
	if got := checkPane(states); got != "SB" {
		t.Errorf("checkPane = %q, want SB", got)
	}
	if got := checkPane(map[string]CheckState{"gamma": states["gamma"]}); got != "" {
		t.Errorf("checkPane with every pane closed = %q, want empty", got)
	}
}

// nextCheckActions begins cycle 1 with every check; then starts every
// unstarted check before any poll; ends a judged check; begins the next
// cycle with only the failed checks; and ends the step once nothing more
// can run.
func TestNextCheckActions(t *testing.T) {
	step := workflow.Step{Name: "review", Checks: []workflow.Check{{Name: "a"}, {Name: "b"}}}
	kinds := func(acts []checkAction) []string {
		var out []string
		for _, a := range acts {
			out = append(out, strconv.Itoa(int(a.kind))+":"+a.check+strings.Join(a.checks, ",")+":"+strconv.Itoa(a.cycle)+":"+a.outcome)
		}
		return out
	}
	sess := func(cycle int, judged string) CheckState {
		c := CheckState{Sessions: []CheckSession{{Role: "check", Cycle: cycle, Session: "S"}}}
		if judged != "" {
			c.Judged = []CheckJudged{{Role: "check", Cycle: cycle, Outcome: judged}}
		}
		return c
	}
	ended := func(c CheckState, cycle int, o string) CheckState {
		c.Outcomes = append(c.Outcomes, CheckOutcome{Cycle: cycle, Outcome: o})
		return c
	}
	cyc := func(n int, checks ...string) CheckCycle { return CheckCycle{Step: "review", Cycle: n, Checks: checks} }
	for _, tt := range []struct {
		name   string
		cycles []CheckCycle
		states map[string]CheckState
		want   []string
	}{
		{"begin", nil, nil, []string{"1:a,b:1:"}},
		{"start all before polling", []CheckCycle{cyc(1, "a", "b")}, map[string]CheckState{"a": sess(1, "")},
			[]string{"2:b:1:", "3:a:1:"}},
		{"end judged", []CheckCycle{cyc(1, "a", "b")}, map[string]CheckState{"a": sess(1, "fail"), "b": sess(1, "")},
			[]string{"4:a:1:failed", "3:b:1:"}},
		{"next cycle with failed only", []CheckCycle{cyc(1, "a", "b")},
			map[string]CheckState{"a": ended(sess(1, "fail"), 1, CheckFailed), "b": ended(sess(1, "pass"), 1, CheckPassed)},
			[]string{"1:a:2:"}},
		{"all passed", []CheckCycle{cyc(1, "a", "b")},
			map[string]CheckState{"a": ended(sess(1, "pass"), 1, CheckPassed), "b": ended(sess(1, "pass"), 1, CheckPassed)},
			[]string{"5::0:pass"}},
		{"out of cycles", []CheckCycle{cyc(1, "a", "b"), cyc(2, "a"), cyc(3, "a")},
			map[string]CheckState{"a": ended(ended(ended(sess(3, "fail"), 1, CheckFailed), 2, CheckFailed), 3, CheckFailed), "b": ended(sess(1, "pass"), 1, CheckPassed)},
			[]string{"5::0:fail"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := kinds(nextCheckActions(step, tt.cycles, tt.states))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("actions = %v, want %v", got, tt.want)
			}
		})
	}
}

// A runner that cannot run checks, or a check session that cannot start,
// ends the step with an error outcome.
func TestChecksStepRunnerErrors(t *testing.T) {
	e, r, _, _ := setup(t, checksWF)
	e.Runner = struct{ Runner }{r}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:error"}) || !strings.Contains(st.Steps[0].Detail, "cannot run checks") {
		t.Errorf("no CheckRunner: %v %q", got, st.Steps[0].Detail)
	}

	e, r, _, _ = setup(t, checksWF)
	r.checkErr = errors.New("no tab")
	st, err = e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:error"}) || !strings.Contains(st.Steps[0].Detail, "check alpha: start its session: no tab") {
		t.Errorf("start fails: %v %q", got, st.Steps[0].Detail)
	}
}

// The checks state survives the step's re-entry: a new visit keeps its
// cycles and outcomes.
func TestChecksStateSurvivesEntered(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("review", 0),
		{Type: TypeEntered, Step: "review", Visit: 1},
		{Type: TypeCheckCycle, Step: "review", Cycle: 1, Checks: []string{"alpha"}},
		{Type: TypeCheckSession, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Session: "S1"},
		{Type: TypeCheckJudged, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Outcome: "pass", Detail: "exit 0"},
		{Type: TypeCheck, Step: "review", Check: "alpha", Cycle: 1, Outcome: CheckPassed},
		{Type: TypePaneClosed, Step: "review", Check: "alpha", Session: "S1"},
		{Type: TypeEnded, Step: "review", Outcome: "pass", Next: "review"},
		{Type: TypeEntered, Step: "review", Visit: 2},
	}}
	st := j.State()
	c, ok := st.Check("review", "alpha")
	if !ok || len(c.Sessions) != 1 || !c.Sessions[0].PaneClosed || len(c.Outcomes) != 1 {
		t.Fatalf("alpha = %+v, %v", c, ok)
	}
	if len(st.Cycles) != 1 || len(st.Leaves) != 0 {
		t.Errorf("Cycles = %+v, Leaves = %+v", st.Cycles, st.Leaves)
	}
	if len(st.Sessions) != 1 || st.Sessions[0].Session != "S1" {
		t.Errorf("Sessions = %+v", st.Sessions)
	}
}

// While a checks step runs, the status report counts its checks passed
// and lists its live check sessions, with their cycle.
func TestDeriveStatusChecks(t *testing.T) {
	def, err := workflow.Parse("wf", []byte(checksWF))
	if err != nil {
		t.Fatal(err)
	}
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("review", 0),
		{Type: TypeEntered, Step: "review", Visit: 1},
		{Type: TypeCheckCycle, Step: "review", Cycle: 1, Checks: threeChecks},
		{Type: TypeCheckSession, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Session: "SA"},
		{Type: TypeCheckSession, Step: "review", Check: "beta", Role: "check", Cycle: 1, Session: "SB"},
		{Type: TypeCheckSession, Step: "review", Check: "gamma", Role: "check", Cycle: 1, Session: "SG"},
		{Type: TypeCheckJudged, Step: "review", Check: "alpha", Role: "check", Cycle: 1, Outcome: "pass"},
		{Type: TypeCheck, Step: "review", Check: "alpha", Cycle: 1, Outcome: CheckPassed},
		{Type: TypeCheckJudged, Step: "review", Check: "beta", Role: "check", Cycle: 1, Outcome: "fail"},
		{Type: TypeCheck, Step: "review", Check: "beta", Cycle: 1, Outcome: CheckFailed},
		{Type: TypeCheckCycle, Step: "review", Cycle: 2, Checks: []string{"beta"}},
		{Type: TypeCheckSession, Step: "review", Check: "beta", Role: "check", Cycle: 2, Session: "SB2"},
	}}
	rep := deriveStatus(j.State(), def, nil)
	if rep.Checks != 3 || rep.Passed != 1 || rep.Progress == nil || *rep.Progress != 33 {
		t.Errorf("report = %+v", rep)
	}
	want := []LeafStatus{
		{N: 2, Handle: "beta", Role: "check", Attempt: 2, Session: "SB2"},
		{N: 3, Handle: "gamma", Role: "check", Attempt: 1, Session: "SG"},
	}
	if !reflect.DeepEqual(rep.Live, want) {
		t.Errorf("Live = %+v\nwant %+v", rep.Live, want)
	}
}
