# ship-ci-fix-plan

source · `01M6MC26DA3HZ8V0PDTPQ93122`

## quotes

```json
[]
```

## raw_sha256

f896093907b53110e63bf91e2293c11cf70c6d9506538977f5c1730bb2484f3a

## session

5bc7fdce-fa6a-4b9d-9294-01a1a8bb94c3

## summary

Autonomous plan run, no owner. The planner read ship-change.toml, review.toml, internal/api/ship.go (ship check and its once-per-head rerun of failed jobs) and the model (milestone-25, #259's change set and its serves decision). It raised one question covering four topics (which intent it serves, the mechanism, which stages fix, the cap). The arbiter answered with Decision ship-change-ci-fix-cycle: serve milestone-25. ship check keeps its flake rerun, then parses failing tests from gh run view --log-failed into findings with no leaf; anything unparsable still fails as today. A ci-fix fixer step runs on pr-green and rebased-green failures, at most 2 cycles per run, and returns through render-pr, open-pr, wait-pr, pr-green and signoff. landed-green and rebased-merge-green are unchanged. The planner wrote two leaves under wait-green-pr: test-failure findings in ship check, and the ci-fix step and its routing.

## title

Plan of Ship-change fixes red PR CI (autonomous)

## Links

### In

- derived-from [ship-change-ci-fix-cycle](../decision/ship-change-ci-fix-cycle.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
