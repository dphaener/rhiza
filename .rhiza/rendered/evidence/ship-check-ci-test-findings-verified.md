# ship-check-ci-test-findings-verified

evidence · `01M6MNXVV84CZWVST1RK2QK1MH` · actor [claude-code-eb823135](../actor/claude-code-eb823135.md)

## result

PASS

## statement

Verified at 91caa11e. go test -trimpath ./internal/api/... passed, including the parseFailedTests test in ship_test.go: a nested TestScripts/review-check-worktree failure gives TestScripts and TestScripts/review-check-worktree, two failures in one package give TestA and TestB, and an infra-only log gives none. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|ship-check-ci-test-findings|describe-golden)' -v passed all 14 scripts. ship-check-ci-test-findings.txtar uses a fake gh whose check fails again after its rerun and whose 'run view 42 --log-failed' prints '--- FAIL: TestScripts/x'. ship check --settled --json lists TestScripts/x under tests and findings and records ci_findings in the state file. With an infra-only log the test is a plain failure and ci_findings is cleared, and a green check clears it too. This fixes attempt 1's stale-findings finding. I read the diff in internal/api/ship.go: shipReadFailures replaces st.CIFindings on each read and clears it for a head that was not rerun, shipRerunFailed clears it on a new rerun, and the green branch clears it. Note: the acceptance's -run 'TestScripts/run-ship' pattern does not match the script's name, so I ran it by name.

## Links

### Out

- verifies [ship-check-ci-test-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-test-findings.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
