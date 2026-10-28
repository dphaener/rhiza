# review-signoff-verified

evidence · `01M600YD3MD5EZNAYZJV9C7SGF` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran at c0b0739 in a detached worktree with RHIZA_* unset; all pass. review-signoff.txtar: rhiza run start stops with current signoff, signoff_waiting true, detail waiting for the owner's sign-off; rhiza run approve prints approved by the actor, run succeeded, run status --json signoffs entry {actor test, outcome approved, rule owner, step signoff, visit 1}, change set stays draft. Outside a waiting review run both run approve and run request-changes exit 2 with no run, a run stopped at wait step ci (journal line count unchanged), and a finished run. review-signoff-changes.txtar: check alpha ran cycles 1, 2 and 3 (cap used); run request-changes with one finding on leaf thing starts role=fixer cycle=4 leaf=thing findings=yes, then its verifier, then only alpha re-runs; the fix commit carries Leaf: thing; the run is stopped with signoff_waiting true and next approve; a third visit is refused exit 2 over max_visits. review-signoff-autonomous.txtar: --autonomous run succeeds, signoff passes approved by rule autonomous, journal signoff entry has rule autonomous and no actor. describe-golden.txtar passes and lists run.approve, run.request-changes, signoffs and signoff_waiting. Unit tests TestSignoff* (internal/run), TestSignoffStep* (internal/workflow), TestSignoffLine (internal/api) pass; go test ./... passes.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./cmd/rhiza -run 'TestScripts/review-signoff' -v -count=1; same env: go test ./internal/run -run TestSignoff -v -count=1; go test ./internal/workflow -run TestSignoffStep -v -count=1; go test ./internal/api -run TestSignoffLine -v -count=1; go test ./cmd/rhiza -run 'TestScripts/describe-golden' -count=1; go test ./... -count=1. Commit c0b0739; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-signoff](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-signoff.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
