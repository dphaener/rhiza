# review-workflow-file-head-verified

evidence · `01M6076MGVF994ZCZDAN84TWZV` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

At d60c3c3, after fix 7a28dd0 (builtin review.toml fixer prompt drops the leaf clause; TestReviewBuiltin asserts no double space for leaf and change-set fixers), TestScripts/review-workflow passes: one check's major is confirmed by the critic, a fixer and the leaf's verifier run, only that check re-runs clean, the qa check's Evidence commit carries qa's own actor, build-check exits 0, run approve ends the run succeeded and the change set stays accepted (not landed). TestReviewBuiltin and TestScripts/describe-golden pass (golden current); go test ./... passes.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./cmd/rhiza -run 'TestScripts/review-workflow' -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/workflow -run TestReviewBuiltin -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./cmd/rhiza -run TestScripts/describe-golden -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./... -count=1; commit d60c3c3 (detached worktree); run by claude-code-0e5a9d92-verifier on 2026-10-28

## Links

### Out

- verifies [review-workflow-file](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-workflow-file.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
