# watchdog-pr-owned-by-ship-record-run

decision · `01M6XBCBKRAR1HDT7ZGPE2FP0B`

## reasoning

watchdogReadPR (internal/api/daemon_watchdog.go) is the single read behind d.ci, used by every watchdog pass; it finds the PR by the change set's branch and only checks prOwnedBy(body, change set), so any ended run of the change set claims the PR. Putting the run check there fixes the reported wrong-run inbox item and also stops the rerun and merge passes acting on a PR a live or resumed run owns (the PR #376 case), with one check instead of three. For legacy records with no run, matching no run is the safe failure: the cost is the watchdog no longer rescuing a pre-change stalled PR (the owner can still act), while guessing an owner could reproduce the exact bug this change set fixes. Legacy records age out as runs write new ones.

## statement

The watchdog matches a PR to a run only when the change set's ship record (.git/rhiza/ship/<cs>.json) names that run as the one that opened or pushed the PR. The check lives in the shared PR read, watchdogReadPR, so the stalled-PR inbox item, the CI-rerun pass and the no-CI merge pass all see no PR for a run that does not own it. A ship record that names no run (written before this change) matches no run: the watchdog skips that PR, and records written from now on carry the run.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-fe0e69a6-1](../question/question-fe0e69a6-1.md)
- derived-from [fix-01m6xb5-plan](../source/fix-01m6xb5-plan.md)

### In

- decided-by [ship-record-names-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/ship-record-names-run.md)
- decided-by [watchdog-pr-owned-by-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-owned-by-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
