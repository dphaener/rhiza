package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// The summary renders the parts no testscript reaches: a failed code
// check, a serious finding with no verdict, a change-set fixer, a
// question raised at the cap, an owner cycle with the owner's findings,
// a request for changes and an approval. It reads the journal only.
func TestReviewSummary(t *testing.T) {
	sha := "0123456789abcdef0123456789abcdef01234567"
	j := run.Journal{ID: "R1", Entries: []run.Entry{
		{Type: run.TypeStarted, Started: &run.Started{Workflow: "review", Change: "C1", Title: "Tidy up"}},
		{Type: run.TypeCodeCheck, Step: "code", Outcome: run.CodeCheckFail, Detail: "leaf thing: no Evidence"},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 1, Checks: []string{"alpha"}},
		{Type: run.TypeFindings, Step: "review", Check: "alpha", Cycle: 1, Findings: []run.Finding{
			{ID: "f1", Severity: run.SeverityBlocker, Place: "main.go:1", Reason: "crashes"},
		}},
		{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 1, Outcome: run.CheckFailed, Detail: "cycle 1: confirmed f1 (blocker)"},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 2, Checks: []string{"alpha"}},
		{Type: run.TypeReviewFix, Step: "review", Role: workflow.RoleFixer, Cycle: 2, Session: "S2"},
		{Type: run.TypeReviewFixJudged, Step: "review", Role: workflow.RoleFixer, Cycle: 2, Outcome: workflow.OutcomeFail, Detail: "no new commit since abc"},
		{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 2, Outcome: run.CheckFailed, Detail: "cycle 2: fix failed: change set: fixer failed: no new commit since abc"},
		{Type: run.TypeReviewQuestion, Step: "review", Check: "alpha", Cycle: 2, Finding: "f1", Question: "Q1", QuestionHandle: "review-failed-r1-1"},
		{Type: run.TypeSignoff, Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffChanges, Findings: []run.Finding{
			{ID: "f1", Check: "alpha", Severity: run.SeverityMinor, Leaf: "L1", LeafHandle: "thing", Place: "thing.txt:1", Reason: "say more"},
		}},
		{Type: run.TypeCheckCycle, Step: "review", Cycle: 3, Checks: []string{"alpha"}, Kind: run.CycleOwner, Findings: []run.Finding{
			{ID: "f1", Check: "alpha", Severity: run.SeverityMinor, Leaf: "L1", LeafHandle: "thing", Place: "thing.txt:1", Reason: "say more"},
		}},
		{Type: run.TypeLeafSession, Step: "review", Leaf: "L1", Handle: "thing", Role: workflow.RoleFixer, Attempt: 5, Cycle: 3, Session: "S3"},
		{Type: run.TypeLeafAttempt, Step: "review", Leaf: "L1", Role: workflow.RoleFixer, Attempt: 5, Outcome: workflow.OutcomeFail, Detail: "no commit naming its leaf"},
		{Type: run.TypeLeafSession, Step: "review", Leaf: "L1", Handle: "thing", Role: workflow.RoleFixer, Attempt: 6, Cycle: 3, Session: "S4"},
		{Type: run.TypeLeafAttempt, Step: "review", Leaf: "L1", Role: workflow.RoleFixer, Attempt: 6, Outcome: workflow.OutcomePass, Commit: sha},
		{Type: run.TypeLeafSession, Step: "review", Leaf: "L1", Role: workflow.RoleVerifier, Attempt: 6, Cycle: 3, Session: "S5"},
		{Type: run.TypeFindings, Step: "review", Check: "alpha", Cycle: 3, Findings: []run.Finding{}},
		{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 3, Outcome: run.CheckPassed},
		{Type: run.TypeSignoff, Step: "signoff", Visit: 2, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved},
		{Type: run.TypeFinished, End: run.Succeeded},
	}}
	want := `# Review of run R1

Change set C1: Tidy up. Workflow review; the run succeeded.

## Code check

- step code: fail: leaf thing: no Evidence

## Checks of step review

### alpha: passed

- cycle 1: failed: confirmed f1 (blocker)
- cycle 2: failed: fix failed: change set: fixer failed: no new commit since abc
- cycle 3: passed

Confirmed:

- cycle 1: f1 blocker at main.go:1: crashes — no verdict, so it stands

Fix cycles:

- cycle 2: re-runs alpha
  - change set: fixer fail (no new commit since abc)
  - question review-failed-r1-1 on check alpha f1
- cycle 3, the owner's changes: re-runs alpha
  - owner's alpha f1 minor at thing.txt:1 on thing: say more
  - leaf thing: fixer fail (no commit naming its leaf)
  - leaf thing, try 2: fixer pass, commit 0123456789ab; verifier open

## Sign-off

- signoff (visit 1): changes asked by ada
  - alpha f1 minor at thing.txt:1 on thing: say more
- signoff (visit 2): approved by ada
`
	if got := reviewSummary(j.State()); got != want {
		t.Errorf("summary:\n%s\nwant:\n%s", got, want)
	}

	waiting := run.Journal{ID: "R2", Entries: []run.Entry{
		{Type: run.TypeStarted, Started: &run.Started{Workflow: "review", Change: "C1"}},
		{Type: run.TypeEntered, Step: "signoff", Visit: 1},
		{Type: run.TypeStopped, Step: "signoff", Detail: run.SignoffWait},
	}}
	want = `# Review of run R2

Change set C1. Workflow review; the run waits at step signoff for the owner's sign-off.

## Code check

None ran.

## Checks

None ran.

## Sign-off

- signoff (visit 1): waiting for the owner: rhiza run approve R2, or rhiza run request-changes <file> R2
`
	if got := reviewSummary(waiting.State()); got != want {
		t.Errorf("summary:\n%s\nwant:\n%s", got, want)
	}
}
