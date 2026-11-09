# ci-triage-infra-waits-verified-2

evidence · `01M6YWZT7BZFC6G8MMA53FTWBS` · actor [claude-code-27fa7891](../actor/claude-code-27fa7891.md)

## result

PASS

## statement

Verifier, attempt 2, at 3ef516fe (a fixer commit that only changes the model). Read the diff: ship-change-ci-fix-step now says ci-triage on.fail goes to wait-ci-rerun, which matches internal/workflow/builtin/ship-change.toml:103. retry-reopens-at-wait now says the retry reopens at the wait the journal visited last, and at the failed step only when the journal visited none, which matches waitFor in internal/run/journal.go:1837. That resolves model findings f1 and f2. Re-ran the acceptance with the run's RHIZA_* and GOTMPDIR unset: go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-rerun-fail|run-ship-infra-wait|run-ship-ci-fix|run-ship-change-fail)' ok; ./internal/run ok; ./internal/workflow ok; ./internal/api: every test PASS, including TestShipCheckFailed, TestShipCheckSettledChanged, TestShipCheckTestFailures and TestShipCheckSettledReadsNoLog. Two failures are environment-only, neither caused by this change. (1) The package-exit guard flagged ~/code/.rhiza-root/usage.live.json changing; usage files are not touched by this change. (2) TestReadBudget (budget_test.go, not touched by this change set) ran past the timeout at VM load average about 14. make self-check fails only on the verifier's own untracked actor record, which needs render.

## Links

### Out

- verifies [ci-triage-infra-waits](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-triage-infra-waits.md)

## History

- `01M6YGFGQT0A9VRZ4AZ9X6VSMJ` Fix rhiza run resume --retry on a ship-change run that failed at ci-triage can never recover once the PR is green. Run 01M6XVTXGJAMQRBZKDZKY0Z01T (PR #393): pr-green failed because GitHub refused every CI job for billing (no test ran), so ci-triage ended 'exit 1, no failing test was recorded for the head: an infra or unparsable failure' and the run failed. After billing was fixed, one rerun turned PR #393 fully green, but 'rhiza run resume --retry' restarts at ci-triage, which re-reads the stale state and fails again in 1s (journal: retried ci-triage, visit 2, outcome fail). The retry should re-evaluate the PR first: when the head (or a new head) is green, continue at pr-green/rebase; an infra-only failure (no failing test) should send the run back to wait on CI rather than end it. Unstuck by hand: merged origin/main in the run worktree, pushed, then rhiza ship check / ship land / ship merge.
