# watchdog-pr-owned-by-run

capability · `01M6XBDTCS3RKPWJHKKRT0EXKH` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run 'TestWatchdog(GreenUnmerged|ReadPR|CIRerun|PRNoCI)', with a fake gh, shows: a succeeded plan run of a change set whose ship record names a later ship-change run gets no inbox item, no resume and no merge (the inbox item 0568a71 case); the same with the ship-change run live or its continue queued; a ship record naming no run gives every run no PR; and the owning run still gets its hold or resume as before

## statement

The watchdog's shared PR read gives a run its change set's PR only when the ship record names that run, and no PR when it names another run or none, so the green-unmerged, CI-rerun and no-CI passes never act on or raise an inbox item about a PR another run owns; the green-unmerged pass also skips a change set any of whose runs is live or has a continue queued

## Links

### Out

- decided-by [watchdog-pr-owned-by-ship-record-run](../../../../../decision/watchdog-pr-owned-by-ship-record-run.md)
- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-pr-owned-by-run-verified](../../../../../evidence/watchdog-pr-owned-by-run-verified.md)
- verifies [pr-382-fix-the-inbox-flags-a-stalled](../../../../../evidence/pr-382-fix-the-inbox-flags-a-stalled.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
