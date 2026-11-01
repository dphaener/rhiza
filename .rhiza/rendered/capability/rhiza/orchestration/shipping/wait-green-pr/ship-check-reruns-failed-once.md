# ship-check-reruns-failed-once

capability · `01M6A0NJWB1HNVZ6QZAR17GDMJ` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A testscript with a fake gh whose pr view reports a failed check with detailsUrl .../actions/runs/42/job/9 until gh run rerun is called, then success: gh.log holds exactly one 'run rerun 42 --failed' and the ship-change run passes pr-green; a second case where the check still fails after the rerun fails at pr-green with exactly one rerun in gh.log; a Go test covers the --main path reading databaseId

## statement

When a PR or main check fails, rhiza ship check reruns the failed jobs of that check's GitHub Actions run once per head, records the rerun in the ship record and reports not yet settled; a check that fails again on the same head fails as before, and a check with no Actions run id is never rerun

## Links

### Out

- decided-by [ci-rerun-once-per-head](../../../../../decision/ci-rerun-once-per-head.md)

### In

- verifies [ship-check-reruns-failed-once-verified](../../../../../evidence/ship-check-reruns-failed-once-verified.md)
- verifies [ship-check-reruns-failed-once-qa-1](../../../../../evidence/ship-check-reruns-failed-once-qa-1.md)
- verifies [pr-155-ci-rerun-once](../../../../../evidence/pr-155-ci-rerun-once.md)

## History

- `01M6A0JQ2WHGH7DDFKJGRXJD85` A failed CI check on a run's PR or main is rerun once per head before the run fails
