package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A sign-off's line names the owner by actor, an autonomous run's by its
// rule, and says changes were asked for; a run waits at its sign-off only
// while stopped there, unfinished.
func TestSignoffLine(t *testing.T) {
	for _, tt := range []struct {
		s    RunSignoff
		want string
	}{
		{RunSignoff{Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved}, "signoff (visit 1): approved by ada"},
		{RunSignoff{Step: "signoff", Visit: 2, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffChanges}, "signoff (visit 2): changes asked by ada"},
		{RunSignoff{Step: "signoff", Visit: 1, Rule: run.RuleAutonomous, Outcome: run.SignoffApproved}, "signoff (visit 1): approved by rule autonomous"},
	} {
		if got := tt.s.line(); got != tt.want {
			t.Errorf("line = %q, want %q", got, tt.want)
		}
	}
	if !signoffWaiting(run.State{Stopped: run.SignoffWait}) {
		t.Error("a run stopped at its sign-off is not waiting")
	}
	if signoffWaiting(run.State{Stopped: run.SignoffWait, Ended: true}) || signoffWaiting(run.State{Stopped: "waiting: false"}) {
		t.Error("a finished run, or one stopped elsewhere, is waiting")
	}
}

// In an owner cycle the fixers' prompts say the owner asked for the
// changes, not that a review cycle confirmed them.
func TestOwnerCycleFixerPrompts(t *testing.T) {
	role := workflow.Role{Prompt: "fix {leaf}"}
	vars := workflow.PromptVars{Run: "R1"}
	fs := []run.Finding{{ID: "f1", Check: "alpha", Severity: "major", Place: "a.go:1", Reason: "say more"}}
	leaf := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "thing"}, Cycle: 4, Owner: true, Findings: fs}, vars)
	cs := changeSetFixerPrompt(role, run.FixSpec{Cycle: 4, Owner: true, Findings: fs}, vars)
	for _, tt := range []struct{ p, want string }{
		{leaf, "The owner requested these changes on leaf thing:\n\n- alpha f1 (major at a.go:1): say more\n"},
		{cs, "The owner requested these changes, which name no leaf of the change set:\n\n- alpha f1 (major at a.go:1): say more\n"},
	} {
		if !strings.Contains(tt.p, tt.want) || strings.Contains(tt.p, "review cycle") {
			t.Errorf("prompt lacks %q or names a review cycle:\n%s", tt.want, tt.p)
		}
	}
	plain := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "thing"}, Cycle: 2, Findings: fs}, vars)
	if !strings.Contains(plain, "Confirmed findings to fix in review cycle 2 on leaf thing:") {
		t.Errorf("a fixer outside an owner cycle is not told its review cycle:\n%s", plain)
	}
	plainCS := changeSetFixerPrompt(role, run.FixSpec{Cycle: 2, Findings: fs}, vars)
	if !strings.Contains(plainCS, "Confirmed findings to fix in review cycle 2, which name no leaf of the change set:") {
		t.Errorf("the change set's fixer outside an owner cycle is not told its review cycle:\n%s", plainCS)
	}
}

// signoffRaceWF is a review: a checks step, then a signoff step whose
// changes go back to it.
const signoffRaceWF = `[step.review]
kind = "checks"
interval = "1s"
max_visits = 3
on.pass = "signoff"
[[step.review.check]]
name = "alpha"
prompt = "check {check}"
[step.review.critic]
prompt = "criticise {check}"
[step.review.fixer]
prompt = "fix {leaf}"
[step.review.verifier]
prompt = "verify {leaf}"
[step.signoff]
kind = "signoff"
max_visits = 3
on.pass = "done"
on.changes = "review"
`

const signoffRaceRun = "01M5WPTZ1V0JVVW8FPX0QYMRWA"

// signoffRaceFixture is a review run stopped at its sign-off, written
// straight to the journal, and its store.
func signoffRaceFixture(t *testing.T) (*sessionFixture, *run.Store) {
	t.Helper()
	f := newSessionFixture(t, "race")
	cs := newChange(t, f.ctx, `{"intent":"X","title":"Race"}`)
	f.ctx.Start = (&shStart{}).start
	def, err := workflow.Parse("review", []byte(signoffRaceWF))
	if err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	at := time.Unix(100, 0)
	if err := run.Start(rs, signoffRaceRun, run.Started{
		Workflow: "review", Definition: signoffRaceWF, Change: cs.ChangeSet.String(), Worktree: f.worktree, PermissionMode: "manual",
	}, at, def); err != nil {
		t.Fatal(err)
	}
	for _, e := range []run.Entry{
		{Type: run.TypeEntered, Step: "review", Visit: 1},
		{Type: run.TypeEnded, Step: "review", Visit: 1, Outcome: "pass", Next: "signoff"},
		{Type: run.TypeEntered, Step: "signoff", Visit: 1},
		{Type: run.TypeStopped, Step: "signoff", Detail: run.SignoffWait},
	} {
		e.At = at
		if err := rs.Append(signoffRaceRun, e); err != nil {
			t.Fatal(err)
		}
	}
	return f, rs
}

