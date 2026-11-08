# question-fe0e69a6-1

question · `01M6XBAT3QEZCNZJ3NFCME8ZME`

## session

fe0e69a6-e20b-498b-876f-01a1bab33d46

## statement

Fix 01M6XB5 (wrong-run stalled-PR inbox item): the watchdog finds a run's PR by its change set's branch, so any ended run of that change set claims the PR. The ship record (.git/rhiza/ship/<cs>.json) will name the run that opened or pushed the PR. Where should that ownership check go? I recommend the shared PR read (watchdogReadPR): the CI-rerun and no-CI merge passes have the same mix-up and would rerun CI or merge main into a PR another run owns. / Fix 01M6XB5: ship records written before this change name no run. How should the watchdog treat a PR whose ship record names no run? I recommend matching it to no run: the watchdog then skips it (no wrong inbox item, no resume), and records written from now on carry the run.

## Links

### In

- answers [watchdog-pr-owned-by-ship-record-run](../decision/watchdog-pr-owned-by-ship-record-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
