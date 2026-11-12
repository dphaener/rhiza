# done-refuses-unrecorded-check-verified

evidence · `01M784V43YHX3ARDHS9JNDV0QB` · actor [claude-code-09c4fb40](../actor/claude-code-09c4fb40.md)

## result

PASS

## statement

At e5443c542, go test -trimpath ./internal/api -run SessionDone -count=1 passes (TestSessionDoneReviewRecords: check session refused naming `rhiza run findings` until a round is journaled for its check and cycle; critic refused naming `rhiza run verdict` while a blocker/major is unjudged, also with only the major left, and succeeds once both have verdicts; fixer never refused). The plan reviewer case (TestReviewVisitNeedsRound) and TestReviewCheckFindingsFlag pass too. The committed test journals a non-empty round, and its refusal assertion checks the message rather than the error code, so I ran a scratch probe test in a copy outside the repo. In it the refusal is CodeConflict, a check session succeeds after an empty {"findings":[]} round, and a critic with one unjudged blocker is refused, then succeeds after its verdict. Read the diff: session_done.go calls reviewRecordMissing only for role check/critic sessions with a run and a check, and skips when the run is ended or the step moved on.

## Links

### Out

- verifies [done-refuses-unrecorded-check](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-check.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
