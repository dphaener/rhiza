# review-critic-verified

evidence · `01M5ZW15T8477W57HNTW5P6TVV` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran the leaf's acceptance tests at ab0a2f3 in a detached worktree; all pass. Clean check: TestCriticNoneForCleanOrMinor gives alpha (no findings) outcome passed, no critic started (started = alpha/beta/gamma check 1 only). Minor+nit only: same test, beta with f1 minor and f2 nit passed, no critic. Major dropped: TestCriticDropsMajor starts 'alpha critic 1' in its own checkout with RHIZA_CHECK/CYCLE/ROLE=critic env, critic drops f1, alpha passed with detail 'cycle 1: critic dropped f1', critic checkout removed. Blocker confirmed: TestCriticConfirmsBlocker, critic confirms f1 blocker, beta cycle 1 failed 'cycle 1: confirmed f1 (blocker)' and re-runs in cycle 2. Verdict adding a finding refused: review-findings.txtar has alpha's critic run 'run verdict --check alpha' on verdict-f9.json, logged 'alpha critic-verdict-unknown exit=2 lines=same' (exit 2, journal unchanged) with stderr 'verdict 1: check alpha's round 1 holds no finding "f9"'. Also passing: TestCriticNoVerdictConfirms, TestCriticFails, TestCriticParallel, TestNextCheckActions, TestEnvCheckCycle, all TestScripts/review-* and go test ./...

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/run -run 'TestCritic|TestNextCheckActions' -v -count=1; same env go test ./internal/api -run TestEnvCheckCycle -v -count=1; same env go test ./cmd/rhiza -run 'TestScripts/review-' -v -count=1; same env go test ./... -count=1; commit ab0a2f3; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-critic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
