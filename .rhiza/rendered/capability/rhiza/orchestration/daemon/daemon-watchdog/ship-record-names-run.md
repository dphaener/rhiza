# ship-record-names-run

capability · `01M6XBDT9GM17AEB5M6N0RB4KX` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShip.*Run|TestShipLand|TestWatchdogPRNoCI' shows the ship record naming the run after ship.open-pr, ship.check, ship.land and ship.rebase run under one RHIZA_RUN, naming a second run after that run's ship.check on the same change set, and naming the run the watchdog acted for after its no-CI merge of main pushes the PR

## statement

Each ship step that opens, checks or pushes a change set's PR (ship.open-pr, ship.check, ship.land, ship.rebase, and the watchdog's merge of main on behalf of a run) records in the change set's ship record the run it ran for (RHIZA_RUN, or the run the watchdog acted for), so the record names the run that opened or last pushed the PR

## Links

### Out

- decided-by [watchdog-pr-owned-by-ship-record-run](../../../../../decision/watchdog-pr-owned-by-ship-record-run.md)

### In

- verifies [ship-record-names-run-verified](../../../../../evidence/ship-record-names-run-verified.md)
- verifies [ship-record-names-run-qa-1](../../../../../evidence/ship-record-names-run-qa-1.md)
- verifies [pr-382-fix-the-inbox-flags-a-stalled](../../../../../evidence/pr-382-fix-the-inbox-flags-a-stalled.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
