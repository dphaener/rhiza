# fix-01m6xb5-plan

source · `01M6XBDC2TNMPHBESTDGGXSM1W`

## quotes

```json
[]
```

## raw_sha256

a11c58f5454d93363c6f6c6388dc9bd185a7d0ab4d320c785244c2e16b1d00ef

## session

fe0e69a6-e20b-498b-876f-01a1bab33d46

## summary

Autonomous plan run, no owner. The planner read internal/api/watchdog_unmerged.go and daemon_watchdog.go: watchdogReadPR finds a run's PR by its change set's branch and checks only that the PR body names the change set, so every ended run of a change set claims the change set's PR. That is how the plan run 01M6VY7E4JWY got the inbox item 'merge it by hand' for PR #376, which ship-change run 01M6WMX236BTT04TYQDACGKBZ4 owned and was being resumed toward its merge. The ship record (.git/rhiza/ship/<cs>.json) is per change set and names no run today. Two questions went to the arbiter, who ruled (assumed, would have escalated): the ownership check lives in the shared PR read, so the green-unmerged, CI-rerun and no-CI passes all see no PR for a run the ship record does not name; a record that names no run matches no run. The planner split the work into two leaves: ship steps record their run in the ship record, and the watchdog's PR read checks it, with the green-unmerged pass also skipping a change set with any live run or queued continue, as the intent asks.

## title

Plan of fix 01M6XB5: watchdog matches a PR only to the run that owns it (autonomous)

## Links

### In

- derived-from [watchdog-pr-owned-by-ship-record-run](../decision/watchdog-pr-owned-by-ship-record-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
