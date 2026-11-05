# ship-check-ci-test-findings

capability · `01M6MC4GXTJ7SWF8NV4SYXD3F9` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/api parses sample --log-failed output (a nested TestScripts/review-check-worktree failure, two failing tests in one package, and an infra-only log with no '--- FAIL:' line) into the expected test names and none for the infra log; a testscript with a fake gh whose check fails again after its rerun and whose 'run view 42 --log-failed' prints a '--- FAIL: TestScripts/x' line shows ship check --json listing that test as a failure and finding; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' only

## statement

When a PR head's checks still fail after their once-per-head rerun, rhiza ship check reads each failed Actions run's log (gh run view <run> --log-failed), reports the failing Go tests and testscripts it names ('--- FAIL: TestX[/sub]' lines) as test failures, and records each as a finding with no leaf for the next fix step; a failed check with no Actions run id, no parsable failing test, or only an infra/setup error stays a plain failure, as today

## Links

### Out

- decided-by [ship-change-ci-fix-cycle](../../../../../decision/ship-change-ci-fix-cycle.md)

### In

- depends-on [ship-change-ci-fix-step](ship-change-ci-fix-step.md)
- verifies [ship-check-ci-test-findings-verified](../../../../../evidence/ship-check-ci-test-findings-verified.md)
- verifies [ship-check-ci-test-findings-qa-1](../../../../../evidence/ship-check-ci-test-findings-qa-1.md)
- verifies [ship-check-ci-test-findings-verified-2](../../../../../evidence/ship-check-ci-test-findings-verified-2.md)
- verifies [pr-280-ship-change-fixes-red-pr-ci](../../../../../evidence/pr-280-ship-change-fixes-red-pr-ci.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
