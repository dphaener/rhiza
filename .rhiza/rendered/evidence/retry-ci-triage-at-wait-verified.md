# retry-ci-triage-at-wait-verified

evidence · `01M6YR97VF29CSEYZA5GKDQAAB` · actor [claude-code-11c0b7c0](../actor/claude-code-11c0b7c0.md)

## result

PASS

## statement

Verified at 1bf8725d. Read the diff: RetryStep(def, st, failed) keeps the sole-wait rule, then walks st.Steps back from the failed visit. It skips the failed step's own visits; retried entries are never in State.Steps. It follows command steps whose on.fail routes into the next step, and returns the sole wait whose on.pass is such a check, else failed. Traced the built-in ship-change: pr-green -> conflict-triage -> ci-triage gives wait-pr; rebased-green -> conflict-triage -> ci-triage gives wait-rebased. engine.go and api/run.go pass the state. Ran with RHIZA_* unset and GOTMPDIR/TMPDIR in a scratch dir. Inherited GOTMPDIR sits under ~/code/rhiza and breaks rhiza init in scripts. go test -trimpath ./internal/run/... ok. TestRetryStepAtWaitBeforeCheck covers: via a -> wait-a; via b (journal picks) -> wait-b; through a middle step -> wait-c; three triage visits -> wait-b; no predecessor, lone step, and previous not routing in -> the failed step. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-resume-retry|run-ship-retry)' passed: run-resume-retry, run-resume-retry-wait, run-ship-retry-pending and the new run-ship-retry-ci-triage. That script uses a fake gh, red then green: dry-run says retried at step wait-pr, resume succeeds, wait-pr visit 2 (retried) passes, then pr-green passes. internal/api: 533 tests passed and none failed, including TestRunResumeRetry and the other Resume/Retry tests. The package overran 8-15 min timeouts on load-heavy tests unrelated to this diff (TestReadBudget 151s, TestLinkEndpointKinds 226s alone, which passed) at VM load 10-16. CI runs the full package. gofmt clean, go vet ./internal/run ok.

## Links

### Out

- verifies [retry-ci-triage-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-ci-triage-at-wait.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
