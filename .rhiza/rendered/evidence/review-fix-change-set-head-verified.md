# review-fix-change-set-head-verified

evidence · `01M6076M5G9CVYZTF2CZNPTZP6` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

At d60c3c3, after fix 3e6a485 (change-set fixer prompt now reads 'Confirmed findings to fix in review cycle N, which name no leaf of the change set:'), TestFixChangeSetEngine passes: two confirmed findings naming no leaf go to one fixer session whose prompt carries both and the Review: R1 trailer, the new commit's Review: trailer names the run, and the failed checks re-run after the fixer is judged; TestFixChangeSetEngineNoCommit and TestFixChangeSetNoCommitFailsCycle show a fixer without such a commit fails the cycle. TestReviewCommit, TestReviewFixFold, TestOwnerCycleFixerPrompts (including the new non-owner-cycle prompt assertion) pass; go test ./... passes.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/run -run 'TestFixChangeSet|TestReviewFixFold' -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/api -run 'TestReviewCommit|TestFixChangeSet|TestOwnerCycleFixerPrompts' -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./... -count=1; commit d60c3c3 (detached worktree); run by claude-code-0e5a9d92-verifier on 2026-10-28

## Links

### Out

- verifies [review-fix-change-set](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-change-set.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