// raceSignoffs is the run's signoff entries.
func raceSignoffs(t *testing.T, rs *run.Store) []run.Entry {
	t.Helper()
	j, err := rs.Read(signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	var out []run.Entry
	for _, e := range j.Entries {
		if e.Type == run.TypeSignoff {
			out = append(out, e)
		}
	}
	return out
}

// A sign-off that passed its first checks before another process took
// the run's lock exits 3 (running in another process) and journals
// nothing: the sign-off is written only under the run's lock.
func TestSignoffUnderLockHeld(t *testing.T) {
	f, rs := signoffRaceFixture(t)
	sr, err := openSignoffRun(f.ctx, "run approve", signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	unlock, ok, err := rs.TryLock(signoffRaceRun)
	if err != nil || !ok {
		t.Fatalf("lock: %v %v", ok, err)
	}
	defer unlock()
	_, err = sr.signoff(f.ctx, run.Entry{
		Type: run.TypeSignoff, At: f.ctx.Now(), Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved,
	}, "")
	wantExit(t, err, 3)
	if !strings.Contains(AsError(err).Message, "running in another process") {
		t.Errorf("message %q", AsError(err).Message)
	}
	if got := raceSignoffs(t, rs); len(got) != 0 {
		t.Errorf("journalled %+v under another's lock", got)
	}
}

// Two sign-offs that both pass the first checks: the second, rechecking
// under the run's lock, sees the run the first finished, exits 2 and
// journals nothing.
func TestSignoffRecheckedUnderLock(t *testing.T) {
	f, rs := signoffRaceFixture(t)
	first, err := openSignoffRun(f.ctx, "run approve", signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	second, err := openSignoffRun(f.ctx, "run request-changes", signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	// The first wins the lock and journals its approval; its run then
	// goes on to its end.
	if _, err := first.signoff(f.ctx, run.Entry{
		Type: run.TypeSignoff, At: f.ctx.Now(), Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved,
	}, ""); err != nil {
		t.Fatal(err)
	}
	_, err = second.signoff(f.ctx, run.Entry{
		Type: run.TypeSignoff, At: f.ctx.Now(), Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffChanges,
		Findings: []run.Finding{{ID: "f1", Check: "alpha", Severity: "major", Place: "a.go:1", Reason: "no"}},
	}, "review")
	wantExit(t, err, 2)
	if !strings.Contains(AsError(err).Message, "is finished") {
		t.Errorf("message %q", AsError(err).Message)
	}
	got := raceSignoffs(t, rs)
	if len(got) != 1 || got[0].Outcome != run.SignoffApproved {
		t.Errorf("signoffs %+v, want the one approval", got)
	}
}

// A sign-off that passed its first checks, then finds under the lock that
// another process journalled a sign-off of the visit (and stopped before
// running it), exits 2 (already signed off) and journals nothing.
func TestSignoffAlreadySignedUnderLock(t *testing.T) {
	f, rs := signoffRaceFixture(t)
	sr, err := openSignoffRun(f.ctx, "run approve", signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	if err := rs.Append(signoffRaceRun, run.Entry{
		Type: run.TypeSignoff, At: f.ctx.Now(), Step: "signoff", Visit: 1, Actor: "bob", Rule: run.RuleOwner, Outcome: run.SignoffApproved,
	}); err != nil {
		t.Fatal(err)
	}
	_, err = sr.signoff(f.ctx, run.Entry{
		Type: run.TypeSignoff, At: f.ctx.Now(), Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved,
	}, "")
	wantExit(t, err, 2)
	if !strings.Contains(AsError(err).Message, "already signed off") {
		t.Errorf("message %q", AsError(err).Message)
	}
	if got := raceSignoffs(t, rs); len(got) != 1 || got[0].Actor != "bob" {
		t.Errorf("signoffs %+v, want bob's alone", got)
	}
}

// A sign-off still waiting when the lock is taken is journalled once.
func TestSignoffUnderLockJournalsOnce(t *testing.T) {
	f, rs := signoffRaceFixture(t)
	resp, err := callOp(t, f.ctx, "run.approve", `{"run":"`+signoffRaceRun+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RunSignoffResponse); r.Run.State != run.Succeeded {
		t.Errorf("run %+v", r.Run)
	}
	if got := raceSignoffs(t, rs); len(got) != 1 {
		t.Errorf("signoffs %+v", got)
	}
}
