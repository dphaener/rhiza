# ship-check-ci-test-findings-qa-1

evidence · `01M6MS3MVG7AC7MZFTN68MGRY8` · actor [claude-code-7f0b608d](../actor/claude-code-7f0b608d.md)

## result

PASS

## statement

QA at dc3c5d5c with a built binary (go build ./cmd/rhiza) in a scratch repo with a bare origin, an accepted change set, a pushed cs/ship-it branch and a fake gh whose PR has two failing checks (test on run 42, lint on run 43). First rhiza ship check: 'not settled: reran the failed jobs of runs 42, 43', no log read. Second: 'failed: test, lint; failing tests: TestScripts, TestScripts/review-check-worktree, TestFoo' (nested timestamped '--- FAIL:' lines parsed, TestFoo named in both logs listed once), gh ran 'run view 42 --log-failed' and 'run view 43 --log-failed', and --json findings carry no leaf with a reason naming run, check and short head. A new pushed head cleared the old findings on its rerun; an infra log ('Error: setup-go failed') after the rerun gave a plain 'failed: test, lint' with no tests; a green check reported '2 checks passed' and the state file no longer held ci_findings. Also ran go test -trimpath ./internal/api/... -run 'TestShip|TestParseFailed' and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|ship-check-ci-test-findings|describe-golden)': ok

## Links

### Out

- verifies [ship-check-ci-test-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-test-findings.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
