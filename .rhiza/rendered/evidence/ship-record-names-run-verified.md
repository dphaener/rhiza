# ship-record-names-run-verified

evidence · `01M6XEQR88XM64J5FNH736NEAE` · actor [claude-code-6e37fcb8](../actor/claude-code-6e37fcb8.md)

## result

PASS

## statement

At 30ccbbc0 ran go test -trimpath ./internal/api -run 'TestShip.*Run|TestShipLand|TestWatchdogPRNoCI|TestShipOpenPR' -count=1: every test passed, including TestShipOpenPR (record Run=R1 under RHIZA_RUN), TestShipCheckNamesRun (01RUNA, then a second run's ship.check on the same change set takes it over as 01RUNB), TestShipLandNamesRun, TestShipRebaseNamesRun (01RUNA) and TestWatchdogPRNoCI (the no-CI merge of main records the watchdog's run 01RUNW). Read the diff: shipState gains Run; ship.open-pr, ship.check, ship.land, ship.rebase set it from RHIZA_RUN, and daemonWatchdog's merge passes st.ID to watchdogMergeMain. The package binary sometimes exits non-zero after PASS because rootguard flags .rhiza-root/usage.live.json, which live agent statuslines write, not these tests (the new tests alone passed on a rerun); filed as draft 01M6XEPTM57C88AEYN4SPME6M5.

## Links

### Out

- verifies [ship-record-names-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/ship-record-names-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
