# ship-record-names-run-qa-1

evidence · `01M6XFCHV43PY0F39Y1JK874T4` · actor [claude-code-68417e7f](../actor/claude-code-68417e7f.md)

## result

PASS

## statement

Built rhiza from the branch head and drove it by hand in a scratch repo (bare origin, fake gh answering pr view as PR #7 green). rhiza ship check with no RHIZA_RUN left the ship record's run unset; RHIZA_RUN=runA ship check wrote run=runA; RHIZA_RUN=runA ship land pushed and kept run=runA; RHIZA_RUN=runB ship check --change switched it to run=runB; a later ship check with no RHIZA_RUN kept runB; after main moved, RHIZA_RUN=runC ship rebase pushed, cleared head and wrote run=runC. internal/run/engine.go sets RHIZA_RUN for every step. go test -trimpath ./internal/api -run 'TestShip.*Run|TestShipLand|TestWatchdogPRNoCI' passes (with the inherited RHIZA_ROOT unset).

## Links

### Out

- verifies [ship-record-names-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/ship-record-names-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.
