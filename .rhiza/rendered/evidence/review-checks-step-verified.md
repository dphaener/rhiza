# review-checks-step-verified

evidence · `01M5ZRXSESD1VTJ3499DEJ0TF9` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran the leaf's tests in a detached worktree at 87ab105 with RHIZA_* unset; all passed. internal/workflow TestChecksStepLoads loads a checks step; TestChecksStepRefused rejects no checks, a duplicate check name and an unknown [[step.X.check]] key (plus critic/fixer keys, missing roles, misplaced keys) with exact messages. internal/run TestChecksStepParallel asserts the first poll comes only after all 3 check sessions started (so all 3 start before any ends) and the journal order check-session x3 before check-judged; TestChecksStepReRunsFailed fails alpha in cycle 1 and asserts the starts are alpha/beta/gamma cycle 1 then only alpha cycle 2. Testscript review-checks-step asserts rhiza run status lists each check's cycles and outcome (beta: cycle 1 failed, cycle 2 passed) in text and --json. go test ./... passed in every package.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/workflow -run 'TestChecksStep' -v -count=1; same env: go test ./internal/run -run 'TestChecksStep|TestNextCheckActions|TestCheckPane' -v -count=1; go test ./cmd/rhiza -run 'TestScripts/review-checks-step' -v -count=1; go test ./... -count=1; commit 87ab105; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-checks-step](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
