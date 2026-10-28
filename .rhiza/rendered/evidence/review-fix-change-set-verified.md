# review-fix-change-set-verified

evidence · `01M5ZXXVWBC99GDC6MQ2251Z4K` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran in a detached worktree at 5b498bc with RHIZA_* unset; all passed. Two confirmed no-leaf findings to ONE fixer: internal/run TestFixChangeSetOneFixer (alpha f1 major + gamma f1 blocker -> len(fixSpecs)==1 with both findings) and internal/api TestFixChangeSetEngine (one fixer session, prompt names both findings). New commit with Review: naming the run: TestFixChangeSetEngine (fixer commits Review: R1, judged pass on that commit) and TestReviewCommit (none / other run / no Review trailer / found / only since start). Failed checks re-run after it: TestFixChangeSetOneFixer (review-fix-judged precedes cycle-2 check sessions; only alpha and gamma re-run) and TestFixChangeSetEngine (judged entry before cycle 2 check). Fixer without such a commit fails the cycle: TestFixChangeSetNoCommitFailsCycle (alpha cycle 2 'fix failed: change set: fixer failed: no new commit', step fails) and TestFixChangeSetEngineNoCommit (commit naming another run -> fail, no re-check). internal/check change only adds a Commit.Reviews field, read by api ReviewCommit, not by any check rule; go test ./internal/check and the check/review/build scripts pass, so rhiza check is unchanged. Full go test ./... passed.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/run -run 'TestFixChangeSet|TestReviewFixFold' -v -count=1; same env: go test ./internal/api -run 'TestReviewCommit|TestFixChangeSet' -v -count=1; go test ./internal/check -count=1; go test ./cmd/rhiza -run 'TestScripts/(review-|build-|check)' -count=1; go test ./... -count=1. Commit 5b498bc; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-fix-change-set](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-change-set.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
