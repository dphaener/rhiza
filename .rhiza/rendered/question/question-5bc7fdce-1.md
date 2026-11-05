# question-5bc7fdce-1

question · `01M6MBX19VWABE6S2R6FSNR46Q`

## session

5bc7fdce-fa6a-4b9d-9294-01a1a8bb94c3

## statement

Which intent should change set 01M6MBPT (Ship-change fixes red PR CI) serve? It follows up #259 (agents run only scoped tests, so CI runs the full suite), and #259's change set served milestone-25 (the daemon empties the backlog on its own) and milestone-24. / How should a red PR CI become a fix in ship-change? Today pr-green/rebased-green/landed-green/rebased-merge-green run `rhiza ship check`, which reruns failed jobs once per head and then fails the run. / Which failing ship-check stages get a CI-fix cycle, and where does the run go after the fix? / What cap should CI-fix cycles have before the run fails?

## Links

### In

- answers [ship-change-ci-fix-cycle](../decision/ship-change-ci-fix-cycle.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
