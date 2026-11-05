# pr-280-ship-change-fixes-red-pr-ci

evidence · `01M6MXRQNH1V234M3Z28AHGY3J` · actor [dphaener](../actor/dphaener.md)

## head

808e2b481fcebea2898e1a9970334fb440c9b883

## pr

https://github.com/dphaener/rhiza/pull/280

## result

1 check passed: check SUCCESS

## statement

Pull request #280 (https://github.com/dphaener/rhiza/pull/280) passed CI at 808e2b4

## status

verified

## Links

### Out

- verifies [ship-check-ci-test-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-test-findings.md)
- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
