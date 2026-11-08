# watchdog-pr-owned-by-run-verified

evidence · `01M6XEZHRKE5V8878DZSXN2DV8` · actor [claude-code-e4cda6c3](../actor/claude-code-e4cda6c3.md)

## result

PASS

## statement

Verifier claude-code-e4cda6c3 at 468bae82 ran go test -trimpath -count=1 ./internal/api -run 'TestWatchdog(GreenUnmerged|ReadPR|CIRerun|PRNoCI)' with the fake gh: TestWatchdogCIRerun, TestWatchdogPRNoCI, TestWatchdogGreenUnmerged, TestWatchdogGreenUnmergedReadsOnce, TestWatchdogReadPR and TestWatchdogGreenUnmergedSkipsBusyChange all PASS; go vet and gofmt clean. Read the diff: watchdogReadPR returns an empty head unless readShipState names st.ID as Run (no record or Run empty gives every run no PR; PLAN gets none while SHIP owns it), and all three passes (green-unmerged, CI-rerun, no-CI) skip an empty head, so a succeeded plan run gets no hold, resume or merge (the 0568a71 case); watchdogReadOnce caches by run ID, so the owner's read does not leak to the plan run; watchdogGreenUnmerged also skips a change set any of whose runs is live or has a continue queued (tested both ways), while an idle succeeded run still gets its hold and the existing owner-run hold/resume tests still pass.

## Links

### Out

- verifies [watchdog-pr-owned-by-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-owned-by-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
